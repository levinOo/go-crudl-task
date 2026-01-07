package main

//	@title			Subscription CRUD API
//	@version		1.0
//	@description	API для управления подписками пользователей

//	@host		localhost:8080
//	@BasePath	/api/v1

import (
	"os"

	"github.com/levinOo/go-crudl-task/internal/app"
)

func main() {
	// Запускаем приложение
	if err := app.Run(); err != nil {
		os.Exit(1)
	}
}
