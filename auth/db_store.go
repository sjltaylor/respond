package auth

import (
	"database/sql"
	"respond"
	"respond/crypto"
	"strings"
	"time"
)

const createStmt string = `

	CREATE TABLE users (
    "id"            SERIAL PRIMARY KEY,
    "email"         text NOT NULL,
    "password_hash" text NOT NULL,
    "password_salt" text NOT NULL,
    "created_at"    timestamp NOT NULL,
    "updated_at"    timestamp NOT NULL
  );

  CREATE UNIQUE INDEX idx_user_email ON users (email);
`

type DBStore struct {
	db                  *sql.DB
	insertUserStmt      *sql.Stmt
	findUserByEmailStmt *sql.Stmt
}

func NewDBStore(db *sql.DB) *DBStore {
	store := &DBStore{db: db}
	return store
}

func (store * DBStore) Create () (err error) {
	_, err = store.db.Exec(createStmt)
	return
}

func (store *DBStore) Reset () (err error) {
	
	if _, err = store.db.Exec(`DROP TABLE IF EXISTS users;`); err != nil {
		return
	}
	
	return store.Create()
}

func (store *DBStore) prepareOrPanic(query string) (stmt *sql.Stmt) {

	var err error

	if stmt, err = store.db.Prepare(query); err != nil {
		panic(err)
	}

	return
}

func (store *DBStore) PrepareOrPanic() {

	store.insertUserStmt = store.prepareOrPanic( `
		INSERT INTO "users" 
			("email", "password_salt", "password_hash", "created_at", "updated_at") 
		VALUES 
			($1, $2, $3, $4, $5) 
		RETURNING "id";
	`)

	store.findUserByEmailStmt = store.prepareOrPanic(`
		SELECT "id", "email", "password_salt", "password_hash", "created_at", "updated_at" FROM "users" WHERE "email" ILIKE $1;
	`)
}

func (store *DBStore) CreateUserInTx(tx *sql.Tx, email, password string) (user *User, err error) {

	createdAt := time.Now()

	salt, hash := crypto.Base64SaltAndHash(password, 64)

	var userId int64
	if err = tx.Stmt(store.insertUserStmt).QueryRow(email, salt, hash, createdAt, createdAt).Scan(&userId); err != nil {
		// and error "sql: no rows in result set" is returned when the insert fails
		// for example, because there is a unique index collision (e.g) email
		return
	}

	user = &User{
		Id:           userId,
		Email:        strings.ToLower(email),
		PasswordSalt: salt,
		PasswordHash: hash,
		CreatedAt:    createdAt,
		UpdatedAt:    createdAt,
	}

	return
}

func (store *DBStore) FindUserByEmail(email string) (user *User, err error) {

	user = &User{}

	if err = store.findUserByEmailStmt.QueryRow(email).Scan(
		&user.Id,
		&user.Email,
		&user.PasswordSalt,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	); err != nil {

		if err.Error() == "sql: no rows in result set" {
			return nil, respond.NewNotFoundError("user with email '%s'", email)
		}

		return nil, err
	}

	user.Email = strings.ToLower(user.Email)

	return
}

func (store *DBStore) FindUserByEmailAgainstPassword(email, password string) (user *User, err error) {

	if user, err = store.FindUserByEmail(email); err != nil {
		return
	}

	if user.PasswordHash != crypto.Base64HashOfSaltedPassword(password, user.PasswordSalt) {

		dataError := respond.NewDataError()
		dataError.Add("Password", "password mismatch")

		return nil, dataError
	}

	return user, nil
}


