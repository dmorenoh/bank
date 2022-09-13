package app

import (
	"bank/pkg/api/dto"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"net/http"
)

func (s *Server) Create() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.CreateAccountRequest
		bindErr := ctx.ShouldBindJSON(&req)
		if bindErr != nil {
			ctx.IndentedJSON(http.StatusBadRequest,
				gin.H{"error": "fields validation failed",
					"desc": fmt.Sprintf("deails %v", bindErr),
				})
			return
		}
		resp, cErr := s.accountService.Create(req)
		if cErr != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx.IndentedJSON(http.StatusCreated, resp)
		return
	}
}

func (s *Server) AddMoney() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountIDPar := ctx.Param("accountID")
		accID, pErr := uuid.Parse(accountIDPar)
		if pErr != nil {
			ctx.IndentedJSON(http.StatusBadRequest, pErr)
			return
		}
		var req dto.UpdateAccountRequest
		bindErr := ctx.ShouldBindJSON(&req)
		if bindErr != nil {
			ctx.IndentedJSON(http.StatusBadRequest, bindErr)
			return
		}
		resp, cErr := s.accountService.AddMoney(accID, req.Amount)
		if cErr != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx.IndentedJSON(http.StatusAccepted, resp)
		return
	}
}

func (s *Server) Transfer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req dto.TransferenceRequest
		bindErr := ctx.ShouldBindJSON(&req)
		if bindErr != nil {
			ctx.IndentedJSON(http.StatusBadRequest, bindErr)
			return
		}
		cErr := s.accountService.Transfer(req.From, req.To, req.Amount)
		if cErr != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx.IndentedJSON(http.StatusAccepted, gin.H{})
		return
	}
}

func (s *Server) Get() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		accountIDPar := ctx.Param("accountID")
		accID, pErr := uuid.Parse(accountIDPar)
		if pErr != nil {
			ctx.IndentedJSON(http.StatusBadRequest, pErr)
			return
		}
		resp, cErr := s.accountService.Get(accID)
		if cErr != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx.IndentedJSON(http.StatusAccepted, resp)
		return
	}
}

func (s *Server) GetAll() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		resp, cErr := s.accountService.GetAll()
		if cErr != nil {
			ctx.AbortWithStatus(http.StatusInternalServerError)
		}
		ctx.IndentedJSON(http.StatusAccepted, resp)
		return
	}
}
