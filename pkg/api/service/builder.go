package service

import (
	"bank/pkg/api/dto"
	"bank/pkg/api/model"
	"github.com/google/uuid"
)

func NewAccount(req dto.CreateAccountRequest) (*model.Account, error) {
	return &model.Account{
		ID:     uuid.New(),
		Name:   req.Name,
		Amount: req.Amount,
	}, nil
}
