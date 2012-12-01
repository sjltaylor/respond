package auth

import (
	"database/sql"
	"respond"
	"respond/crypto"
	"respond/sqldb"
	"strings"
	"time"
)

type DBStore struct {
	db                  *sql.DB
	insertUserStmt      *sql.Stmt
	findUserByEmailStmt *sql.Stmt
}

func NewDBStore(db *sql.DB) *DBStore {
	store := &DBStore{db: db}
	store.prepareUserStmts(db)
	return store
}

func (store *DBStore) prepareUserStmts(db *sql.DB) {

	store.insertUserStmt = sqldb.PrepareOrPanic(db, `
		INSERT INTO "users" 
			("email", "password_salt", "password_hash", "created_at", "updated_at") 
		VALUES 
			($1, $2, $3, $4, $5) 
		RETURNING "id";
	`)

	store.findUserByEmailStmt = sqldb.PrepareOrPanic(db, `
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

func (store *DBStore) CreateUser(email, password string) (user *User, err error) {

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
