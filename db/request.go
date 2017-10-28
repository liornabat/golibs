package db

import (
	"context"
	_ "database/sql"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
)

type SqlRequest struct {
	sp   string
	dest interface{}
	db   *sqlx.DB
}

func (sc *SqlClient) NewRequest() *SqlRequest {
	sr := &SqlRequest{
		db: sc.db,
	}
	return sr
}

func (sr *SqlRequest) SetStoredProcedure(sp string) *SqlRequest {
	sr.sp = sp
	return sr
}

func (sr *SqlRequest) SetDestination(dest interface{}) *SqlRequest {
	sr.dest = dest
	return sr
}

func (sr *SqlRequest) Run(ctx context.Context) (err error) {
	err = sr.db.SelectContext(ctx, sr.dest, sr.sp)
	return
}

func (sr *SqlRequest) RunSingle(ctx context.Context) (err error) {
	err = sr.db.GetContext(ctx, sr.dest, sr.sp)
	return
}
