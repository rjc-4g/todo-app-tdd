package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Message string `json:"message"`
}

// アプリケーションの依存関係を保持
type Server struct {
	taskFilePath string
}

// 新しいServerインスタンスを生成
func NewServer(taskFilePath string) *Server {
	return &Server{taskFilePath: taskFilePath}
}

func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Hello World!",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func main() {
	// 本番用のファイルパスでサーバーを初期化
	server := NewServer("data/tasks.json")

	http.HandleFunc("/", server.helloHandler)

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
