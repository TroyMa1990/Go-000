package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"golang.org/x/sync/errgroup"
)

func main() {
	var ShutDownTime int = 3
	g, ctx := errgroup.WithContext(context.Background())
	http.Handle("/", http.FileServer(http.Dir("."))
	server := &http.Server{
		Addr:    ":8080",
		Handler: http.DefaultServeMux,
	}
	c := make(chan os.Signal)
	s := make(chan error)
	
	g.Go(func() error {
			s <- server.ListenAndServe()
			fmt.Println("Http Server Throw Error",<-s)
			forceCloseSignal:="Http Server Close Signal Listen Graceful"
			c <- forceCloseSignal
		}
	})
	g.Go(func() error {
		signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		fmt.Println("Receive Signal:",<-c)
		timeoutCtx,_ := context.WithTimeout(ctx, ShutDownTime*time.Second)
		forceCloseHttpServer:="Http Server Close Http Listen Goroutine Graceful"
		s <- forceCloseHttpServer
		server.Shutdown(timeoutCtx)
		close(s)
		close(c)
		fmt.Printf("Http Shutdown Completed In %d Second",ShutDownTime)
	})
	g.Wait()
}
