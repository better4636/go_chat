package hub

import (
	"github.com/better4636/go_chat/client"
)

type Message struct {
	ChatRoomId string
	Data       []byte
}

type Hub struct {
	clients   map[*client.Client]bool // 허브에서 관리되고있는 클라이언트
	OnConnect chan *client.Client     // 클라이언트 연결 시 client가 들어오는 채널
	OnClose   chan *client.Client     // 연결 종료되는 클라이언트가 들어오는 채널
	Broadcast chan Message            // hub에 존재하는 모든 클라이언트에게 메시지를 브로드 캐스트 시키기위한 채널
}

func NewHub() *Hub {
	return &Hub{
		clients:   make(map[*client.Client]bool),
		OnConnect: make(chan *client.Client),
		OnClose:   make(chan *client.Client),
		Broadcast: make(chan Message),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.OnConnect: // client가 커넥트 되었을때
			h.clients[client] = true
		case client := <-h.OnClose: // client가 연결을 종료하였을 때
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.Send) // client가 더이상 존재하지 않을 것이니 Send채널을 닫는다.
			}
		case message := <-h.Broadcast:
			for client := range h.clients {
				if client.ChatRoomId == message.ChatRoomId {
					select {
					case client.Send <- string(message.Data): // Send 채널에 데이터를 보낼 수 있으면 전송
					default: // 브로드 캐스트 대상인데 Send가 가능한 client가 아니라면 클라이언트 관리 끝냄
						delete(h.clients, client)
						close(client.Send) // client가 더이상 존재하지 않을 것이니 Send채널을 닫는다.
					}
				}
			}
		}
	}
}
