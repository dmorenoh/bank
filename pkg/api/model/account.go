package model

import (
	"errors"
	"github.com/google/uuid"
)

type Account struct {
	ID     uuid.UUID
	Name   string
	Amount float64
}

func (a *Account) AddMoney(amount float64) error {
	if amount < 0 {
		return errors.New("not valid amount")
	}
	a.Amount = a.Amount + amount
	return nil
}

func (a *Account) Withdraw(amount float64) error {
	if a.Amount < amount {
		return errors.New("not enough balance")
	}
	a.Amount = a.Amount - amount
	return nil
}
