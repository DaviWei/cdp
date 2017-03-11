// Package rpcc implements the communication layer for communicating
// with a debugging target over the Chrome Debugging Protocol.
package rpcc

import "context"

type rpcCall struct {
	Method string
	Args   interface{}
	Reply  interface{}
	Error  chan error
}

func (c *rpcCall) done(err error) {
	c.Error <- err
}

// Invoke sends an RPC request and blocks until the response is received.
// This function is called by generated code but can be used to issue
// requests manually.
func Invoke(ctx context.Context, method string, args, reply interface{}, conn *Conn) error {
	if ctx == nil {
		ctx = context.Background()
	}

	call := &rpcCall{
		Method: method,
		Args:   args,
		Reply:  reply,
		Error:  make(chan error, 1), // Do not block.
	}
	go func() {
		conn.send <- call
	}()

	select {
	case <-conn.ctx.Done():
		return ErrConnClosing
	case <-ctx.Done():
		return ctx.Err()
	case err := <-call.Error:
		return err
	}
}
