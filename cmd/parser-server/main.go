package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	parserpb "parser/gen/pb"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

// server is used to implement parserServer.
type server struct {
	parserpb.parserServer
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	host := flag.String("host", "localhost", "hostname on which the server will listen")
	port := flag.Int("port", 8888, "port on which the server will listen")
	flag.Parse()

	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", *host, *port))
	if err != nil {
		return fmt.Errorf("failed to listen %s:%d: %w", *host, *port, err)
	}

	return runServer(context.Background(), ln)
}

func runServer(ctx context.Context, ln net.Listener) error {
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	s := grpc.NewServer()
	parserpb.RegisterparserServer(s, &server{})
	log.Printf("server listening at %v", ln.Addr())
	if err := s.Serve(ln); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	eg, ctx := errgroup.WithContext(ctx)
	eg.Go(func() error {
		if err := s.Serve(ln); err != nil {
			if err != grpc.ErrServerStopped {
				return fmt.Errorf("failed to run server: %w", err)
			}
		}
		return nil

	})
	eg.Go(func() error {
		<-ctx.Done()
		_, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		s.GracefulStop()
		return nil
	})
	log.Printf("server listening on %s", ln.Addr())
	return eg.Wait()
}
