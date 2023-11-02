package postgres

import (
	"context"
	"encoding"
	"fmt"
	"strings"

	"github.com/jackc/pgx"

	"github.com/pshvedko/zit/service/loader"
)

type Conn struct {
	*pgx.Conn
}

type Rows struct {
	*pgx.Rows
}

func (c Conn) Load(ctx context.Context) (loader.Streamer, error) {
	rows, err := c.QueryEx(ctx, "SELECT id, ip FROM log", nil)
	if err != nil {
		return nil, err
	}
	return &Rows{rows}, nil
}

func (c Conn) Receive(ctx context.Context, args ...any) error {
	notification, err := c.WaitForNotification(ctx)
	if err != nil {
		return err
	}
	for i, f := range strings.Fields(notification.Payload) {
		switch v := args[i].(type) {
		case encoding.TextUnmarshaler:
			err = v.UnmarshalText([]byte(f))
		default:
			_, err = fmt.Sscan(f, v)
		}
		if err != nil {
			break
		}
	}
	return err
}

func New(name string) (loader.Loader, error) {
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
