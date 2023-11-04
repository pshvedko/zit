package loader

import (
	"context"
	"io"
	"net"
	"strconv"
	"strings"

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
		err = w.Push(id, ip)
		if err != nil {
			return err
		}
	}
	return r.Err()
}

func (c Conn) Update(ctx context.Context, w Pusher) error {
	n, err := c.WaitForNotification(ctx)
	if err != nil {
		return err
	}
	f := strings.Fields(n.Payload)
	if len(f) != 2 {
		return io.EOF
	}
	id, err := strconv.ParseInt(f[0], 10, 64)
	if err != nil {
		return err
	}
	ip, _, err := net.ParseCIDR(f[1])
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
	return &Conn{conn}, nil
}
