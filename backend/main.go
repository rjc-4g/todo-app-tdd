package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

// TODOタスクを表す構造体。
type Task struct {
	ID      int       `json:"id"`      // 自動採番される一意なID
	Name    string    `json:"name"`    // タスク名（最大40文字）
	Status  int       `json:"status"`  // 状態（0: 未完了 / 1: 完了）
	Created time.Time `json:"created"` // 作成日時
	Updated time.Time `json:"updated"` // 最終更新日時
	Deleted bool      `json:"deleted"` // 論理削除フラグ
}

var (
	mu       sync.RWMutex   // 並行アクセスからstoreを保護するロック
	store    []Task         // タスクのインメモリストア（JSONファイルと同期）
	nextID   = 1            // 次に発行するタスクID
	dataFile = "tasks.json" // 永続化先のJSONファイルパス（テストで差し替え可能）
)

// JSONファイルからタスクをメモリに読み込む。
// ファイルが存在しない場合は空のストアで初期化する。
func loadStore() error {
	mu.Lock()
	defer mu.Unlock()

	data, err := os.ReadFile(dataFile)
	if err != nil {
		if os.IsNotExist(err) {
			store = make([]Task, 0)
			nextID = 1
			return nil
		}
		return err
	}

	if len(data) == 0 {
		store = make([]Task, 0)
		nextID = 1
		return nil
	}

	if err := json.Unmarshal(data, &store); err != nil {
		return err
	}

	// 既存データの最大IDを元にnextIDを設定（ID重複を防ぐ）
	nextID = 1
	for _, t := range store {
		if t.ID >= nextID {
			nextID = t.ID + 1
		}
	}
	return nil
}

// メモリ上のストアをJSONファイルに書き出す。
// 呼び出し元がすでに mu.Lock を保持している前提で動作する（二重ロック防止）。
func saveStore() {
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		log.Printf("saveStore: marshal error: %v", err)
		return
	}
	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		log.Printf("saveStore: write error: %v", err)
	}
}

// ストアをクリアしIDカウンターをリセットする。
// テストで特定のケースのクリーンな状態を確立するために使用する。
func resetStore() {
	mu.Lock()
	defer mu.Unlock()
	store = make([]Task, 0)
	nextID = 1
}

// レスポンスにステータスコードとJSONボディを書き出す。
func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v) //nolint:errcheck
}

// エラーメッセージをJSON形式でレスポンスに書き出す。
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

// URLパスの末尾からタスクIDを取り出す。
// 例: /api/v1/tasks/123 → 123
func extractIDFromPath(r *http.Request) (int, error) {
	path := strings.TrimRight(r.URL.Path, "/")
	parts := strings.Split(path, "/")
	return strconv.Atoi(parts[len(parts)-1])
}

// リクエストボディのname・statusを検証してタスクを作成する。
func createTaskHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name   string `json:"name"`
		Status *int   `json:"status"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name is required")
		return
	}

	// 文字数チェック
	if utf8.RuneCountInString(req.Name) > 40 {
		writeError(w, http.StatusBadRequest, "name must be 40 characters or less")
		return
	}

	// statusが未指定の場合はデフォルト値0（未完了）を使用する
	status := 0
	if req.Status != nil {
		if *req.Status < 0 || *req.Status > 1 {
			writeError(w, http.StatusBadRequest, "status must be 0 or 1")
			return
		}
		status = *req.Status
	}

	now := time.Now()
	mu.Lock()
	task := Task{
		ID:      nextID,
		Name:    req.Name,
		Status:  status,
		Created: now,
		Updated: now,
		Deleted: false,
	}
	nextID++
	store = append(store, task)
	saveStore()
	mu.Unlock()

	writeJSON(w, http.StatusCreated, task)
}

// 論理削除されていないタスクの一覧を返す。
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	mu.RLock() // 読み取り専用ロック（複数の読み取りを並行可能）
	result := make([]Task, 0)
	for _, t := range store {
		if !t.Deleted {
			result = append(result, t)
		}
	}
	mu.RUnlock()

	writeJSON(w, http.StatusOK, result)
}

// statusまたはnameを更新し、更新後のタスクを返す。
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	var req struct {
		Status *int    `json:"status"`
		Name   *string `json:"name"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON")
		return
	}

	// statusとnameのどちらも指定されていない場合はエラー
	if req.Status == nil && req.Name == nil {
		writeError(w, http.StatusBadRequest, "status is required")
		return
	}

	if req.Status != nil && (*req.Status < 0 || *req.Status > 1) {
		writeError(w, http.StatusBadRequest, "status must be 0 or 1")
		return
	}

	if req.Name != nil {
		if *req.Name == "" {
			writeError(w, http.StatusBadRequest, "name is required")
			return
		}
		if utf8.RuneCountInString(*req.Name) > 40 {
			writeError(w, http.StatusBadRequest, "name must be 40 characters or less")
			return
		}
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range store {
		if store[i].ID == id {
			if req.Status != nil {
				store[i].Status = *req.Status
			}
			if req.Name != nil {
				store[i].Name = *req.Name
			}
			store[i].Updated = time.Now()
			saveStore()
			writeJSON(w, http.StatusOK, store[i])
			return
		}
	}

	writeError(w, http.StatusNotFound, "task not found")
}

// deleted フラグを true にする論理削除を行う。
func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := extractIDFromPath(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task ID")
		return
	}

	mu.Lock()
	defer mu.Unlock()

	for i := range store {
		if store[i].ID == id {
			store[i].Deleted = true
			store[i].Updated = time.Now()
			saveStore()
			writeJSON(w, http.StatusOK, store[i])
			return
		}
	}

	writeError(w, http.StatusNotFound, "task not found")
}

func main() {
	// 起動時にJSONファイルからタスクをメモリに読み込む
	if err := loadStore(); err != nil {
		log.Printf("Warning: could not load data file (%s): %v", dataFile, err)
	}

	mux := http.NewServeMux()

	// タスク一覧取得・新規作成
	mux.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getTasksHandler(w, r)
		case http.MethodPost:
			createTaskHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// タスクの更新・削除
	mux.HandleFunc("/api/v1/tasks/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPatch:
			updateTaskHandler(w, r)
		case http.MethodDelete:
			deleteTaskHandler(w, r)
		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Printf("Server starting on :8080 (data file: %s)", dataFile)
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}
}
