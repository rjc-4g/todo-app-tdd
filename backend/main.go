package main

import (
	"fmt"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	// レスポンスとして文字列を書き込む
	fmt.Fprint(w, "Hello, World!")
}

func main() {
	// ルートパス "/" にアクセスしたときに helloHandler を実行する設定
	http.HandleFunc("/", helloHandler)

	fmt.Println("Server is running on http://localhost:8080")
	
	// 8080ポートでサーバーを起動
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}