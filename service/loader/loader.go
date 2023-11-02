package loader

import "context"

type Streamer interface {
	Close()
	Next() bool
	Scan(...any) error
}

type Loader interface {
	Load(ctx context.Context) (Streamer, error)
	Listen(string) error
	Unlisten(string) error
	Receive(context.Context, ...any) error
	Close() error
}
