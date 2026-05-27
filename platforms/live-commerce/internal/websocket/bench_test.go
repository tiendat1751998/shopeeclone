package websocket

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"
)

func BenchmarkJSONMarshalStdlib(b *testing.B) {
	msg := map[string]interface{}{
		"type":      "chat",
		"payload":   map[string]string{"user_id": "u1", "content": "hello world this is a test message"},
		"timestamp": time.Now().UnixMilli(),
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = json.Marshal(msg)
	}
}

func BenchmarkJSONUnmarshalStdlib(b *testing.B) {
	data := []byte(`{"type":"chat","payload":{"user_id":"u1","content":"hello world"},"timestamp":1234567890}`)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var msg map[string]interface{}
		_ = json.Unmarshal(data, &msg)
	}
}

func BenchmarkSprintfKey(b *testing.B) {
	roomID := "room-12345"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("room:%s:viewers", roomID)
	}
}

func BenchmarkStringConcatKey(b *testing.B) {
	roomID := "room-12345"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = "room:" + roomID + ":viewers"
	}
}

func BenchmarkBroadcastToRoom(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	numClients := 50
	for i := 0; i < numClients; i++ {
		client := NewClient(fmt.Sprintf("c%d", i), "user1", "user1", nil, hub)
		client.RoomID = "room-bench"
		hub.Register <- client
	}

	data := []byte(`{"type":"chat","payload":{"content":"hello"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		hub.BroadcastToRoomBytes("room-bench", data, "")
	}
	b.StopTimer()
}

func BenchmarkBroadcastToRoomParallel(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	numClients := 50
	for i := 0; i < numClients; i++ {
		client := NewClient(fmt.Sprintf("c%d", i), "user1", "user1", nil, hub)
		client.RoomID = "room-bench"
		hub.Register <- client
	}

	data := []byte(`{"type":"chat","payload":{"content":"hello"}}`)

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			hub.BroadcastToRoomBytes("room-bench", data, "")
		}
	})
	b.StopTimer()
}

func BenchmarkHubRegister(b *testing.B) {
	hub := NewHub()
	go hub.Run()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client := NewClient(fmt.Sprintf("client-%d", i), "user1", "user1", nil, hub)
		client.RoomID = "bench-room"
		hub.Register <- client
	}
	b.StopTimer()
}

func BenchmarkRoomBroadcastBytes(b *testing.B) {
	room := NewRoom("test-room")
	for i := 0; i < 100; i++ {
		c := NewClient(fmt.Sprintf("c%d", i), "user1", "user1", nil, nil)
		room.Add(c)
	}
	data := []byte(`{"type":"chat","payload":{"content":"hello"}}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		room.BroadcastBytes(data, "")
	}
}

func BenchmarkRoomBroadcastBytesParallel(b *testing.B) {
	room := NewRoom("test-room")
	for i := 0; i < 100; i++ {
		c := NewClient(fmt.Sprintf("c%d", i), "user1", "user1", nil, nil)
		room.Add(c)
	}
	data := []byte(`{"type":"chat","payload":{"content":"hello"}}`)

	b.ResetTimer()
	var mu sync.Mutex
	_ = mu
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			room.BroadcastBytes(data, "")
		}
	})
}
