package api

import (
    "github.com/gin-gonic/gin"
    db "github.com/ibnumei/digitalBankGo/db/sqlc"
)

// Server serves HTTP request for out banking service
type Server struct {
    store db.Store
    router *gin.Engine
}

//NewServer creates a new HTTP server and setup routing
func NewServer(store db.Store) *Server  {
    server := &Server{store: store}
    router := gin.Default()

    //List all routes

    router.POST("/accounts", server.createAccount)
    router.GET("/accounts/:id", server.getAccount)
    router.GET("/accounts", server.listAccount) //get list accounts with pagination


    router.POST("/transfers", server.createTransfer)

    server.router = router
    return server
}

//Start runs the http server on a specific address
func (server *Server) Start(address string) error {
    return server.router.Run(address)
}

//mapping error
func errorResponse(err error) gin.H {
    return gin.H{"error": err.Error()}
}