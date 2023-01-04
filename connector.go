package unixodbc

import (
	"context"
	"database/sql/driver"
)

var _ ConnectionStringFactory = (*staticString)(nil)

type staticString struct {
	string
}

func (s *staticString) ConnectionString() (string, error) {
	return s.string, nil
}

type ConnectionStringFactory interface {
	ConnectionString() (string, error)
}

var _ driver.Connector = (*connector)(nil)

type connector struct {
	initError error
	factory   ConnectionStringFactory
	tracer    Tracer
	config    *EnvConfig
}

type Option func(conn *connector) error

func Connector(options ...Option) driver.Connector {
	c := new(connector)
	for _, option := range options {
		c.initError = option(c)
		if c.initError != nil {
			return c
		}
	}
	return c
}

func WithConnectionStringFactory(factory ConnectionStringFactory) Option {
	return func(conn *connector) error {
		conn.factory = factory
		return nil
	}
}
func WithConnectionString(connStr string) Option {
	return func(conn *connector) error {
		conn.factory = &staticString{connStr}
		return nil
	}
}

func WithEnvConfig(config *EnvConfig) Option {
	return func(conn *connector) error {
		conn.config = config
		return nil
	}
}

// Connect implements driver.Connector
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {

	if c.initError != nil {
		return nil, c.initError
	}

	hndEnv, err := env(c.config)
	if err != nil {
		return nil, err
	}
	trace := tracer.Start(ctx, "Connect")
	defer trace.End()

	type result struct {
		conn *conn
		err  error
	}

	ch := make(chan *result)

	go func() {

		cn, err := hndEnv.Connection()
		if err != nil {
			ch <- &result{conn: nil, err: err}
			return
		}
		connStr, err := c.factory.ConnectionString()
		if err != nil {
			defer cn.Free()
			ch <- &result{conn: nil, err: err}
			return
		}

		if err := cn.DriverConnect(connStr); err != nil {
			defer cn.Free()
			ch <- &result{conn: nil, err: err}
			return
		}

		connection := &conn{
			psCache: psCache{
				cache: make(map[string]*cachedStatement),
			},
			connector: c,
			hnd:       cn,
			invalid:   false,
		}

		if err = ctx.Err(); err != nil {
			connection.Close()
			return
		}

		ch <- &result{conn: connection, err: nil}
	}()

	select {
	case res := <-ch:
		return res.conn, res.err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

// Driver implements driver.Connector
func (c *connector) Driver() driver.Driver {
	return driverInstance
}
