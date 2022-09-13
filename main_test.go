package main

import (
	"bank/pkg/api/dto"
	"bank/pkg/api/repositories"
	"bank/pkg/api/service"
	"bank/pkg/app"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCreateBankAccount(t *testing.T) {
	db, err := setup()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&repositories.AccountEntity{}))

	t.Run("Given a create bank account endpoint", func(t *testing.T) {
		repo := repositories.NewDBRepository(db)
		accountService := service.NewAccountService(repo)
		router := gin.Default()
		server := app.NewServer(router, accountService)
		router.POST("/v1/account/", server.Create())

		t.Run("When request to create an account with an invalid request", func(t *testing.T) {
			request := dto.CreateAccountRequest{
				Name: "a",
			}
			jsonValue, _ := json.Marshal(request)
			req, _ := http.NewRequest("POST", "/v1/account/", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			t.Run("Then fails as bad request", func(t *testing.T) {
				assert.Equal(t, http.StatusBadRequest, w.Code)
			})
		})

		t.Run("When request to create an account with a valid request", func(t *testing.T) {
			requestOk := dto.CreateAccountRequest{
				Name:   "bob smith",
				Amount: 100.00,
			}
			jsonValueOk, _ := json.Marshal(requestOk)
			req, _ := http.NewRequest("POST", "/v1/account/", bytes.NewBuffer(jsonValueOk))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			t.Run("Then successfully created", func(t *testing.T) {
				assert.Equal(t, http.StatusCreated, w.Code)
			})
		})
	})
}

func TestUpdateBankAccount(t *testing.T) {
	db, err := setup()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&repositories.AccountEntity{}))

	t.Run("Given an existing bank account endpoint", func(t *testing.T) {
		existingAcc := repositories.AccountEntity{ID: uuid.New(), Name: "bill smith", Amount: 0.00}
		require.NoError(t, db.Create(&existingAcc).Error)

		repo := repositories.NewDBRepository(db)
		accountService := service.NewAccountService(repo)
		router := gin.Default()
		server := app.NewServer(router, accountService)
		router.PATCH("/v1/account/:accountID/money", server.AddMoney())

		t.Run("When request to add money on that account ", func(t *testing.T) {
			request := dto.UpdateAccountRequest{
				Amount: 100.00,
			}
			jsonValue, _ := json.Marshal(request)
			req, _ := http.NewRequest("PATCH", "/v1/account/"+existingAcc.ID.String()+"/money", bytes.NewBuffer(jsonValue))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			t.Run("Then update success", func(t *testing.T) {
				assert.Equal(t, http.StatusAccepted, w.Code)
			})
		})
	})
}

func TestTransfer(t *testing.T) {
	db, err := setup()
	require.NoError(t, err)
	require.NoError(t, db.AutoMigrate(&repositories.AccountEntity{}))

	t.Run("Given two existing accounts", func(t *testing.T) {
		fromAccount := repositories.AccountEntity{ID: uuid.New(), Name: "billy smith", Amount: 100.00}
		require.NoError(t, db.Create(&fromAccount).Error)

		toAccount := repositories.AccountEntity{ID: uuid.New(), Name: "jhon smith", Amount: 100.00}
		require.NoError(t, db.Create(&toAccount).Error)

		t.Run("And a transfer service api", func(t *testing.T) {
			repo := repositories.NewDBRepository(db)
			accountService := service.NewAccountService(repo)
			router := gin.Default()
			server := app.NewServer(router, accountService)
			router.POST("/v1/transfer/", server.Transfer())

			t.Run("When one account request transfer more money than current balance to the other", func(t *testing.T) {
				request := dto.TransferenceRequest{
					From:   fromAccount.ID,
					To:     toAccount.ID,
					Amount: 200.00,
				}

				jsonValue, _ := json.Marshal(request)
				req, _ := http.NewRequest("POST", "/v1/transfer/", bytes.NewBuffer(jsonValue))
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				t.Run("Then transaction fails", func(t *testing.T) {
					assert.Equal(t, http.StatusInternalServerError, w.Code)
					var fromCurrent repositories.AccountEntity
					var toCurrent repositories.AccountEntity

					require.NoError(t, db.Find(&fromCurrent, fromAccount.ID).Error)
					require.NoError(t, db.Find(&toCurrent, toAccount.ID).Error)

					assert.Equal(t, fromCurrent.Amount, fromAccount.Amount)
					assert.Equal(t, toCurrent.Amount, toAccount.Amount)
				})
			})

			t.Run("When one account request transfer money to the other", func(t *testing.T) {
				request := dto.TransferenceRequest{
					From:   fromAccount.ID,
					To:     toAccount.ID,
					Amount: 50.00,
				}

				jsonValue, _ := json.Marshal(request)
				req, _ := http.NewRequest("POST", "/v1/transfer/", bytes.NewBuffer(jsonValue))
				w := httptest.NewRecorder()
				router.ServeHTTP(w, req)

				t.Run("Then transaction fails", func(t *testing.T) {
					assert.Equal(t, http.StatusAccepted, w.Code)

					var fromCurrent repositories.AccountEntity
					var toCurrent repositories.AccountEntity

					require.NoError(t, db.Find(&fromCurrent, fromAccount.ID).Error)
					require.NoError(t, db.Find(&toCurrent, toAccount.ID).Error)

					assert.Equal(t, fromCurrent.Amount, fromAccount.Amount-50.00)
					assert.Equal(t, toCurrent.Amount, toAccount.Amount+50.00)
				})
			})
		})
	})
}

func setup() (*gorm.DB, error) {
	dsn := "test:test@tcp(localhost:3306)/bank"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}
