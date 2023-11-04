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

type Service struct {
	http.Server
	sync.WaitGroup
	sync.RWMutex
	ids map[int64]map[int32]struct{}
}

func (s *Service) Push(id int64, ip net.IP) error {
	log.Println(id, ip)
	var ipv4 int32
	err := binary.Read(bytes.NewReader(ip), binary.BigEndian, &ipv4)
	if err != nil {
		return err
	}
	s.Lock()
	defer s.Unlock()
	if s.ids == nil {
		s.ids = map[int64]map[int32]struct{}{id: {ipv4: struct{}{}}}
	} else {
		ips, ok := s.ids[id]
		if !ok {
			ips = map[int32]struct{}{}
			s.ids[id] = ips
		}
		ips[ipv4] = struct{}{}
	}
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
	err = json.NewEncoder(w).Encode(Response{Dupes: s.intersected(one, two)})
	if err != nil {
		log.Println(err)
	}
}

func (s *Service) intersected(id1, id2 int64) bool {
	if id1 == id2 {
		return true
	}
	s.RLock()
	defer s.RUnlock()
	ips1, ips2 := s.ids[id1], s.ids[id2]
	if len(ips1) > len(ips2) {
		ips1, ips2 = ips2, ips1
	}
	var n int
	for ipv4 := range ips1 {
		_, ok := ips2[ipv4]
		if ok {
			n++
			if n == 2 {
				return true
			}
		}
	}
	return false
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
	err := r.Load(ctx, s)
	if err != nil {
		return err
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
