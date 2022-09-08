package services

import (
	"database/sql"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"istorage/logger"
)

func RunMigrations() {
	db, err := sql.Open("postgres", GetDataSourceName("sslmode=disable"))
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Fatal(err)
	}

	defer driver.Close()

	m, err := migrate.NewWithDatabaseInstance("file://migrations", os.Getenv("POSTGRES_DB"), driver)
	if err != nil {
		logger.Fatal(err)
	}

	defer m.Close()

	m.Up()
}
