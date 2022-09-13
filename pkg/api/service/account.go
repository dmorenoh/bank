package service

import (
	"bank/pkg/api/dto"
	"bank/pkg/api/repositories"
	"errors"
	"github.com/google/uuid"
	"sync"
)

type AccountService interface {
	Create(req dto.CreateAccountRequest) (dto.CreateAccountResponse, error)
	AddMoney(accountID uuid.UUID, amount float64) (dto.UpdateAccountResponse, error)
	Transfer(accFrom uuid.UUID, accTo uuid.UUID, amount float64) error
	Get(accountID uuid.UUID) (dto.GetAccountResponse, error)
	GetAll() (dto.GetAllAccountResponse, error)
}

func NewAccountService(repository repositories.AccountRepository) AccountService {
	return &accountService{
		repository: repository,
	}
}

type accountService struct {
	mux        sync.RWMutex
	repository repositories.AccountRepository
}

func (a *accountService) Get(accountID uuid.UUID) (dto.GetAccountResponse, error) {
	account, err := a.repository.Get(accountID)
	if err != nil {
		return dto.GetAccountResponse{}, err
	}

	return dto.GetAccountResponse{
		ID:     account.ID,
		Name:   account.Name,
		Amount: account.Amount,
	}, nil
}

func (a *accountService) GetAll() (dto.GetAllAccountResponse, error) {
	accounts, err := a.repository.GetAll()
	if err != nil {
		return dto.GetAllAccountResponse{}, err
	}

	accResponses := make([]dto.GetAccountResponse, len(accounts))
	for i := range accounts {
		account := accounts[i]
		accResponses[i] = dto.GetAccountResponse{
			ID:     account.ID,
			Name:   account.Name,
			Amount: account.Amount,
		}
	}

	return dto.GetAllAccountResponse{Accounts: accResponses}, nil
}

func (a *accountService) Create(req dto.CreateAccountRequest) (dto.CreateAccountResponse, error) {
	newAccount, nErr := NewAccount(req)
	if nErr != nil {
		return dto.CreateAccountResponse{}, nErr
	}

	if cErr := a.repository.Create(newAccount); cErr != nil {
		return dto.CreateAccountResponse{}, cErr
	}

	return dto.CreateAccountResponse{
		ID:     newAccount.ID,
		Name:   newAccount.Name,
		Amount: newAccount.Amount,
	}, nil
}

func (a *accountService) AddMoney(accountID uuid.UUID, amount float64) (dto.UpdateAccountResponse, error) {
	// used for pessimistic locking
	a.mux.Lock()
	defer a.mux.Unlock()

	acc, gErr := a.repository.Get(accountID)
	if gErr != nil {
		return dto.UpdateAccountResponse{}, gErr
	}

	if addErr := acc.AddMoney(amount); addErr != nil {
		return dto.UpdateAccountResponse{}, addErr
	}

	if updErr := a.repository.Update(acc); updErr != nil {
		return dto.UpdateAccountResponse{}, updErr
	}

	return dto.UpdateAccountResponse{
		ID:            acc.ID,
		Name:          acc.Name,
		CurrentAmount: acc.Amount,
	}, nil
}

func (a *accountService) Transfer(fromID uuid.UUID, toID uuid.UUID, amount float64) error {
	// used for pessimistic locking
	a.mux.Lock()
	defer a.mux.Unlock()

	if fromID == toID || amount <= 0 {
		return errors.New("inconsistent data")
	}

	from, fErr := a.repository.Get(fromID)
	if fErr != nil {
		return fErr
	}

	to, tErr := a.repository.Get(toID)
	if tErr != nil {
		return tErr
	}

	if wErr := from.Withdraw(amount); wErr != nil {
		return wErr
	}

	if aErr := to.AddMoney(amount); aErr != nil {
		return aErr
	}

	// update both changes in both accounts as transactional
	if err := a.repository.UpdatesTx(from, to); err != nil {
		return err
	}

	return nil
}
