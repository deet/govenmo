package govenmo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type User struct {
	Username          string  `json:"username"`
	Id                string  `json:"id"`
	Email             *string `json:"email"`
	DisplayName       string  `json:"display_name"`
	FirstName         string  `json:"first_name"`
	LastName          string  `json:"last_name"`
	Phone             *string `json:"phone"`
	About             string  `json:"about"`
	ProfilePictureUrl string  `json:"profile_picture_url"`
	FriendsCount      int64   `json:"friends_count"`
	IsFriend          *bool   `json:"is_friend"`
	DateJoined        Time    `json:"date_joined"`
}

// FetchFriends retrieves all Venmo friends for an Account.
// It follows 'next' links.
func (account *Account) FetchFriends() (friends []User, err error) {
	next := ""

	for {
		url := ""
		if next != "" {
			url = next
		} else {
			url = apiRoot() + "/users/" + account.Id + "/friends?"
		}
		url += "&access_token=" + account.AccessToken
		logger.Println("Fetching url for user's friends:", url)

		resp, err := http.Get(url)
		if err != nil {
			logger.Println("Could get response from Venmo:", err)
			return friends, err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logger.Println("Could not parse response from Venmo:", err)
			return friends, err
		}

		logger.Println("Received response from GET /friends: ", string(body), "", "")

		var parsedResponse *userFriendsResponse = &userFriendsResponse{}
		err = json.Unmarshal(body, &parsedResponse)
		if err != nil {
			logger.Println("Could not parse response from Venmo:", err)
			return friends, err
		}

		if parsedResponse.Error.Code != 0 {
			logger.Println("Venmo friend fetch returned error:", parsedResponse.Error.Message)
			return friends, errors.New("Venmo friend fetch returned error: " + parsedResponse.Error.Message)
		}

		logger.Println("Next: ", parsedResponse.Pagination.Next)

		for _, friend := range parsedResponse.Data {
			logger.Printf("Received friend: %+v\n", friend)
			friends = append(friends, friend)
		}

		next = parsedResponse.Pagination.Next
		if next == "" {
			break
		} else {
			logger.Println("Fetching more friends")
		}
	}

	err = nil

	return
}
