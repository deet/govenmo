package govenmo

import (
	"testing"
)

func TestAccountRefresh(t *testing.T) {
	account := &Account{}
	account.AccessToken = "faketoken"

	Environment = "local_sandbox"
	EnableLogging(nil)

	err := account.Refresh()

	if err != nil {
		t.Error("/me should not have errored:", err)
	}

	if account.Id != "123245678901232456789" || account.Balance != 1.23 || *account.Email != "email@example.com" {
		t.Error("Parsed user info is wrong")
	}
}
