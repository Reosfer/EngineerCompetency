package pgsql

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
)

type SqlInterface interface {
	DB() *sql.DB
}

type sqlStruct struct {
	db *sql.DB
}

func InitSql(driver string, host string, port string, username string, password string, database string) (SqlInterface, error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", username, password, host, port, database)
	db, err := sql.Open(driver, connectionString)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return &sqlStruct{
		db: db,
	}, nil

}

func InitPgSql(url string, driver string, host string, port string, username string, password string, database string) (SqlInterface, error) {
	connectionString := ""
	if len(url) > 0 {
		connectionString = url
	} else {
		connectionString = fmt.Sprintf("postgresql://%s:%s@%s:%s/%s?sslmode=disable", username, password, host, port, database)
	}

	db, err := sql.Open(driver, connectionString)

	if err != nil {
		fmt.Println("error cn")
		return nil, err
	}

	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(50)

	return &sqlStruct{
		db: db,
	}, nil

}

func (m *sqlStruct) DB() *sql.DB {
	return m.db
}
