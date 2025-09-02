package main

import "github.com/Citadelas/api-gateway/internal/app"

//TODO:
// Implement logger
// Better grpc error handling
// Write validate token function in sso
// Refactor server initialization
// Fix bug with priority not delivers

func main() {
	application, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	application.Run()
}
