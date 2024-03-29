package main

//go:generate go run github.com/swaggo/swag/cmd/swag init -o ./swagger
//go:generate go run github.com/swaggo/swag/cmd/swag fmt

import "github.com/oechsler-it/identity/app"

//	@title			identity
//	@version		0.0.1
//	@description	A minimal identity provider

//	@basePath	/

// @securityDefinitions.apikey	TokenAuth
// @in							header
// @name						Authorization
// @description				Bearer token authentication
func main() {
	app.New().Run()
}
