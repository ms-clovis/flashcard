package session

import "database/sql"

type MongoDB struct {
	db *sql.DB
}
