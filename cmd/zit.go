package main

import (
	"context"
	"errors"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/pshvedko/zit/service"
	"github.com/pshvedko/zit/service/loader"
	"github.com/pshvedko/zit/storage"
)

func main() {
	var addrFlag string
	var portFlag string
	var baseFlag string

	s := service.Service{
		Storage: &storage.ArrayIntersection{},
	}
	defer s.Wait()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := &cobra.Command{
		Use:  "zit",
		Long: "Duplicate detection microservice",
		PreRunE: func(*cobra.Command, []string) error {
			r, err := loader.New(baseFlag)
			if err != nil {
				return err
			}
			return s.Load(ctx, r)
		},
		RunE: func(*cobra.Command, []string) error {
			err := s.Run(ctx, addrFlag, portFlag)
			switch {
			case errors.Is(err, http.ErrServerClosed):
				return nil
			default:
				return err
			}
		},
	}

	c.Flags().StringVar(&addrFlag, "addr", "", "bind address")
	c.Flags().StringVar(&portFlag, "port", "8080", "bind port")
	c.Flags().StringVar(&baseFlag, "db", "postgres://postgres:postgres@postgres:5432/zit", "data base")

	err := c.Execute()
	if err != nil {
		os.Exit(1)
	}
}
