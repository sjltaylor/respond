package auth

import `database/sql`

const resetStmt string = `

	DROP TABLE IF EXISTS users;

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

func RecreateUserStore(db *sql.DB) (err error) {

	if _, err = db.Exec(resetStmt); err != nil {
		return
	}

	return
}
