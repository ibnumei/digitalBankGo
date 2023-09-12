package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/ibnumei/digitalBankGo/db/sqlc"
	"github.com/ibnumei/digitalBankGo/token"
	"github.com/ibnumei/digitalBankGo/util"
)

// Server serves HTTP request for out banking service
type Server struct {
	config     util.Config
	store      db.Store
	tokenMaker token.Maker
	router     *gin.Engine
}

//NewServer creates a new HTTP server and setup routing
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server.setupRouter()
	return server, nil
}

func (server *Server) setupRouter() {
	router := gin.Default()

	//List all routes
	router.POST("/users", server.createUser)
	router.POST("/users/login", server.loginUser)

	authRoutes := router.Group("/").Use(authMiddleware(server.tokenMaker))

	//Rule: A logged-in user can only create an account for him/herself
	authRoutes.POST("/accounts", server.createAccount)
	//Rule: A logged-in user can only get accounts that he/she owns
	authRoutes.GET("/accounts/:id", server.getAccount)
	//Rule: A logged-in user can only list accounts that belong to him/her
	authRoutes.GET("/accounts", server.listAccount) //get list accounts with pagination
	//Rule: A logged-in user can only send money from his/her own account
	authRoutes.POST("/transfers", server.createTransfer)

	server.router = router
}

//Start runs the http server on a specific address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

//mapping error
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}
