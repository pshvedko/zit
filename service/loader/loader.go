package loader

import (
	"context"
	"net"
)

type Pusher interface {
	Push(int64, net.IP) error
}

type Loader interface {
	Load(context.Context, Pusher) error
	Listen(string) error
	Unlisten(string) error
	Update(context.Context, Pusher) error
	Close() error
}
