package auth

import (
	"database/sql"
	_ "github.com/bmizerany/pq"
	"respond"
	"strings"
	"testing"
	"time"
)

var store *DBStore

func init() {
	store = dbstore()
}

func dbstore() *DBStore {

	if store == nil {

		if db, err := sql.Open("postgres", "dbname=respond_auth_testing sslmode=disable"); err != nil {
			panic(err)
		} else {

			store = NewDBStore(db)
			store.Reset()
			store.PrepareOrPanic()
		}
	}

	return store
}

func isNotFoundError(err error) bool {
	_, ok := err.(*respond.NotFoundError)
	return ok
}

func (store *DBStore) createUser(email, password string) (user *User, err error) {

	var tx *sql.Tx

	if tx, err = store.db.Begin(); err != nil {
		return
	}

	if user, err = store.CreateUserInTx(tx, email, password); err != nil {
		return
	}

	err = tx.Commit()

	return
}

func TestCreateUser(t *testing.T) {

	testStartNano := time.Now().UnixNano()

	var countBefore int64
	var countAfter int64

	if err := store.db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&countBefore); err != nil {
		t.Fatal(err)
	}

	var user *User
	email := emailGenerator.Next()

	if u, err := store.createUser(email, "password"); err != nil {
		t.Fatal(err)
	} else {
		user = u
	}

	if err := store.db.QueryRow(`SELECT COUNT(*) FROM users;`).Scan(&countAfter); err != nil {
		t.Fatal(err)
	}

	if countAfter != (countBefore + 1) {
		t.Fatal("The user count did not increment.")
	}

	if user.Email != strings.ToLower(email) {
		t.Fatalf("email not set on returned user (%s) or not converted to lowercase", user.Email)
	}

	if user.PasswordSalt == "" {
		t.Fatal("password salt not set on returned user")
	}

	if user.PasswordHash == "" {
		t.Fatal("password hash not set on returned user")
	}

	if user.CreatedAt.UnixNano() < testStartNano {
		t.Fatal("created at not set to now on returned user")
	}

	if user.UpdatedAt.UnixNano() < testStartNano {
		t.Fatal("updated at not set to now on returned user")
	}

	if user.UpdatedAt.UnixNano() != user.CreatedAt.UnixNano() {
		t.Fatal("created at and updated at times should be the same on returned user")
	}

	if user.Id == 0 {
		t.Fatal("user id not set")
	}
}

func TestFindUser(t *testing.T) {

	email := strings.ToUpper(emailGenerator.Next())

	if _, err := store.createUser(email, "password"); err != nil {
		t.Fatal("error creating user: ", err)
	}

	var user *User

	if u, err := store.FindUserByEmail(email); err != nil {
		t.Fatal(err)
	} else {
		user = u
	}

	if _, err := store.FindUserByEmail(strings.ToLower(email)); isNotFoundError(err) {
		t.Fatal("not case insensitive")
	}

	if user.Email != strings.ToLower(email) {
		t.Fatalf("email not set on returned user (%s) or not converted to lowercase", user.Email)
	}

	if user.PasswordSalt == "" {
		t.Fatal("password salt not set on returned user")
	}

	if user.PasswordHash == "" {
		t.Fatal("password hash not set on returned user")
	}

	if user.CreatedAt == *new(time.Time) {
		t.Fatal("created at not set to now on returned user")
	}

	if user.UpdatedAt == *new(time.Time) {
		t.Fatal("updated at not set to now on returned user")
	}

	if user.Id == 0 {
		t.Fatal("user id not set")
	}
}

func TestFindUserWithInvalidEmail(t *testing.T) {

	email := emailGenerator.Next()

	if _, err := store.createUser(email, "password"); err != nil {
		t.Fatal("error creating user: ", err)
	}

	if _, err := store.FindUserByEmail("iamnotthere@all.com"); !isNotFoundError(err) {
		t.Fatal("does not return an NotFoundError when a non-existant email is passed")
	}
}

func TestFindUserAgainstPassword(t *testing.T) {

	email := emailGenerator.Next()

	if _, err := store.createUser(email, "very-secret"); err != nil {
		t.Fatal(err)
	}

	user, err := store.FindUserByEmailAgainstPassword(email, "very-secret")

	if err != nil {
		t.Fatal(err)
	}

	if user == nil {
		t.Fatal("user not returned")
	}

	user, err = store.FindUserByEmailAgainstPassword(email, "WRONG!!PASSWORD")

	if dataError, ok := err.(*respond.DataError); ok {
		if dataError.Errors["Password"][0] != "password mismatch" {
			t.Fatal("password mismatch error not returned")
		}
	} else {
		t.Fatal("returned error was not a *respond.DataError")
	}

	if user != nil {
		t.Fatal("user returned against invalid password")
	}
}