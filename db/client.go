package db

import (
	_ "database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type SqlClientType string

const (
	MsSql    SqlClientType = "mssql"
	MySql    SqlClientType = "mysql"
	Postgres SqlClientType = "postgres"
	SqlLite3 SqlClientType = "sqlite3"
)

type SqlConnection struct {
	Host         string `json:"host"`
	Port         int    `json:"port"`
	User         string `json:"user"`
	Password     string `json:"password"`
	Database     string `json:"database"`
	Timeout      int    `json:"timeout"`
	FileName     string `json:"file_name"`
	InMemory     bool   `json:"in_memory"`
	MaxIdleConns int
	MaxOpenConns int
}

type SqlClient struct {
	driver SqlClientType
	conn   *SqlConnection
	db     *sqlx.DB
}

func (scn *SqlConnection) ConnectionString(driver SqlClientType) string {
	switch driver {
	case MsSql:
		return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s&encrypt=disable&connection+timeout=%d", scn.User, scn.Password, scn.Host, scn.Port, scn.Database, scn.Timeout)

	case MySql:

	case Postgres:

	}
	return ""
}

func NewSqlClient(driver SqlClientType, conn *SqlConnection) (client *SqlClient, err error) {
	client = &SqlClient{
		driver: driver,
		conn:   conn,
	}
	client.db, err = sqlx.Open(string(driver), conn.ConnectionString(driver))
	if conn.MaxIdleConns != 0 {
		client.db.SetMaxIdleConns(conn.MaxIdleConns)
	}
	if conn.MaxOpenConns != 0 {
		client.db.SetMaxOpenConns(conn.MaxOpenConns)
	}

	return
}
func (sc *SqlClient) Ping() error {
	return sc.db.Ping()
}
