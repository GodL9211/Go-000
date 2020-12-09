package main

import (
"context"
"fmt"
"golang.org/x/sync/errgroup"
"log"
"net/http"
"os"
"os/signal"
"syscall"
)


func Server1(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "errgroup")
})
	return Server(ctx, ":8080", mux)
}

func Server2(ctx context.Context) error {
	return Server(ctx, ":8081", http.DefaultServeMux)
}

func Server(ctx context.Context, addr string, handler http.Handler) error {
	s := http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		fmt.Println("stop")
		s.Shutdown(ctx)
	}()
	return s.ListenAndServe()
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g, _ := errgroup.WithContext(ctx)

	g.Go(func() error {
		if err := Server1(ctx); err != nil {
			cancel()
			return err
		}
		return nil
	})

	g.Go(func() error {
		if err := Server2(ctx); err != nil {
			cancel()
			return err
		}
		return nil
	})

	g.Go(func() error {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		select {
		case <-stop:
			log.Println("shutdown")
			cancel()
		case <-ctx.Done():
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}

}
