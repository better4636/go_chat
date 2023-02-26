package client

type Client struct {
	ChatRoomId string      // 클라이언트가 연결된 채팅방
	Send       chan string // 메시지를 전송할 채널
}
