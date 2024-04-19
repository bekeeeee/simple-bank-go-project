package api

import (
	db "github/bekeeeee/simplebank/db/sqlc"

	"github.com/gin-gonic/gin"
)
type Server struct{
	store *db.Store
	router *gin.Engine
}


func NewServer(store *db.Store)*Server{
	

	router := gin.Default();
	server := &Server{
		store: store,
		router: router,
	}
	router.POST("/accounts",server.createAccount)
	router.GET("/accounts/:id",server.getAccount)
	router.GET("/accounts",server.listAccounts)



	return server

}


func (server *Server) Start(address string)error{
	return server.router.Run(address)
}

func errorResponse(err error) gin.H{
	return gin.H{"error": err.Error()}
}