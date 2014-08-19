package govenmo

import (
	"testing"
)

func TestPayOrCharge(t *testing.T) {
	account := &Account{}
	account.AccessToken = "faketoken"
	target := Target{}

	Environment = "local_sandbox"
	EnableLogging(nil)

	payment, err := account.PayOrCharge(target, 0.09, "", "public")

	if err == nil {
		t.Error("Sandbox should have errored on invalid amount")
	}
	if payment.Id != "" {
		t.Error("Payment should not have ID")
	}

	payment, err = account.PayOrCharge(target, 0.10, "", "public")
	if err == nil {
		t.Error("Sandbox should have errored with non sandbox user")
	}

	target.Email = "venmo@venmo.com"
	payment, err = account.PayOrCharge(target, 0.10, "", "public")
	if err != nil {
		t.Error("Sandbox should not have errored with non sandbox user")
	}

	if payment.Id != "1322585332520059420" || payment.Actor.Username != "delavara" || payment.Target.User.Id != "145434160922624933" {
		t.Error("Wrong payment ID, actor username, or target user ID")
	}

	if payment.Amount != 0.10 {
		t.Error("Wrong payment amount")
	}

}

func TestPaymentRefresh(t *testing.T) {
	account := &Account{}
	account.AccessToken = "faketoken"

	Environment = "local_sandbox"
	EnableLogging(nil)

	payment := &Payment{}

	payment.Id = "ddd"
	err := account.RefreshPayment(payment)

	if err == nil {
		t.Error("Sandbox should have errored on payment ID")
	}

	payment.Id = "1111111111111111111"
	err = account.RefreshPayment(payment)

	if err != nil {
		t.Error("Sandbox should not have errored test payment ID")
	}

	if payment.Id != "1111111111111111111" ||
		payment.Target.User.Username != "someone-else" ||
		payment.Actor.DisplayName != "Keith Brisson" ||
		payment.Note != "The Meatball Shop" ||
		payment.Action != "pay" {

		t.Error("Wrong payment info")
	}

	if payment.Amount != 6 {
		t.Error("Wrong payment amount")
	}

}
