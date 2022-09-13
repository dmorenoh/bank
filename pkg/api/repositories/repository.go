package repositories

import (
	"bank/pkg/api/model"
	"github.com/google/uuid"
)

type AccountRepository interface {
	// Create account
	Create(account *model.Account) error
	// Update account
	Update(account *model.Account) error
	// Get account
	Get(accountID uuid.UUID) (*model.Account, error)
	// GetAll accounts
	GetAll() ([]*model.Account, error)
	// UpdatesTx updates a list of accounts in a transaction
	UpdatesTx(account ...*model.Account) error
}
