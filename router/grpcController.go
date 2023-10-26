package router

import (
	"fmt"
	"os"
	grpc "solution/grpc"

	"encoding/json"
	logger "solution/logger"

	"github.com/gofiber/fiber/v2"
)

// datatype for response sending to client
type Result struct {
	Alarm        bool
	Label        string
	Tagging_rate float32
}

type ResultWithSwitch struct {
	Alarm        bool
	Label        string
	Tagging_rate float32
	Switch       bool
}

// datatype for goroutine results
type ResultRoutine struct {
	Alarm        bool
	Label        string
	Tagging_rate float32
	Err          error
}

// for parsing client's request
type SignInCredentials struct {
	sound []byte
}

// datatype to round robin grpc connection
type RoundRobin struct {
	Index int
	Links []string
}

// variable to send requests as round robin
var Robin RoundRobin

// initialize round robin object
func init() {
	Robin = RoundRobin{
		Index: 0,
		Links: []string{
			"localhost:8080",
			"localhost:8080",
		},
	}
}

// round robin function.
// return url to send request, and prepare for next function call
func (Lsts *RoundRobin) next() string {
	fmt.Printf("%d\t", Lsts.Index)
	Lsts.Index += 1
	if Lsts.Index == len(Lsts.Links) {
		Lsts.Index = 0
	}
	return Lsts.Links[Lsts.Index]
}

// routing to check if ml server is connected
func PingPong(c *fiber.Ctx) error {
	logger.MyLogger.Printf("request from", c.IP())
	resultChan := make(chan bool)

	go func() {
		res := grpc.Ping(Robin.next())
		resultChan <- res
	}()
	res := <-resultChan
	if res {
		return c.SendString("Pong!")
	} else {
		return c.SendString("Unable to connect ML server")
	}
}

// routing to send sound format of bytes to ml server.
func MlServer(c *fiber.Ctx) error {
	logger.MyLogger.Printf("request from", c.IP())
	//parsing sound file from request
	body := c.Body()
	parsed := SignInCredentials{
		sound: []byte(body),
	}

	//because we use goroutine, we use channel to get result
	resultChan := make(chan ResultRoutine)

	//using goroutine, send request to ml server, and get response
	go func() {
		alarm, label, tagging_rate, err := grpc.GRPC(Robin.next(), parsed.sound)
		response := ResultRoutine{
			Alarm:        alarm,
			Label:        label,
			Tagging_rate: tagging_rate,
			Err:          err,
		}
		resultChan <- response
	}()

	//check result, and check if error exist
	res := <-resultChan
	if res.Err != nil {
		return c.SendString("GRPC error")
	}

	//if no error exist, parse the data as response type
	// response := Result{
	// 	Alarm:        res.Alarm,
	// 	Label:        res.Label,
	// 	Tagging_rate: res.Tagging_rate,
	// }
	data, err := os.ReadFile("switch.json")
	if err != nil {
		return c.SendString("Error reading the file" + err.Error())
	}

	var labelSwitches map[string]bool
	err = json.Unmarshal(data, &labelSwitches)
	if err != nil {
		return c.SendString("Error parsing JSON: " + err.Error())
	}

	switchVal, exists := labelSwitches[res.Label]
	if !exists {
		return c.SendString("Error: label not found in JSON switch file")
	}

	response := ResultWithSwitch{
		Alarm:        res.Alarm,
		Label:        res.Label,
		Tagging_rate: res.Tagging_rate,
		Switch:       switchVal,
	}

	u, err := json.Marshal(response)
	if err != nil {
		logger.MyLogger.Printf("[json] json parsing error")
		return c.SendStatus(400)
	}
	return c.SendString(string(u))
}
