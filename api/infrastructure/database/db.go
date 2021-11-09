package database

import (
	"app/api/constants"
	"app/api/infrastructure/nosql"
	"app/api/llog"
	"database/sql"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
)

type SQLHandler interface {
	Exec(query string, args ...interface{}) (SQLResult, error)
	Query(query string, args ...interface{}) (SQLRows, error)
	QueryRow(query string, args ...interface{}) SQLRow
}

type sqlHandler struct {
	DB *sql.DB
}

func New() (SQLHandler, error) {
	nosql.New()
	dbuser := os.Getenv("MYSQL_USER")
	if dbuser == "" {
		dbuser = constants.DBUser
	}
	dbpassword := os.Getenv("MYSQL_PASSWORD")
	if dbpassword == "" {
		dbpassword = constants.DBPassword
	}
	dbhost := os.Getenv("MYSQL_HOST")
	if dbhost == "" {
		dbhost = constants.DBHost
	}
	dbport := os.Getenv("MYSQL_PORT")
	if dbport == "" {
		dbport = constants.DBPort
	}
	dbname := os.Getenv("MYSQL_DATABASE")
	if dbname == "" {
		dbname = constants.DBName
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbuser, dbpassword, dbhost, dbport, dbname))
	if err != nil {
		return nil, errors.Wrap(err, "failed to connect db")
	}

	llog.Info(fmt.Sprintf("connect db: user=%s, target=%s:%s, dbName=%s", dbuser, dbhost, dbport, dbname))
	return &sqlHandler{DB: db}, nil
}

type SQLResult interface {
}

type sqlResult struct {
	Result sql.Result
}

type SQLRows interface {
	Scan(dest ...interface{}) error
	Next() bool
	CheckNoRows(err error) bool
}

type sqlRows struct {
	Rows *sql.Rows
}

type SQLRow interface {
	Scan(dest ...interface{}) error
	CheckNoRows(err error) bool
}

type sqlRow struct {
	Row *sql.Row
}

func (sh *sqlHandler) Exec(query string, args ...interface{}) (SQLResult, error) {
	res, err := sh.DB.Exec(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to execute query. query=%s", query))
	}
	return &sqlResult{
		Result: res,
	}, nil
}

func (sh *sqlHandler) Query(query string, args ...interface{}) (SQLRows, error) {
	res, err := sh.DB.Query(query, args...)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to execute query. query=%s", query))
	}
	return &sqlRows{
		Rows: res,
	}, nil
}

func (sh *sqlHandler) QueryRow(query string, args ...interface{}) SQLRow {
	res := sh.DB.QueryRow(query, args...)
	return &sqlRow{
		Row: res,
	}
}

func (r *sqlRows) Scan(dest ...interface{}) error {
	return r.Rows.Scan(dest...)
}

func (r *sqlRows) Next() bool {
	return r.Rows.Next()
}

func (r *sqlRows) CheckNoRows(err error) bool {
	return err == sql.ErrNoRows
}

func (r *sqlRow) Scan(dest ...interface{}) error {
	return r.Row.Scan(dest...)
}

func (r *sqlRow) CheckNoRows(err error) bool {
	return err == sql.ErrNoRows
}
