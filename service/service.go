package service

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/pshvedko/zit/service/loader"
)

type Storage interface {
	Insert(id int64, ip int32) bool
	Intersected(id1, id2 int64) bool
}

type Service struct {
	http.Server
	sync.WaitGroup
	Storage
}

func (s *Service) Push(id int64, ip net.IP) error {
	log.Println(id, ip)
	var ipv4 int32
	err := binary.Read(bytes.NewReader(ip), binary.BigEndian, &ipv4)
	if err != nil {
		return err
	}
	s.Insert(id, ipv4)
	return nil
}

type Response struct {
	Dupes bool `json:"dupes"`
}

func (s *Service) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	one, err := strconv.ParseInt(chi.URLParam(r, "one"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	two, err := strconv.ParseInt(chi.URLParam(r, "two"), 10, 64)
	if err != nil {
		http.NotFound(w, r)
		log.Println(err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(Response{
		Dupes: s.Intersected(one, two),
	})
	if err != nil {
		log.Println(err)
	}
}

func (s *Service) Run(ctx context.Context, addr, port string) error {
	h := chi.NewRouter()
	h.Handle("/{one:[0-9]+}/{two:[0-9]+}", s)
	s.Handler = h
	s.Addr = net.JoinHostPort(addr, port)
	s.BaseContext = func(net.Listener) context.Context {
		return ctx
	}
	go s.WaitForContextCancel(ctx)
	return s.ListenAndServe()
}

func (s *Service) WaitForContextCancel(ctx context.Context) {
	s.Add(1)
	defer s.Done()
	<-ctx.Done()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = s.Shutdown(ctx)
}

func (s *Service) Load(ctx context.Context, r loader.Loader) error {
	err := r.Load(ctx, s)
	if err != nil {
		return err
	}
	err = r.Listen("log")
	if err != nil {
		return err
	}
	go s.WaitForNotification(ctx, r)
	return nil
}

func (s *Service) WaitForNotification(ctx context.Context, r loader.Loader) {
	s.Add(1)
	defer s.Done()
	for {
		err := r.Update(ctx, s)
		if err != nil {
			log.Println(err)
			break
		}
	}
	_ = r.Unlisten("log")
	_ = r.Close()
}
