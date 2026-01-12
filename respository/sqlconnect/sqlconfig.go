package sqlconnect

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"http-api.com/utils"
)

// ConnectDB establishes a connection to MariaDB using environment config
func ConnectDB() (*sql.DB, error) {
	config := utils.LoadDBConfig()

	// Build DSN (Data Source Name)
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	fmt.Printf("✓ Database connected successfully to %s:%s/%s\n",
		config.Host, config.Port, config.Database)

	return db, nil
}
