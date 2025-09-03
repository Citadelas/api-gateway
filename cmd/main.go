package main

import "github.com/Citadelas/api-gateway/internal/app"

//TODO:
// Better grpc error handling
// Write validate token function in sso
// Fix bug with priority not delivers

func main() {
	application, err := app.NewApp()
	if err != nil {
		panic(err)
	}
	application.Run()
}
