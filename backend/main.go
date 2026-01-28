package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"
)

// Response はAPIレスポンスの汎用的なメッセージ構造体
type Response struct {
	Message string `json:"message"`
}

// Task は1つのタスクを表すデータモデル
type Task struct {
	Id      int       `json:"id"`      // タスクID (一意)
	Name    string    `json:"name"`    // タスク名
	Status  int       `json:"status"`  // タスクのステータス (0: 未, 1: 済)
	Created time.Time `json:"created"` // 作成日時
	Updated time.Time `json:"updated"` // 最終更新日時
	Deleted bool      `json:"deleted"` // 削除フラグ (論理削除用)
}

// Server はアプリケーションの依存関係（設定や状態）を保持
type Server struct {
	taskFilePath string // タスクを永続化するJSONファイルのパス
}

// NewServer は新しいServerインスタンスを生成するコンストラクタ
// taskFilePathには、タスクデータを保存するファイルのパスを指定する。
func NewServer(taskFilePath string) *Server {
	return &Server{taskFilePath: taskFilePath}
}

// helloHandler は動作確認用のエンドポイント
// "Hello World!" というメッセージを含むJSONを返却する。
func (s *Server) helloHandler(w http.ResponseWriter, r *http.Request) {
	response := Response{
		Message: "Hello World!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// postTasksHandler は新しいタスクを作成するリクエストを処理
// HTTP POSTリクエストを受け取り、JSONファイルに新しいタスクを追加する。
func (s *Server) postTasksHandler(w http.ResponseWriter, r *http.Request) {

	// リクエストボディからタスク情報をデコード
	var reqTask Task
	if err := json.NewDecoder(r.Body).Decode(&reqTask); err != nil {
		http.Error(w, "リクエストボディの解析に失敗しました: "+err.Error(), http.StatusBadRequest)
		return
	}

	// 既存のタスク一覧をファイルから読み込む
	var tasks []Task
	file, err := os.ReadFile(s.taskFilePath)
	// ファイルが存在しないエラーは初回リクエストとして許容し、それ以外の読み込みエラーは500を返す
	if err != nil && !os.IsNotExist(err) {
		http.Error(w, "タスクファイルの読み込みに失敗しました: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ファイルに中身があれば（初回リクエストでなければ）、JSONをデコードしてtasksスライスに格納
	if err == nil && len(file) > 0 {
		if err := json.Unmarshal(file, &tasks); err != nil {
			http.Error(w, "タスクデータの解析に失敗しました: "+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// 新しいタスクのIDを決定
	newId := 1
	if len(tasks) > 0 {
		// 既存タスクがあれば、最後のタスクのIDに1を加えることでIDを採番
		newId = tasks[len(tasks)-1].Id + 1
	}

	// 新しいタスクを作成
	now := time.Now()
	newTask := Task{
		Id:      newId,
		Name:    reqTask.Name,
		Status:  0,
		Created: now,
		Updated: now,
		Deleted: false,
	}

	// タスク一覧に新しいタスクを追加
	tasks = append(tasks, newTask)

	// 更新されたタスク一覧を可読性の高いJSON形式（インデント付き）に変換
	tasksJSON, err := json.MarshalIndent(tasks, "", "  ")
	if err != nil {
		http.Error(w, "タスクデータのJSON変換に失敗しました: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// ファイルに書き込み（パーミッションは0644）
	if err := os.WriteFile(s.taskFilePath, tasksJSON, 0666); err != nil {
		http.Error(w, "タスクファイルの書き込みに失敗しました: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// レスポンスとして作成されたタスク情報を返す
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(newTask)
}

// main はアプリケーションのエントリーポイント
func main() {

	// 本番用のファイルパスでサーバーを初期化
	server := NewServer("data/tasks.json")

	// URLパスとハンドラ関数をマッピング
	http.HandleFunc("/", server.helloHandler)
	http.HandleFunc("/api/v1/tasks", server.postTasksHandler)

	// サーバーをポート8080で起動
	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
