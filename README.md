# govenmo

A Venmo API client in Golang.

See the Venmo API documentation here at https://developer.venmo.com/docs/oauth

## Concepts

Since all Venmo API requests must be authenticated using a user's Oauth token, the basic object is the Account, which holds information for the authenticated user and tokens.

The library assumes you have already obtained a user access token somehow, for example using the OAuth flow. Helpers for completing the OAuth flow might be added in the future.

## Usage

### Create account with user access token

Obtaining the token is not included in the library but is easy to implement.

	account := Account{
		AccessToken:  "...",
		RefreshToken: "...",
	}

### Fetch user information including balance

	err := account.Refresh()
	if err != nil {
		// Handle error ...
	}

### Pay someone

	target := Target{}
	target.Email = "kbrisson@gmail.com"
	// OR target.Phone = "..."
	// OR target.User.Id = "..."

	sentPayment, err := account.Pay(target, 5.27, "Thanks for the govenmo library!", "public")
	if err != nil {
		// Handle error ...
	}

### Complete a charge

	updatedPayment, err := account.Complete("paymentID", "approve")
	if err != nil {
		// Handle error ...
	}

### Deny a charge

	updatedPayment, err := account.Complete("paymentID", "deny")
	if err != nil {
		// Handle error ...
	}

### Cancel a charge

	updatedPayment, err := account.Complete("paymentID", "cancel")
	if err != nil {
		// Handle error ...
	}

### Fetch friends

	friends, err := account.FetchFriends()
	if err != nil {
		// Handle error ...
	}

	for _, friend := range friends {
		// Each friend is an Account with the User fields populated
		log.Println("Found friend:", friend.DisplayName)
	}

## Settings

Enable Venmo sandbox mode. Note that the Venmo sandbox doesn't behave exactly like the production API.

	govenmo.Environment = "sandbox"

Use local sandbox

	// In your client
	govenmo.Environment = "local"

Set a maximum payment or charge amount. (Why?... because when you first start using the production API you probably don't want a bug in your code to be able to send thousands of dollars.)

	max := float64(50)
	govenmo.MaxPayment = &max

Enable logging

	govenmo.EnableLogging(nil)  // you can also pass a Logger

### Run the local sandbox

The package local_sandbox mimics the real Venmo sandbox so that you don't have to hit it as much during testing. It returns the sandbox's hardcoded POST /payments responses and proxies other requests to the real sandbox.

Like the real sandbox, it's not a replcate of the Venmo production API and the values returned in the responses might not be the same as what you send in your request.

	cd local_sandbox
	go run main.go

The local sandbox must be running for the (limited) tests.

Local sandbox responses deviate from real sandbox responses in the error checking is not as strict and that some bugs have been fixed to match production API better.

## License

MIT. See LICENSE.md.
