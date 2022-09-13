package main

import (
	"bank/pkg/api/repositories"
	"bank/pkg/api/service"
	"bank/pkg/app"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\\n", err)
		os.Exit(1)
	}
}

func run() error {
	println("loading accountService...")

	dsn := "test:test@tcp(db:3306)/bank"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.AutoMigrate(&repositories.AccountEntity{})
	if err != nil {
		log.Fatal("failed to load table")
	}

	dbRepo := repositories.NewDBRepository(db)
	accountService := service.NewAccountService(dbRepo)

	router := gin.Default()

	server := app.NewServer(router, accountService)
	sErr := server.Run()
	if err != nil {
		return err
	}

	return sErr
}
