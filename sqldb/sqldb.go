package sqldb

import "database/sql"

func PrepareOrPanic(db *sql.DB, query string) (stmt *sql.Stmt) {

	var err error

	if stmt, err = db.Prepare(query); err != nil {
		panic(err)
	}

	return stmt
}
