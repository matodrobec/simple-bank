package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "github.com/matodrobec/simplebank/db/sqlc"
	"github.com/matodrobec/simplebank/token"
	"github.com/matodrobec/simplebank/util"
)

type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot craete token maker: %v", err)
	}

	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	// router := gin.Default()

	// validators
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}

	server.setupRoutes()

	// // route
	// router.POST("/accounts", server.createAccount)
	// router.GET("/accounts/:id", server.getAccount)
	// router.DELETE("/accounts/:id", server.deleteAccount)
	// router.PUT("/accounts/:id", server.updateAccount)
	// router.GET("/accounts", server.listAccount)

	// router.POST("/transfers", server.createTransfer)

	// router.POST("/users", server.createUser)
	// router.POST("/users/login", server.loginUser)

	// server.router = router
	return server, nil
}

func (server *Server) setupRoutes() {
	server.router = gin.Default()

	// route
	server.router.POST("/users", server.createUser)
	server.router.POST("/users/login", server.loginUser)
	server.router.POST("/tokens/renew_access", server.renewToken)

	authRoutes := server.router.Group("/", authMiddleware(server.tokenMaker))

	authRoutes.POST("/accounts", server.createAccount)
	authRoutes.GET("/accounts/:id", server.getAccount)
	// authRoutes.DELETE("/accounts/:id", server.deleteAccount)
	authRoutes.PUT("/accounts/:id", server.updateAccount)
	authRoutes.GET("/accounts", server.listAccount)

	authRoutes.POST("/transfers", server.createTransfer)

}

func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
