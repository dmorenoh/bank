package app

import "github.com/gin-gonic/gin"

func (s *Server) Routes() *gin.Engine {
	router := s.router

	accV1 := router.Group("/v1/account")
	{
		accV1.POST("/", s.Create())
		accV1.PATCH("/:accountID/money", s.AddMoney())
		accV1.GET("/", s.GetAll())
		accV1.GET("/:accountID", s.Get())
	}

	transferV1 := router.Group("/v1/transfer")
	{
		transferV1.POST("/", s.Transfer())
	}

	return router
}
