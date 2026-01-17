package sqlconnect

import (
	"database/sql"
	"fmt"
	"sync"

	_ "github.com/go-sql-driver/mysql"
	"http-api.com/utils"
)

var (
	db   *sql.DB
	once sync.Once
)

// GetDB returns the shared database connection (initializes on first call)
func GetDB() (*sql.DB, error) {
	var initErr error
	once.Do(func() {
		config := utils.LoadDBConfig()
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
			config.User, config.Password, config.Host, config.Port, config.Database)

		var err error
		db, err = sql.Open("mysql", dsn)
		if err != nil {
			initErr = fmt.Errorf("failed to open database: %w", err)
			return
		}

		db.SetMaxOpenConns(config.MaxOpenConns)
		db.SetMaxIdleConns(config.MaxIdleConns)
		db.SetConnMaxLifetime(config.ConnMaxLifetime)

		if err := db.Ping(); err != nil {
			initErr = fmt.Errorf("failed to ping database: %w", err)
			return
		}
		fmt.Printf("✓ Database connected to %s:%s/%s\n", config.Host, config.Port, config.Database)
	})
	return db, initErr
}

// ConnectDB is kept for backward compatibility - returns shared connection
func ConnectDB() (*sql.DB, error) {
	return GetDB()
}
