package main

import (
	"database/sql"
	"github/bekeeeee/simplebank/api"
	db "github/bekeeeee/simplebank/db/sqlc"
	"github/bekeeeee/simplebank/util"
	"log"

	_ "github.com/lib/pq"
)


func main(){

	config,err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("Cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver,config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server:", err)
	}
}