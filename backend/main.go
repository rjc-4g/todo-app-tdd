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

const dbFileName = "tasks.json" // データファイル名

// CORS対応のミドルウェア
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// preflight requestに対応
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

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

	return corsMiddleware(mux)
}

// --- ハンドラの実装 ---

func createTask(w http.ResponseWriter, r *http.Request) {
	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// 既存のタスク一覧を読み込む
	var tasks []Task
	data, err := os.ReadFile(dbFileName)
	if err == nil {
		json.Unmarshal(data, &tasks)
	} // ファイルが存在しない場合は空配列で初期化

	// 新しいタスクに ID を割り当てる（最大IDの次の値）
	maxID := 0
	for _, task := range tasks {
		if task.ID > maxID {
			maxID = task.ID
		}
	}
	newTask.ID = maxID + 1
	newTask.Deleted = false
	if newTask.Status == 0 {
		newTask.Status = 0 // ステータス 0: 未完了
	}

	// 新しいタスクを既存リストに追加
	tasks = append(tasks, newTask)

	// ファイルに保存
	fileData, err := json.Marshal(tasks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to marshal tasks"})
		return
	}

	err = os.WriteFile(dbFileName, fileData, 0644)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
		return
	}

	// 作成したタスクをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newTask)
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	// JSONファイルから全内容を読み込む
	data, err := os.ReadFile(dbFileName)
	
	// ファイルが存在しない場合は空配列を返却
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("[]"))
		return
	}
	
	// JSONファイルの全内容を返却
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}

func updateTask(w http.ResponseWriter, r *http.Request, id int) {
	// リクエストボディからsatusを取得
	var updatePayload struct {
		Status int `json:"status"`
	}
	err := json.NewDecoder(r.Body).Decode(&updatePayload)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request body"})
		return
	}

	// JSONファイルから全内容を読み込む
	data, err := os.ReadFile(dbFileName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tasks file not found"})
		return
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to unmarshal tasks"})
		return
	}

	// 指定されたIDのタスクを探してステータスを更新
	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Status = updatePayload.Status
			found = true
			break
		}
	}

	if !found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	// ファイルに保存
	fileData, err := json.Marshal(tasks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to marshal tasks"})
		return
	}

	err = os.WriteFile(dbFileName, fileData, 0644)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
		return
	}

	// 更新されたタスクをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	for _, task := range tasks {
		if task.ID == id {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
}

func deleteTask(w http.ResponseWriter, r *http.Request, id int) {
	// JSONファイルから全内容を読み込む
	data, err := os.ReadFile(dbFileName)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Tasks file not found"})
		return
	}

	var tasks []Task
	err = json.Unmarshal(data, &tasks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to unmarshal tasks"})
		return
	}

	// 指定されたIDのタスクを探して削除フラグを更新
	found := false
	for i, task := range tasks {
		if task.ID == id {
			tasks[i].Deleted = true
			found = true
			break
		}
	}

	if !found {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"error": "Task not found"})
		return
	}

	// ファイルに保存
	fileData, err := json.Marshal(tasks)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to marshal tasks"})
		return
	}

	err = os.WriteFile(dbFileName, fileData, 0644)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": "Failed to write file"})
		return
	}

	// 更新されたタスクをレスポンスとして返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	for _, task := range tasks {
		if task.ID == id {
			json.NewEncoder(w).Encode(task)
			return
		}
	}
}

func main() {
	router := setupRouter()
	http.ListenAndServe(":8080", router)
}