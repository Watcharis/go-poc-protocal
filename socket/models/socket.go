package models

const (
	NAMESPACE_CHAT_APP    string = "/chat-app"
	NAMESPACE_CHAT_SERVER string = "/chat-server"
)

const (
	EVENT_NAME_CHAT         string = "chat"
	EVENT_NAME_REPLY        string = "reply"
	EVENT_NAME_WELCOME      string = "welcome"
	EVENT_NAME_JOIN_ROOM    string = "join_room"
	EVENT_NAME_ROOM_JOINED  string = "room_joined"
	EVENT_NAME_ROOM_MESSAGE string = "room_message"
	EVENT_NAME_ROOM_CYCLE   string = "room_cycle"
	EVENT_NAME_LEAVE_ROOM   string = "leave_room"
)

type Message struct {
	Text        string `json:"text"`
	RoomName    string `json:"room_name"`
	AccessToken string `json:"access_token"`
	OwnerChatID string `json:"owner_chat_id"`
}
