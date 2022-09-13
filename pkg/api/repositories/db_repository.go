package repositories

import (
	"bank/pkg/api/model"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AccountEntity struct {
	ID     uuid.UUID `gorm:"column:id;PRIMARY_KEY"`
	Name   string
	Amount float64
}

type dbRepository struct {
	db *gorm.DB
}

func (d *dbRepository) GetAll() ([]*model.Account, error) {
	var accountsEnt []AccountEntity
	if err := d.db.Find(&accountsEnt).Error; err != nil {
		return nil, err
	}

	accounts := make([]*model.Account, len(accountsEnt))
	for i := range accountsEnt {
		accounts[i] = &model.Account{
			ID:     accountsEnt[i].ID,
			Name:   accountsEnt[i].Name,
			Amount: accountsEnt[i].Amount,
		}
	}

	return accounts, nil
}

func NewDBRepository(db *gorm.DB) AccountRepository {
	return &dbRepository{
		db: db,
	}
}

func (d *dbRepository) Create(account *model.Account) error {
	entity := AccountEntity{
		ID:     account.ID,
		Name:   account.Name,
		Amount: account.Amount,
	}

	if cErr := d.db.Create(&entity).Error; cErr != nil {
		return cErr
	}

	return nil
}

func (d dbRepository) Update(account *model.Account) error {
	accEnt := AccountEntity{
		ID:     account.ID,
		Name:   account.Name,
		Amount: account.Amount,
	}

	if err := d.db.Save(&accEnt).Error; err != nil {
		return err
	}

	return nil
}

func (d dbRepository) Get(accountID uuid.UUID) (*model.Account, error) {
	var accEnt AccountEntity

	if err := d.db.First(&accEnt, accountID).Error; err != nil {
		return nil, err
	}

	return &model.Account{
		ID:     accEnt.ID,
		Name:   accEnt.Name,
		Amount: accEnt.Amount,
	}, nil
}

func (d *dbRepository) UpdatesTx(accounts ...*model.Account) error {
	txErr := d.db.Transaction(func(tx *gorm.DB) error {
		for _, acc := range accounts {

			ent := AccountEntity{
				ID:     acc.ID,
				Name:   acc.Name,
				Amount: acc.Amount,
			}

			if saveToErr := tx.Save(&ent).Error; saveToErr != nil {
				return saveToErr
			}
		}
		return nil
	})

	return txErr
}
