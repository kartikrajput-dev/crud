package main

import "github.com/KARTIKrocks/apikit"

func main() {
	InitDB()

	app := apikit.New()

	app.POST("/users", CreateUser)
	app.GET("/users", GetUsers)
	app.GET("/users/:id", GetUser)
	app.PUT("/users/:id", UpdateUser)
	app.DELETE("/users/:id", DeleteUser)

	app.Run(":8080")
}