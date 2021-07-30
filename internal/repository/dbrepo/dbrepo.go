package dbrepo

import (
	"database/sql"

	"github.com/surakshith-suvarna/bookings/internal/config"
	"github.com/surakshith-suvarna/bookings/internal/repository"
)

type postgresDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

//Repo type for tests
type testDBRepo struct {
	App *config.AppConfig
	DB  *sql.DB
}

//NewPostgresRepo passing connection pool and app config returning a repository (This is DB repository). If we want another DB create a new type and new function newmysqlrepo
func NewPostgresRepo(conn *sql.DB, a *config.AppConfig) repository.DatabaseRepo {
	return &postgresDBRepo{
		App: a,
		DB:  conn,
	}
}

func NewTestingRepo(a *config.AppConfig) repository.DatabaseRepo {
	return &testDBRepo{
		App: a,
	}
}
