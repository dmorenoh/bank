package service_test

import (
	"bank/pkg/api/dto"
	"bank/pkg/api/repositories"
	"bank/pkg/api/service"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"sync"
	"testing"
)

func TestAccountService_Create(t *testing.T) {
	db := setup(t)

	t.Run("Given a create service", func(t *testing.T) {
		dbRepo := repositories.NewDBRepository(db)
		accService := service.NewAccountService(dbRepo)

		t.Run("When request to create account", func(t *testing.T) {
			req := dto.CreateAccountRequest{Name: "bill smith", Amount: 100.00}
			resp, err := accService.Create(req)
			require.NoError(t, err)

			t.Run("Then creates that new account", func(t *testing.T) {
				assert.Equal(t, resp.Name, req.Name)
				assert.Equal(t, resp.Amount, req.Amount)
				assert.NotNil(t, resp.ID)
			})
		})
	})
}

func TestAccountService_Update_Concurrent(t *testing.T) {
	db := setup(t)

	t.Run("Given an existing account", func(t *testing.T) {
		accEnt := repositories.AccountEntity{ID: uuid.New(), Name: "billy", Amount: 0.00}
		require.NoError(t, db.Create(&accEnt).Error)

		t.Run("And an account service", func(t *testing.T) {
			dbRepo := repositories.NewDBRepository(db)
			accService := service.NewAccountService(dbRepo)

			t.Run("When adding money multiple times at the same moment", func(t *testing.T) {
				signal := make(chan int)
				var wg sync.WaitGroup
				for i := 0; i < 10; i++ {
					wg.Add(1)
					go func(ind int) {
						defer wg.Done()
						<-signal
						_, err := accService.AddMoney(accEnt.ID, 100.00)
						require.NoError(t, err)
					}(i)
				}
				close(signal)
				wg.Wait()

				t.Run("Then result should be consistent", func(t *testing.T) {
					var currentAccEnt repositories.AccountEntity
					require.NoError(t, db.Find(&currentAccEnt, accEnt.ID).Error)
					assert.Equal(t, currentAccEnt.Amount, 1000.00)
				})
			})
		})
	})
}

func TestAccountService_Transfer(t *testing.T) {
	db := setup(t)
	t.Run("Given two existing accounts", func(t *testing.T) {
		from := repositories.AccountEntity{ID: uuid.New(), Name: "billy", Amount: 100.00}
		require.NoError(t, db.Create(&from).Error)

		to := repositories.AccountEntity{ID: uuid.New(), Name: "jhon", Amount: 100.00}
		require.NoError(t, db.Create(&to).Error)

		t.Run("And an account service", func(t *testing.T) {
			dbRepo := repositories.NewDBRepository(db)
			accService := service.NewAccountService(dbRepo)

			t.Run("When requesting to transfer money with no enough balance", func(t *testing.T) {
				err := accService.Transfer(from.ID, to.ID, 200.00)

				t.Run("Then fails", func(t *testing.T) {
					assert.Error(t, err)
				})
			})
			t.Run("When requesting to transfer money with enough balance", func(t *testing.T) {
				transferAmount := 100.00
				err := accService.Transfer(from.ID, to.ID, transferAmount)
				require.NoError(t, err)

				t.Run("Then success", func(t *testing.T) {
					var fromCurrent repositories.AccountEntity
					var toCurrent repositories.AccountEntity

					require.NoError(t, db.Find(&fromCurrent, from.ID).Error)
					require.NoError(t, db.Find(&toCurrent, to.ID).Error)

					assert.Equal(t, fromCurrent.Amount, from.Amount-transferAmount)
					assert.Equal(t, toCurrent.Amount, to.Amount+transferAmount)
				})
			})
		})
	})
}

func setup(t *testing.T) *gorm.DB {
	dsn := "test:test@tcp(localhost:3306)/bank"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&repositories.AccountEntity{}))
	return db
}
