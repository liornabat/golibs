package db

import "database/sql"

type SqlResult struct {
	IsExecuted   bool
	Error        error
	LastInsertId int64
	RowsAffected int64
}

func GetResult(res sql.Result, err error) *SqlResult {
	sr := &SqlResult{
		IsExecuted: true,
		Error:      err,
	}
	sr.LastInsertId, _ = res.LastInsertId()
	sr.RowsAffected, _ = res.RowsAffected()
	if err != nil {
		sr.IsExecuted = false
	}
	return sr
}
