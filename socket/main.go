package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"watcharis/go-poc-protocal/pkg"
	"watcharis/go-poc-protocal/socket/models"
	"watcharis/go-poc-protocal/socket/repositories/socket"

	socketio "github.com/googollee/go-socket.io"
	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/session"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"
)

const (
	SEVER_PORT = ":8999"
)

func main() {
	// สร้าง Socket.IO server
	ctx := context.Background()

	server := InitSocketIo()

	socketRepository := socket.NewSocketIoRepository(server)

	// กำหนด event เมื่อ client เชื่อมต่อ
	// server.OnConnect("/chat-app", func(s socketio.Conn) error {
	// 	fmt.Println("Client connected:", s.ID())
	// 	s.Emit("welcome", "Welcome to the Socket.IO server")
	// 	return nil
	// })

	var connectBehavier socket.OnConnectFunc = func(s socketio.Conn) error {
		fmt.Println("Client connected:", s.ID())
		socketRepository.Emit(s, models.EVENT_NAME_WELCOME, "Welcome to the Socket.IO server")
		return nil
	}

	// connect to namespace "/chat-app"
	if err := socketRepository.Connect(ctx, models.NAMESPACE_CHAT_APP, connectBehavier); err != nil {
		fmt.Println("[Error] socket connect err :", err)
		log.Panic(err)
	}

	// connect to namespace "/chat-server"
	if err := socketRepository.Connects(ctx, models.NAMESPACE_CHAT_SERVER, func(s socketio.Conn) error {
		fmt.Println("Client connected:", s.ID())
		message := fmt.Sprintf("server is connect namespace: %s success.", models.NAMESPACE_CHAT_SERVER)
		socketRepository.Emit(s, models.EVENT_NAME_WELCOME, message)
		return nil
	}); err != nil {
		log.Panic(err)
	}

	// กำหนด event สำหรับรับข้อความจาก client
	// server.OnEvent("/chat-app", "chat", func(s socketio.Conn, msg string) {
	// 	fmt.Printf("Message from client (ch.reply): %s\n", msg)
	// 	s.Emit("reply", "Server received : "+msg)
	// })

	var onEvent socket.OnEvent = func(s socketio.Conn, message string) {
		fmt.Printf("Message from client (ch.reply): %s\n", message)
		socketRepository.Emit(s, models.EVENT_NAME_REPLY, message)
	}

	if err := socketRepository.OnEvent(ctx, models.NAMESPACE_CHAT_APP, models.EVENT_NAME_CHAT, onEvent); err != nil {
		log.Panicf("[ERROR] OnEvent failed : %+v", err)
	}

	server.OnEvent(models.NAMESPACE_CHAT_SERVER, models.EVENT_NAME_CHAT, func(s socketio.Conn, msg string) {
		fmt.Printf("Message from client (ch.reply): %s\n", msg)
		s.Emit(models.EVENT_NAME_REPLY, "Server received : "+msg)
	})

	//  ------------------------------- Room chat -------------------------------------
	server.OnEvent(models.NAMESPACE_CHAT_APP, models.EVENT_NAME_JOIN_ROOM, func(s socketio.Conn, message string) {
		fmt.Println("Channel.[ join_room ] message :", message)
		var data models.Message
		if err := json.Unmarshal([]byte(message), &data); err != nil {
			log.Printf("[ERROR] json.Unmarshal failed err : %+v\n", err)
		}

		room := data.RoomName
		fmt.Println("event join_room | room_name:", room)
		// ให้ client เข้าร่วม Room
		s.Join(room)

		// แจ้งให้ทราบว่า Room ถูกสร้าง (หรือ client เข้าร่วมสำเร็จ)
		fmt.Printf("Client %s joined room: %s\n", s.ID(), room)
		s.Emit(models.EVENT_NAME_ROOM_JOINED, message)
	})

	server.OnEvent(models.NAMESPACE_CHAT_APP, models.EVENT_NAME_ROOM_MESSAGE, func(s socketio.Conn, room, message string) {
		// ส่งข้อความไปยัง Room
		fmt.Printf("(ch.room_message) recieve event -> room: %s, message: %s\n", room, message)
		fmt.Printf("ch.room_message socker-ID : %s\n", s.ID())

		var data models.Message
		if err := json.Unmarshal([]byte(message), &data); err != nil {
			log.Printf("[ERROR] json.Unmarshal failed err : %+v\n", err)
		}

		fmt.Println("rooms :", s.Rooms())

		send_message := fmt.Sprintf(`[Room_%s] %s : %s`, data.RoomName, data.OwnerChatID, data.Text)

		server.BroadcastToRoom(models.NAMESPACE_CHAT_APP, room, models.EVENT_NAME_ROOM_MESSAGE, send_message)
	})

	// Event สำหรับออกจาก Room
	server.OnEvent(models.NAMESPACE_CHAT_APP, models.EVENT_NAME_LEAVE_ROOM, func(s socketio.Conn, room string) {
		// ออกจาก room
		s.Leave(room)
		fmt.Printf("Client %s left room: %s\n", s.ID(), room)
		s.Emit(models.EVENT_NAME_ROOM_CYCLE, "You left room: "+room)
	})

	// กำหนด event เมื่อ client ตัดการเชื่อมต่อ
	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("Client disconnected:", s.ID(), "Reason:", reason)
	})

	// เริ่ม server HTTP
	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("Socket.IO server error: %v", err)
		}
	}()

	// close socket.io
	defer server.Close()

	handler := InitRouter(server)

	// ports := []string{":8998", ":8999"}
	// for _, port := range ports {

	httpServer := http.Server{
		Addr:    SEVER_PORT,
		Handler: handler,
	}

	go func(port string) {
		defer httpServer.Close()

		fmt.Printf("Socket.IO server running on http://localhost%s\n", port)
		if err := httpServer.ListenAndServe(); err != nil {
			log.Println("[error] cannot start server :", err)
			log.Panic(err)
		}
	}(httpServer.Addr)
	// }

	wg := new(sync.WaitGroup)
	signal := make(chan os.Signal, 1)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		s := <-signal
		fmt.Println("signal :", s)
	}()
	wg.Wait()
}

var allowOriginFunc = func(r *http.Request) bool {
	fmt.Printf("allow-origin-func reuqest : %+v\n", r)
	return true
}

func InitSocketIo() *socketio.Server {
	opts := &engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
		SessionIDGenerator: &session.DefaultIDGenerator{},
		ConnInitor: engineio.ConnInitorFunc(func(r *http.Request, c engineio.Conn) {
			fmt.Printf("request after create connection : %+v\n", r)
			url := c.URL()
			id := c.ID()
			fmt.Printf("url : %+v, engineio_id : %+v\n", url, id)
		}),
	}
	server := socketio.NewServer(opts)
	return server
}

func InitRouter(server *socketio.Server) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", pkg.HealthCheck)

	// เชื่อม Socket.IO เข้ากับ HTTP server
	// mux.HandleFunc("/socket.io/", server.ServeHTTP)
	mux.Handle("/socket.io/", server)
	mux.Handle("/", http.FileServer(http.Dir("./public")))

	return mux
}
