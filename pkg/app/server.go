package app

import (
	"bank/pkg/api/service"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	accountService service.AccountService
	router         *gin.Engine
}

func NewServer(router *gin.Engine, service service.AccountService) *Server {
	return &Server{
		router:         router,
		accountService: service,
	}
}

func (s *Server) Run() error {
	r := s.Routes()

	err := r.Run()
	if err != nil {
		log.Printf("error calling Run on router: %v", err)
		return err
	}

	return nil
}
