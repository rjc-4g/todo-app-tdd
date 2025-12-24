package main

import (
	"os"
	"path/filepath"
	"encoding/json"
	"reflect"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// getTasksHandler のテスト
func TestGetTasksHandler(t *testing.T) {
	// 1. テスト用のデータと一時ファイルを作成
	expectedTasks := []Task{
		{ID: 1, Name: "タスク01", Status: 1, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
		{ID: 2, Name: "タスク02", Status: 2, Created: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Deleted: false},
		{ID: 3, Name: "タスク03", Status: 3, Created: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Deleted: false},
	}
	tasksJSON, err := json.Marshal(expectedTasks)
	if err != nil {
		t.Fatalf("テストデータのJSON変換に失敗: %v", err)
	}

	// t.TempDir() によりテスト終了時に自動でクリーンアップされる一時ディレクトリを作成
	tempDir := t.TempDir()
	tempFilePath := filepath.Join(tempDir, "test_tasks.json")
	if err := os.WriteFile(tempFilePath, tasksJSON, 0666); err != nil {
		t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
	}

	// 2. テスト対象のサーバーを一時ファイルパスで初期化
	server := NewServer(tempFilePath)

	// 3. テスト用のHTTPリクエストを作成
	req, err := http.NewRequest("GET", "/api/v1/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 4. レスポンスを記録するためのRecorderを作成
	rr := httptest.NewRecorder()

	// 5. テスト対象のハンドラーを実行
	handler := http.HandlerFunc(server.getTasksHandler)
	handler.ServeHTTP(rr, req)

	// 6. ステータスコードの検証
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusOK)
	}

	// 7. レスポンスボディの検証
	var tasks []Task
	if err := json.NewDecoder(rr.Body).Decode(&tasks); err != nil {
		t.Fatalf("レスポンスボディのデコードに失敗: %v", err)
	}
	if !reflect.DeepEqual(tasks, expectedTasks) {
		t.Errorf("レスポンスボディのテスト失敗\n実測値: %v\n期待値: %v", tasks, expectedTasks)
	}
}
