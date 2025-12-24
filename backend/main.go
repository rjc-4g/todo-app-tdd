package main

import (
	"encoding/json" // 追加
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	filePath := os.Getenv("DATA_FILE_PATH")
	if filePath == "" {
		filePath = "data/todos.json"
	}

	repo := &JSONRepository{FilePath: filePath}

	http.HandleFunc("/todos", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "application/json")

		if r.Method == http.MethodGet {
			todos, err := repo.Load()
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			
			// 修正箇所: todos を使ってレスポンスを返す
			if err := json.NewEncoder(w).Encode(todos); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
	})

	fmt.Println("Server started at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}