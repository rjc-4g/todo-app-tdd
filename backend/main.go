package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// Task 構造体の定義
type Task struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Status  int    `json:"status"`
	Deleted bool   `json:"deleted"`
}

const dbFileName = "tasks_test.json" // テストと共通のファイル名

// 1. ルーターの設定（テストからも呼び出せるように外出しする）
func setupRouter() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			createTask(w, r)
		case http.MethodGet:
			getTasks(w, r)
		}
	})

	// ID付きのパス（簡易的な実装）
	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/api/v1/tasks/")
		id, _ := strconv.Atoi(idStr)

		switch r.Method {
		case http.MethodPatch:
			updateTask(w, r, id)
		case http.MethodDelete:
			deleteTask(w, r, id)
		}
	})

	return mux
}

// --- ハンドラの実装 ---

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	json.NewDecoder(r.Body).Decode(&task)
	task.ID = 1 // 簡易的にID固定

	// ファイル保存（本来は既存データを読み込んで追加するが、ここでは簡易化）
	data, _ := json.Marshal([]Task{task})
	os.WriteFile(dbFileName, data, 0644)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(task)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	data, _ := os.ReadFile(dbFileName)
	w.Header().Set("Content-Type", "application/json")
	w.Write(data)
}

func updateTask(w http.ResponseWriter, r *http.Request, id int) {
	w.WriteHeader(http.StatusOK) // 成功のみ返す簡易実装
}

func deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	w.WriteHeader(http.StatusOK) // 成功のみ返す簡易実装
}

func main() {
	router := setupRouter()
	http.ListenAndServe(":8080", router)
}