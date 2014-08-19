package govenmo

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Payment stores a payment retrieved from the Venmo API. See the Venmo API docs.
type Payment struct {
	Id            string
	Status        string
	Action        string
	Actor         User
	Amount        float64
	Audience      string
	DateCompleted *Time `json:"date_completed"`
	DateCreated   *Time `json:"date_created"`
	Note          string
	Target        Target
	Fee           *float64
	Refund        *string
	Medium        string
}

//

type recentPaymentsResponse struct {
	Pagination Pagination
	Data       []Payment
}

// PaymentsSince fetches payments for an Account updated since a Time. Note that
// Venmo's 'updated at' logic is somewhat imprecise.
// There is currently no way to specify a limit, and PaymentsSince will follow 'next'
// links to retrieve the entire result set.
func (a *Account) PaymentsSince(updatedSince time.Time) (payments []Payment, err error) {
	next := ""

	for {
		url := ""
		if next != "" {
			url = next
		} else {
			url = apiRoot() + "/payments?"
			url += "after=" + updatedSince.Format(VenmoTimeFormat)
		}
		url += "&access_token=" + a.AccessToken
		logger.Println("Fetching url for recent transactions:", url)

		resp, err := http.Get(url)
		if err != nil {
			logger.Println("Could get response from Venmo:", err)
			return payments, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Println("Could not parse response from Venmo:", err)
			return payments, err
		}

		logger.Println("Received response from GET /payments: ", string(body), "", "")

		var parsedResponse *recentPaymentsResponse = &recentPaymentsResponse{}
		err = json.Unmarshal(body, &parsedResponse)
		if err != nil {
			logger.Println("Could not parse response from Venmo:", err)
			return payments, err
		}

		for _, payment := range parsedResponse.Data {
			payments = append(payments, payment)
		}

		logger.Println("Next: ", parsedResponse.Pagination.Next)

		next = parsedResponse.Pagination.Next
		if next == "" {
			break
		} else {
			logger.Println("Fetching more")
		}
	}

	err = nil
	return
}

// PayOrCharge creates a Venmo payment with the Account as a Actor.
func (a *Account) PayOrCharge(target Target, amount float64, note string, audience string) (sentPayment Payment, err error) {
	logger.Println("Sending venmo payment")

	if MaxPayment != nil {
		if amount > *MaxPayment || amount < *MaxPayment {
			logger.Println("Will not do venmo transactions over", *MaxPayment, "for now. Tried to do amount:", amount)
			err = errors.New("Venmo transactions are limited in size for now.")
			return
		}
	}

	params := url.Values{}
	params.Set("access_token", a.AccessToken)

	if target.Email != "" {
		params.Set("email", target.Email)
	}

	if target.User.Id != "" {
		params.Set("user_id", target.User.Id)
	}

	url := apiRoot() + "/payments"

	params.Set("note", note)
	params.Set("amount", fmt.Sprintf("%f", amount))
	params.Set("audience", audience)

	logger.Printf("Sending venmo payment: %+v\n", params)
	resp, err := http.PostForm(url, params)
	if err != nil {
		logger.Println("Could post payment to Venmo:", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("Could not parse response from Venmo:", err)
		return
	}

	logger.Println("Venmo payment body:", string(body))

	var parsedResponse *postPaymentResponse = &postPaymentResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		logger.Println("Could not parse response from Venmo:", err)
		return
	}

	sentPayment = parsedResponse.Data.Payment

	if parsedResponse.Error.Message != "" {
		err = errors.New(parsedResponse.Error.Message)
		return
	}

	return
}

// CompletePayment allows you to 'approve', 'deny', or 'cancel' a pending charge request.
func (a *Account) CompletePayment(paymentId, action string) (updatedPayment Payment, err error) {
	logger.Println("Completing venmo payment", paymentId, "with action", action)

	params := url.Values{}

	url := apiRoot() + "/payments/" + paymentId + "?access_token=" + a.AccessToken

	params.Set("action", action)

	logger.Printf("Complete venmo payment: %+v\n", params)
	logger.Println("Using URL:", url)

	req, err := http.NewRequest("PUT", url, strings.NewReader(params.Encode()))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		logger.Println("Could not create PUT request to complete Venmo payment:", err)
		return
	}

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logger.Println("Could not PUT to complete Venmo payment:", err)
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("Could not parse response from Venmo:", err)
		return
	}

	logger.Println("Venmo payment body:", string(body))

	var parsedResponse *completePaymentResponse = &completePaymentResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		logger.Println("Could not parse response from Venmo:", err)
		return
	}

	if parsedResponse.Error.Message != "" {
		err = errors.New(parsedResponse.Error.Message)
		return
	}

	updatedPayment = parsedResponse.Data

	if updatedPayment.Id != "" {
		logger.Println("Updated venmo payment with ID:", updatedPayment.Id, "and status", updatedPayment.Status)
	} else {
		err = errors.New("Could not complete venmo payment")
	}

	return
}

// RefreshPayment updates a Payment object with the most current state from the Venmo API.
// For multiple requests using PaymentsSince would be advisable.
func (a *Account) RefreshPayment(payment *Payment) error {
	if payment == nil {
		return errors.New("Cannot refresh nil payment")
	}

	url := apiRoot() + "/payments/" + payment.Id
	resp, err := http.Get(url + "?access_token=" + a.AccessToken)
	if err != nil {
		logger.Println("Could get response from Venmo:", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logger.Println("Could not read response from Venmo:", err)
		return err
	}

	logger.Println("Venmo payment body:", string(body))

	var parsedResponse *getPaymentResponse = &getPaymentResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		logger.Println("Could not parse response from Venmo:", err)
		return err
	}

	if parsedResponse.Error.Message != "" {
		logger.Println("Error from Venmo API when refreshing payment:", parsedResponse.Error.Message)
		return errors.New(parsedResponse.Error.Message)
	}

	*payment = parsedResponse.Data

	return nil
}
