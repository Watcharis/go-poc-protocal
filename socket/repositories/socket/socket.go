package socket

import (
	"context"
	"fmt"

	socketio "github.com/googollee/go-socket.io"
)

type OnConnectFunc func(s socketio.Conn) error
type OnDisconnectFunc func(socketio.Conn, string)
type OnEvent interface{}

type SocketIoRepository interface {
	// ----- customs generic type version ------
	Connect(ctx context.Context, namespace string, f OnConnectFunc) error
	OnEvent(ctx context.Context, namespace string, event string, f OnEvent) error
	Behavier(s socketio.Conn, f any, data ...interface{}) error

	// --------- general version ----------
	Connects(ctx context.Context, namespace string, f func(s socketio.Conn) error) error
	DisConnect(ctx context.Context, namespace string, f func(socketio.Conn, string)) error
	Emit(s socketio.Conn, eventName, message string)
}

type socketIoRepository struct {
	server *socketio.Server
}

func NewSocketIoRepository(server *socketio.Server) SocketIoRepository {
	return &socketIoRepository{
		server: server,
	}
}

func (r *socketIoRepository) Behavier(s socketio.Conn, f any, data ...interface{}) error {
	switch fn := any(f).(type) {
	case OnConnectFunc:
		if err := fn(s); err != nil {
			return err
		}
	case OnDisconnectFunc:
		reason, ok := data[0].(string)
		if !ok {
			return fmt.Errorf("OnDisconnectFunc - invalid data type is string")
		}
		fn(s, reason)
		return nil
	}
	return nil
}

func (r *socketIoRepository) Connect(ctx context.Context, namespace string, f OnConnectFunc) error {
	r.server.OnConnect(namespace, func(s socketio.Conn) error {
		r.Behavier(s, f)
		return nil
	})
	return nil
}

func (r *socketIoRepository) DisConnects(ctx context.Context, namespace string, f OnConnectFunc) error {
	r.server.OnDisconnect(namespace, func(s socketio.Conn, reason string) {
		r.Behavier(s, f, reason)
	})
	return nil
}

func (r *socketIoRepository) OnEvent(ctx context.Context, namespace string, event string, f OnEvent) error {
	r.server.OnEvent(namespace, event, f)
	return nil
}
