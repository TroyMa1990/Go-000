package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	pb "app/api/user/v1"
	di "app/internal/di"
	"app/internal/pkg/grpc"
	"app/internal/service"

	"golang.org/x/sync/errgroup"
)

const (
	address = ":18031"
)

func main() {
	usr := di.InitFindUser()
	service := service.NewUserService(usr)

	s := grpc.NewServer(address)
	pb.RegisterUserServer(s.Server, service)
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() error {
		return s.Start(ctx)
	})

	// signal
	g.Go(func() error {
		exitSignals := []os.Signal{os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT} // SIGTERM is POSIX specific
		sig := make(chan os.Signal, len(exitSignals))
		signal.Notify(sig, exitSignals...)
		for {
			fmt.Println("signal")
			select {
			case <-ctx.Done():
				fmt.Println("signal ctx done")
				return ctx.Err()
			case <-sig:
				return nil
			}
		}
	})

	err := g.Wait() // first error return
	fmt.Println(err)
}
