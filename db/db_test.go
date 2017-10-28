package db

import (
	"github.com/stretchr/testify/require"
	"testing"

	"context"
	"os"
)

var (
	msSqlConnection = &SqlConnection{
		Host:     "10.20.30.21",
		Port:     1421,
		User:     "Reports",
		Password: "reports",
		Database: "cheetah",
		Timeout:  30}
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestSqlClient_Ping(t *testing.T) {
	require := require.New(t)
	client, err := NewSqlClient(MsSql, msSqlConnection)
	require.NoError(err)
	err = client.Ping()
	require.NoError(err)
}

func TestSqlRequest_Run(t *testing.T) {
	require := require.New(t)
	client, err := NewSqlClient(MsSql, msSqlConnection)
	require.NoError(err)

	type rate struct {
		InstrumentID   int     `db:"InstrumentID"`
		InstrumentName string  `db:"InstrumentName"`
		Symbol         string  `db:"Symbol"`
		Bid            float64 `db:"Bid"`
		Ask            float64 `db:"Ask"`
		Coif           float64 `db:"Coif"`
		PipsValue      float64 `db:"PipsValue"`
		Spread         float64 `db:"Spread"`
		Dir            int     `db:"Dir"`
	}
	rates := &[]*rate{}
	ctx := context.Background()
	err = client.NewRequest().
		SetStoredProcedure("exec dbo.p_GLB_GetRates 1").
		SetDestination(rates).
		Run(ctx)
	require.NoError(err)
	require.Equal(33, len(*rates))
}
