package infrastructure

import (
	"database/sql"
	"fmt"

	"github.com/nstoker/bookish-palm-domain/src/interfaces"
	"github.com/sirupsen/logrus"

	// To get the sqlite3 goodness
	_ "github.com/mattn/go-sqlite3"
)

// SqliteHandler structure
type SqliteHandler struct {
	Conn *sql.DB
}

// Execute statement
func (handler *SqliteHandler) Execute(statement string) {
	handler.Conn.Exec(statement)
}

// Query expression
func (handler *SqliteHandler) Query(statement string) interfaces.Row {
	logrus.Info("sqliteHandler Query('%s')", statement)
	rows, err := handler.Conn.Query(statement)
	if err != nil {
		fmt.Println(err)
		return new(SqliteRow)
	}
	row := new(SqliteRow)
	row.Rows = rows
	return row
}

// SqliteRow structure
type SqliteRow struct {
	Rows *sql.Rows
}

// Scan a row
func (r SqliteRow) Scan(dest ...interface{}) {
	r.Rows.Scan(dest...)
}

// Next record
func (r SqliteRow) Next() bool {
	return r.Rows.Next()
}

// NewSqliteHandler get a new sqlite handler
func NewSqliteHandler(dbfilename string) *SqliteHandler {
	conn, _ := sql.Open("sqlite3", dbfilename)
	sqliteHandler := new(SqliteHandler)
	sqliteHandler.Conn = conn

	return sqliteHandler
}
