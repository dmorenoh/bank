package dto

import "github.com/google/uuid"

type CreateAccountRequest struct {
	Name   string  `json:"name" binding:"required,min=3"`
	Amount float64 `json:"amount" binding:"required"`
}

type CreateAccountResponse struct {
	ID     uuid.UUID
	Name   string
	Amount float64
}
type UpdateAccountRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}
type UpdateAccountResponse struct {
	ID            uuid.UUID
	Name          string
	CurrentAmount float64
}

type TransferenceRequest struct {
	From   uuid.UUID `json:"from" binding:"required"`
	To     uuid.UUID `json:"to" binding:"required"`
	Amount float64   `json:"amount" binding:"required"`
}

type GetAccountResponse struct {
	ID     uuid.UUID
	Name   string
	Amount float64
}

type GetAllAccountResponse struct {
	Accounts []GetAccountResponse
}
