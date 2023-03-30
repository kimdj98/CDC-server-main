package logger

import (
	"io"
	"log"
	"os"
)

var MyLogger *log.Logger

func Start() {
	//using logger to save log
	fpLog, err := os.OpenFile("./logs/logfile.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	MyLogger = log.New(fpLog, "INFO: ", log.Ldate|log.Ltime)

	multiWriter := io.MultiWriter(fpLog, os.Stdout)
	MyLogger.SetOutput(multiWriter)
	MyLogger.Print("")
	MyLogger.Print("Program start")
}
