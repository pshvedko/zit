package loader

import (
	"context"
	"fmt"
	"net"

	"github.com/jackc/pgx"
)

type Conn struct {
	*pgx.Conn
}

func (c Conn) Load(ctx context.Context, w Pusher) error {
	r, err := c.QueryEx(ctx, "SELECT id, ip FROM log", nil)
	if err != nil {
		return err
	}
	defer r.Close()
	var id int64
	var ip net.IP
	for r.Next() {
		err = r.Scan(&id, &ip)
		if err != nil {
			return err
		}
		err = w.Push(id, ip.To4())
		if err != nil {
			return err
		}
	}
	return r.Err()
}

type IP struct {
	net.IP
}

func (ip *IP) Scan(r fmt.ScanState, _ rune) error {
	b, err := r.Token(true, nil)
	if err != nil {
		return err
	}
	return ip.UnmarshalText(b)
}

func (c Conn) Update(ctx context.Context, w Pusher) error {
	n, err := c.WaitForNotification(ctx)
	if err != nil {
		return err
	}
	var id int64
	var ip IP
	_, err = fmt.Sscan(n.Payload, &id, &ip)
	if err != nil {
		return err
	}
	return w.Push(id, ip.To4())
}

func New(name string) (Loader, error) {
	conf, err := pgx.ParseConnectionString(name)
	if err != nil {
		return nil, err
	}
	conn, err := pgx.Connect(conf)
	if err != nil {
		return nil, err
	}
	return Conn{conn}, nil
}
