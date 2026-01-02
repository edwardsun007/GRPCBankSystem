package api

import (
	db "github.com/techschool/simple-bank/db2/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
    store *db.Store   // pointer to the store object that contains the database connection, directly access it without copying the value
	router *gin.Engine // HTTP request router
}

// Constructor to create a new server instance, return a pointer to that instance
func NewServer(store *db.Store) *Server {
	server := &Server{store: store}
	router := gin.Default() // this use the gin default middleware


	// add routes to the router
	router.POST("/accounts", server.createAccount)
	// router.GET("/accounts/:id", server.getAccountById)
	// router.GET("/accounts/:id", server.getAccount)
	// router.GET("/accounts", server.listAccounts)
	// router.PUT("/accounts/:id", server.updateAccount)
	// router.DELETE("/accounts/:id", server.deleteAccount)
	// router.POST("/transfers", server.createTransfer)
	// router.GET("/transfers/:id", server.getTransfer)
	// router.GET("/transfers", server.listTransfers)
	server.router = router // assign the router to the server instance

	return server
}

// helper function to return error response in JSON format
func errorResponse(err error) gin.H {
	return gin.H{"error": err.Error()}
}

// Start runs the HTTP server on a specific address.
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}