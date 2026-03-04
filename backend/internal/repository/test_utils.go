package repository

import (
	"real-time-forum/internal/database"
	"testing"
)

func SetupTestDB(t *testing.T) *Repositories {
	t.Helper()

	db, err := database.NewDatabase(":memory:")
	if err != nil {
		panic("Failed to connect to test database: " + err.Error())
	}

	if err := database.RunMigrations(db); err != nil {
		panic("Failed to run migrations: " + err.Error())
	}

	t.Cleanup(func() {
		db.Close()
	})

	return NewRepositories(db)
}
