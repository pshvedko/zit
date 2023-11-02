package service

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pshvedko/zit/service/loader"
)

type Service struct {
	http.Server
	sync.WaitGroup
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	one := chi.URLParam(r, "one")
	two := chi.URLParam(r, "two")

	w.Header().Set("If-Schedule-Tag-Match", fmt.Sprint(one, ",", two))
}

func (s *Service) Run(ctx context.Context, addr, port string) error {
	h := chi.NewRouter()
	h.Handle("/{one:[0-9]+}/{two:[0-9]+}", s)
	s.Handler = h
	s.Addr = net.JoinHostPort(addr, port)
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	go s.waitForContextCancel(ctx)
	return s.ListenAndServe()
}

func (s *Service) waitForContextCancel(ctx context.Context) {
	s.Add(1)
	defer s.Done()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
}

func (s *Service) Load(ctx context.Context, r loader.Loader) error {
	var id int64
	var ip net.IP
	e, err := r.Load(ctx)
	if err != nil {
		return err
	}
	defer e.Close()
	for e.Next() {
		err = e.Scan(&id, &ip)
		if err != nil {
			return err
		}
		log.Println(id, ip)
	}
	err = r.Listen("log")
	if err != nil {
		return err
	}
	go s.waitForNotification(ctx, r)
	return nil
}

func (s *Service) waitForNotification(ctx context.Context, r loader.Loader) {
	s.Add(1)
	defer s.Done()
	var id int64
	var ip net.IP
	for {
		err := r.Receive(ctx, &id, &ip)
		if err != nil {
			break
		}
		log.Println(id, ip)
	}
	_ = r.Unlisten("log")
	_ = r.Close()
}
