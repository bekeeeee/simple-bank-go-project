package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

// import _ "github.com/lib/pq" // Import the latest installed version
var testQuery *Queries

func TestMain(m *testing.M){
const dbDriver = "postgres"
const dbSource = "postgresql://root:secret@localhost:5455/simple_bank?sslmode=disable"

	fmt.Println("hereeeeeeeeee")
	conn,err := sql.Open(dbDriver,dbSource)
	if err != nil{
		log.Fatal("Cannot connect to database: ",err)
	}
	testQuery = New(conn)
	os.Exit(m.Run())
}