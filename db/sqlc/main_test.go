package db

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/matodrobec/simplebank/util"
)

// const (
// 	dbDriver = "postgres"
// 	dbSoruce = "postgresql://postgres:test@localhost:5432/bank?sslmode=disable"
// )

// var testQueries *Queries
// var testDB *sql.DB

var testStore Store

func TestMain(m *testing.M) {
	// var err error

	// config, err := util.LoadConfig("./../..")
	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot load config:", err)
	}

	// testDB, err = sql.Open(dbDriver, dbSoruce)
	connPool, err := pgxpool.New(context.Background(), config.DBSource)
	if err != nil {
		log.Fatal("cannot conntect to db: ", err)
	}

	// testQueries = New(testDB)
	testStore =  NewStore(connPool)

	os.Exit(m.Run())
}
