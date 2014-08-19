package govenmo

type userGetResponse struct {
	Data Account
}

type paymentPostData struct {
	Balance float64 `json:",string"`
	Payment Payment
}

type postPaymentResponse struct {
	Data  paymentPostData
	Error Error
}

type completePaymentResponse struct {
	Data  Payment
	Error Error
}

type getPaymentResponse struct {
	Data  Payment
	Error Error
}

type userFriendsResponse struct {
	Pagination Pagination
	Error      Error
	Data       []User
}

func apiRoot() string {
	switch Environment {
	case "local_sandbox":
		return "http://localhost:4000"
	case "sandbox":
		return "https://sandbox-api.venmo.com/v1"
	default:
		return "https://api.venmo.com/v1"
	}
}
