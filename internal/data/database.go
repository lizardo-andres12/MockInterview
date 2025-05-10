package data

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

func GenerateConfig() (*mysql.Config, error) {
	user := os.Getenv("DBUser")
	pwd := os.Getenv("DBPassword")
	net := "tcp"
	addr := "localhost:3306"
	name := "test"

	if user == "" {
		return nil, fmt.Errorf("Couldn't load database username...")
	}
	if pwd == "" {
		return nil, fmt.Errorf("Couldn't load database password...")
	}

	return &mysql.Config{
		User: user,
		Passwd: pwd,
		Net: net,
		Addr: addr,
		DBName: name,
	}, nil
}

func NewDatabase(cfg *mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	// Set up connections.

	return db, nil
}
