package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/nakatamixi/dev-span-pu-scaler/internal/scaler"
	"github.com/nakatamixi/dev-span-pu-scaler/internal/spanner"
)

type csvInstances []string

func (v *csvInstances) Set(s string) error {
	is := strings.Split(s, ",")
	if len(is) == 0 {
		return errors.New("parse instances error")
	}
	for _, s := range is {
		if s == "" {
			return errors.New("parse instances error")
		}
		*v = append(*v, s)
	}
	return nil
}

func (v *csvInstances) String() string {
	return strings.Join(*v, ",")
}

var (
	instances []string
	project   = flag.String("project", "", "gcp project")
	buffer    = flag.Int("buffer", 3, "buffer db count to scale pu")
	server    = flag.Bool("server", false, "run server(for Cloud Run)")
)

func main() {
	flag.Var((*csvInstances)(&instances), "instances", "comma separated instances")
	flag.Parse()
	if *project == "" {
		log.Fatal("need project")
	}
	exec := func() {
		for _, instance := range instances {
			err := scaleInstance(*project, instance, *buffer)
			if err != nil {
				log.Fatal(err)
			}
		}
	}

	if *server {
		ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			exec()
		})
		port := os.Getenv("PORT")
		if port == "" {
			port = "8080"
		}
		log.Printf("Listening on port %s", port)
		server := &http.Server{
			Addr:    ":" + port,
			Handler: nil,
		}
		go func() {
			<-ctx.Done()
			fmt.Println("shutting down server")
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			server.Shutdown(ctx)
		}()
		fmt.Println("start server")
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	} else {
		exec()
	}
}

func scaleInstance(projectID, instanceID string, buffer int) error {
	ctx := context.Background()
	client, err := spanner.NewClient(ctx, projectID, instanceID)
	if err != nil {
		return err
	}
	defer client.Close()

	s := scaler.NewScaler(client)
	if err := s.Execute(ctx, buffer); err != nil {
		return err
	}

	return nil
}
