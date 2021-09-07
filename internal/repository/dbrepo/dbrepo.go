package dbrepo

import (
	"database/sql"

	"github.com/prayagsingh/bookings/internal/config"
	"github.com/prayagsingh/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// for testcases
type testPostgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

// NewPostgresRepo creates a new instance of PostgresDbRepo
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {

	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

// NewTestPostgresRepo creates a new instance of PostgressDbRepo for testcases
func NewTestPostgresRepo(a *config.AppConfig) repository.DatabaseRepo {

	return &testPostgresDBRepo{
		App: a,
	}
}
