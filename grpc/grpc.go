package grpc

import (
	"context"
	"fmt"
	pb "solution/grpc/sound/sound"
	"time"

	logger "solution/logger"

	"google.golang.org/grpc"
)

/*
using defined grpc, check connection with ml server.
*/
func Ping(path string) bool {

	//open grpc connection
	conn, err := grpc.Dial(path, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.MyLogger.Printf("did not connect: %v", err)
		return false
	}
	defer conn.Close()
	c := pb.NewFileClient(conn)

	//checking timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	//send ping pong, check if error exist
	r, err := c.Connect(ctx, &pb.Ping{Ping: "Ping!"})
	if err != nil {
		logger.MyLogger.Printf("could not request: %v", err)
		return false
	}

	//using logger, if error not defined, return true for connection success.
	logger.MyLogger.Printf("Config: %v", r)
	return true
}

/*
using defined grpc, send sound file and return results.
*/
func GRPC(path string, file []byte) (bool, string, float32, error) {

	//open grpc connection
	conn, err := grpc.Dial(path, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		logger.MyLogger.Printf("[GRPC error] did not connect: %v", err)
		return false, "", 0.0, err
	}
	defer conn.Close()
	c := pb.NewFileClient(conn)

	//checking timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	//send sound file, and check if error exist
	r, err := c.Define(ctx, &pb.SoundRequest{Sound: file})
	if err != nil {
		logger.MyLogger.Printf("[GRPC error] could not request: %v", err)
		return false, "", 0.0, err
	}

	//after some logging, return results
	s := fmt.Sprintf("%t", r.GetAlarm())
	logger.MyLogger.Printf("Config : { Alarm : "+s+", Label : "+r.GetRes()+", Tagging_rate : %f }", r.GetTaggingRate())
	return r.GetAlarm(), r.GetRes(), r.GetTaggingRate(), nil
}
