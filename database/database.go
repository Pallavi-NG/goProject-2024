package database

import (
	"context"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Database struct {
	Client *sqlx.DB
}

//without .env file
/* func NewDatabase() (*Database, error) {
	// Opening a database connection.
	db, err := sqlx.Open("mysql", "root:mindPalace@23@tcp(localhost:3308)/students?parseTime=true")
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!")
	return &Database{
		Client: db,
	}, nil
} */

// using .env file
func NewDatabase() (*Database, error) {

	LoadEnvVariables()
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true",
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	db, err := sqlx.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	fmt.Println("Connected!")
	return &Database{
		Client: db,
	}, nil
}

func (d *Database) Ping(ctx context.Context) error {
	return d.Client.DB.PingContext(ctx)
}
