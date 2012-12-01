package auth

import (
	"regexp"
	"respond"
	"respond/crypto"
)

const EmailValidationRegex string = ".+@.+" // very leanient, enough.

type UserStore interface {
	CreateUser(email string, password string) (*User, error)
	FindUserByEmail(string) (*User, error)
}

type Auth struct {
	store UserStore
}

func New(store UserStore) *Auth {
	return &Auth{store: store}
}

func (a *Auth) FindUserAgainstPassword(email, password string) (user *User, err error) {

	if user, err = a.store.FindUserByEmail(email); err != nil {

		if _, ok := err.(*respond.NotFoundError); ok {

			dataError := respond.NewDataError()
			dataError.Add("Email", "email address not found")

			return nil, dataError
		}

		return nil, err
	}

	if user.PasswordHash != crypto.Base64HashOfSaltedPassword(password, user.PasswordSalt) {

		dataError := respond.NewDataError()
		dataError.Add("Password", "password mismatch")

		return nil, dataError
	}

	return user, nil
}

func (a *Auth) validateEmailAddressForNewUser(email string) (string, bool) {

	if !(len(email) > 0) {
		return "email address required", false
	}

	match, err := regexp.MatchString(EmailValidationRegex, email)

	if err != nil {
		// rely on testing to pickup regexp parsing errors, the regex is static
		panic(err)
	}

	if !match {
		return "email address not valid", false
	}

	_, err = a.store.FindUserByEmail(email)

	if _, ok := err.(*respond.NotFoundError); !ok {
		return "email address taken", false
	}

	return "", true
}

func (a *Auth) validatePassword(password string) (string, bool) {

	if !(len(password) > 0) {
		return "password is required", false
	}

	return "", true
}

func (a *Auth) SignupUser(email, password string) (*User, error) {

	dataError := respond.NewDataErrorWithMessage("validation failed")

	if errorMessage, valid := a.validatePassword(password); !valid {
		dataError.Add("Password", errorMessage)
	}

	if errorMessage, valid := a.validateEmailAddressForNewUser(email); !valid {
		dataError.Add("Email", errorMessage)
	}

	if dataError.HasDetails() {
		return nil, dataError
	}

	return a.store.CreateUser(email, password)
}
