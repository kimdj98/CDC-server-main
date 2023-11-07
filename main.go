package main

import (
	"context"
	"fmt"
	logger "solution/logger"
	router "solution/router"

	fiber "github.com/gofiber/fiber/v2"
)

var ctx = context.Background()

func main() {
	app := fiber.New(fiber.Config{AppName: "Hearsitter!"})

	app.Get("/", router.HelloWorld)   //route to check server on
	app.Get("/ping", router.PingPong) //route to check ml server on
	//app.Get("/history/latest", router.Polling)
	app.Post("/uint", router.MlServer)       //route to send byte[] to ml server
	app.Post("/switch", router.ModifySwitch) //route to change switch.json
	// app.Post("/file", router.Files)    //route to send file to ml server
	logger.Start()
	fmt.Println("Server Start!")
	app.Listen(":3000")
}
