package socket

import (
	"context"

	socketio "github.com/googollee/go-socket.io"
)

func (r *socketIoRepository) Connects(ctx context.Context, namespace string, f func(s socketio.Conn) error) error {
	r.server.OnConnect(namespace, f)
	return nil
}

func (r *socketIoRepository) DisConnect(ctx context.Context, namespace string, f func(socketio.Conn, string)) error {
	r.server.OnDisconnect(namespace, f)
	return nil
}

func (r *socketIoRepository) Emit(s socketio.Conn, eventName, message string) {
	s.Emit(eventName, message)
}
