package auth

import (
	"respond"
	"testing"
)

var userStore UserStore
var auth *Auth

type testUserStore struct {

}

func init () {
	userStore = dbstore()
	auth = New(userStore)
}

func TestSignupUser(t *testing.T) {

	_, err := auth.SignupUser(emailGenerator.Next(), "")

	if validationError, ok := err.(*respond.DataError); !ok {
		t.Fatal("signup with an empty password should result in a validation error")
	} else {
		if validationError.Errors["Password"][0] != "password is required" {
			t.Fatal("empty password validation message not set")
		}
	}

	_, err = auth.SignupUser("invalidemailaddress", "password")

	if validationError, ok := err.(*respond.DataError); !ok {
		t.Fatal("signup with an invalid email address should result in a validation error")
	} else {
		if validationError.Errors["Email"][0] != "email address not valid" {
			t.Fatal("invalid email address validation message not set")
		}
	}

	_, err = auth.SignupUser("", "password")

	if validationError, ok := err.(*respond.DataError); !ok {
		t.Fatal("signup with an invalid email address should result in a validation error")
	} else {
		if validationError.Errors["Email"][0] != "email address required" {
			t.Fatal("empty email address validation message not set")
		}
	}

	var user *User
	emailAddress := emailGenerator.Next()

	user, err = auth.SignupUser(emailAddress, "password")

	if err != nil {
		t.Fatalf("user creation failed: %s", err)
	}

	if user == nil {
		t.Fatal("user not returned")
	}

	// TODO: check there is a new lovelist

	_, err = auth.SignupUser(emailAddress, "password")

	if validationError, ok := err.(*respond.DataError); !ok {
		t.Fatalf("signup with the email address of an existing user should result in a validation error, got: %s", err)
	} else {
		if validationError.Errors["Email"][0] != "email address taken" {
			t.Fatal("email address collision validation message not set")
		}
	}

}

func TestFindUserAgainstPassword(t *testing.T) {

	email := emailGenerator.Next()

	if _, err := userStore.CreateUser(email, "very-secret"); err != nil {
		t.Fatal(err)
	}

	user, err := auth.FindUserAgainstPassword(email, "very-secret")

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user not returned")
	}

	user, err = auth.FindUserAgainstPassword(email, "WRONG!!PASSWORD")

	if dataError, ok := err.(*respond.DataError); ok {
		if dataError.Errors["Password"][0] != "password mismatch" {
			t.Fatal("password mismatch error not returned")
		}
	} else {
		t.Fatal("returned error was not a *data.DataError")
	}

	if user != nil {
		t.Fatal("user returned against invalid password")
	}
}
