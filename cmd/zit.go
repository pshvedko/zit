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
)

func main() {
	var s service.Service
	var addrFlag string
	var portFlag string
	var baseFlag string

	defer s.Wait()
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	c := &cobra.Command{
		Use:  "zit",
		Long: "Duplicate detection microservice",
		PreRunE: func(cmd *cobra.Command, args []string) error {
			r, err := loader.New(baseFlag)
			if err != nil {
				return err
			}
			return s.Load(cmd.Context(), r)
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := s.Run(cmd.Context(), addrFlag, portFlag)
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

	err := c.ExecuteContext(ctx)
	if err != nil {
		os.Exit(1)
	}
}
