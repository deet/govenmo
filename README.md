# govenmo

A Venmo API client library in Golang.

See the Venmo API documentation here at https://developer.venmo.com/docs/oauth.

See README for usage or GoDoc.

<a href="https://godoc.org/github.com/deet/govenmo"><img src="https://godoc.org/github.com/deet/govenmo?status.svg" alt="GoDoc"></a>

## Concepts

Since all Venmo API requests must be authenticated using a user's OAuth token, the basic object is the Account, which holds information for the authenticated user and tokens.

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

	sentPayment, err := account.PayOrCharge(target, 5.27, "Thanks for the govenmo library!", "public")
	if err != nil {
		// Handle error ...
	}

### Refresh single payment

	payment := &Payment{}
	payment.Id = "1111111111111111111"
	err = account.RefreshPayment(payment)
	if err != nil {
		// Handle error ...
	}

### Fetch multiple payments

	var updatedSince time.Time
	payments, err := account.PaymentsSince(updatedSince)
	if err != nil {
		// Handle error ...
	}
	for _, payment := range payments {
		// Each payment is an instance of Payment.
		log.Println("Found payment:", payment.Note)
	}

### Complete a charge

	updatedPayment, err := account.CompletePayment("paymentID", "approve")
	if err != nil {
		// Handle error ...
	}

### Deny a charge

	updatedPayment, err := account.CompletePayment("paymentID", "deny")
	if err != nil {
		// Handle error ...
	}

### Cancel a charge

	updatedPayment, err := account.CompletePayment("paymentID", "cancel")
	if err != nil {
		// Handle error ...
	}

### Fetch friends

	friends, err := account.FetchFriends()
	if err != nil {
		// Handle error ...
	}

	for _, friend := range friends {
		// Each friend is an instance of User.
		log.Println("Found friend:", friend.DisplayName)
	}

## Settings

Enable Venmo sandbox mode. Note that the Venmo sandbox doesn't behave exactly like the production API.

	govenmo.Environment = "sandbox"

Use local sandbox

	// In your client
	govenmo.Environment = "local_sandbox"

Set a maximum payment or charge amount. (Why?... because when you first start using the production API you probably don't want a bug in your code to be able to send thousands of dollars.)

	max := float64(50)
	govenmo.MaxPayment = &max

Enable logging

	govenmo.EnableLogging(nil)  // you can also pass a Logger

### Run the local sandbox

The package local_sandbox mimics the real Venmo sandbox so that you don't have to hit it as much during testing. 

Local sandbox returns the sandbox's hardcoded POST /payments responses. It also mimics GET /me and GET /payments/1111111111111111111.  Other requests are proxied to the real sandbox and would require a valid token.

Like the real sandbox, it's not a replica of the Venmo production API and the values returned in the responses might not be the same as what you send in your request.

	cd local_sandbox
	go run main.go

The local sandbox must be running for the (limited) tests.

## License

MIT. See LICENSE.md.
