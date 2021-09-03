package driver

import (
	"database/sql"
	"time"

	// driver to connect to db. don't remove it
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgconn"
	_ "github.com/jackc/pgx/v4/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbConn = &DB{}

const maxOpenDbConn = 10
const maxIdleDbConn = 5
const maxDbLifeTime = 5 * time.Minute

// ConnectSQL create a db connection-pool for postgres
func ConnectSQL(dsn string) (*DB, error) {

	db, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}

	// setting parameters to db
	db.SetMaxOpenConns(maxOpenDbConn)
	db.SetMaxIdleConns(maxIdleDbConn)
	db.SetConnMaxLifetime(maxDbLifeTime)

	dbConn.SQL = db

	err = testDb(db)
	if err != nil {
		return nil, err
	}
	return dbConn, nil
}

// tries to ping the db
func testDb(d *sql.DB) error {

	err := d.Ping()
	if err != nil {
		return err
	}

	return nil
}

// NewDatabase creates a new db for the application
func NewDatabase(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
