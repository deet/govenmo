// Package govenmo provides a Venmo client.
// Use it to retrieve payments, fetch account infromation,
// make payments, complete charges, and list Venmo friends.
// You must provide your own OAuth tokens.
package govenmo

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

// Account is the basic type used for all API calls in govenmo. To make an API call
// you should create and Account with valid OAuth tokens. Account includes User.
type Account struct {
	AccessToken  string  `json:"access_token"`
	RefreshToken string  `json:"refresh_token"`
	Balance      float64 `json:"balance,string"`
	ExpiresIn    int64   `json:"expires_in"`
	TokenType    string  `json:"bearer"`
	User         `json:"user"`
}

// Refresh retrieves account information, including balance and biographical info
// from the Venmo api.
func (a *Account) Refresh() error {
	url := apiRoot() + "/me?access_token=" + a.AccessToken
	logger.Println("account refresh using URL:", url)
	resp, err := http.Get(url)
	if err != nil {
		logger.Println("Could get response from Venmo:", err)
		return err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	logger.Println("Received response from GET /me: ", string(body), "", "")

	if err != nil {
		logger.Println("Could not read response from Venmo:", err)
		return err
	}

	var parsedResponse *userGetResponse = &userGetResponse{}
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		logger.Println("Could not parse response from Venmo:", err)
		return err
	}

	logger.Printf("Parsed response from GET /me: %+v\n", *parsedResponse)

	a.User = parsedResponse.Data.User
	a.Balance = parsedResponse.Data.Balance

	//venmoAccount.Username = parsedResponse.Data.Username
	//venmoAccount.About = parsedResponse.Data.About
	//venmoAccount.Balance = parsedResponse.Data.Balance
	//venmoAccount.FriendsCount = parsedResponse.Data.FriendsCount
	//venmoAccount.FirstName = parsedResponse.Data.FirstName
	//venmoAccount.LastName = parsedResponse.Data.LastName
	//venmoAccount.Email = parsedResponse.Data.Email
	//venmoAccount.Phone = parsedResponse.Data.Phone
	//venmoAccount.ProfilePictureUrl = parsedResponse.Data.ProfilePictureUrl

	return nil
}
