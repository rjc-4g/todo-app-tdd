package main

import (
	"encoding/json"
	"net/http"
)

// 一覧表示（Read）のハンドラー：まずはテストを通すために空を返す
func getTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks := []Task{} // 本来はJSONファイルから読み込むが、今は空
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

// TODO: 他のハンドラーも同様に定義していく
// func createTasksHandler...
// func updateTaskHandler...
// func deleteTaskHandler...

func main() {
    // ルーティング設定
    http.HandleFunc("/api/v1/tasks", func(w http.ResponseWriter, r *http.Request) {
        switch r.Method {
        case http.MethodGet:
            getTasksHandler(w, r)
        case http.MethodPost:
            // createTasksHandler(w, r)
        default:
            w.WriteHeader(http.StatusMethodNotAllowed)
        }
    })
    // サーバー起動処理...
}