package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetTasks(t *testing.T) {
	// 1. テスト用のリクエスト作成
	req, err := http.NewRequest("GET", "/api/v1/tasks", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 2. レスポンスを記録するための recorder
	rr := httptest.NewRecorder()
	
	// 3. ハンドラーの実行（まだ中身は空、あるいは未定義なので失敗するはず）
	handler := http.HandlerFunc(getTasksHandler)
	handler.ServeHTTP(rr, req)

	// 4. ステータスコードの検証
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 5. 本来はここで期待するJSON（[]）が返ってきているかチェックする
}