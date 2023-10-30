package main

import (
	"errors"
	"log"
	"net"
	"os"
	"os/signal"
	"path/filepath"
	"strconv"
	"syscall"

	"google.golang.org/grpc"

	hc "test/grpc-demo/app"

	pb "google.golang.org/grpc/health/grpc_health_v1"
	cli "gopkg.in/urfave/cli.v1"
)

const (
	VERSION = "1.0.1"
	USAGE   = "grpc health check server"
)

var app *cli.App

func init() {
	app = cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Version = VERSION
	app.Usage = USAGE
	app.Flags = []cli.Flag{
		cli.UintFlag{Name: "port, p", Usage: "端口"},
	}
	app.Action = func(ctx *cli.Context) error {
		p := ctx.GlobalUint("port")
		if p == 0 {
			log.Fatalf("Missing port!")
			return errors.New("Missing port!")
		}
		grpcServer := grpc.NewServer()
		lis, err := net.Listen("tcp", ":"+strconv.Itoa(int(p)))
		if err != nil {
			log.Fatalf("Failed to listen:%+v", err)
			return err
		}
		pb.RegisterHealthServer(grpcServer, hc.New())
		go func() {
			sigs := make(chan os.Signal, 1)
			signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
			_ = <-sigs
			grpcServer.GracefulStop()
		}()
		log.Printf("service started")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %+v", err)
			return err
		}
		return nil
	}
}
func main() {
	if err := app.Run(os.Args); err != nil {
		os.Exit(1)
	}
}
