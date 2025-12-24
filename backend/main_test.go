package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHelloHandler(t *testing.T) {
	// 1. テスト用のHTTPリクエストを作成 (GET /)
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// 2. レスポンスを記録するためのRecorderを作成
	rr := httptest.NewRecorder()

	// 3. テスト対象のハンドラーを実行
	handler := http.HandlerFunc(helloHandler)
	handler.ServeHTTP(rr, req)

	// 4. ステータスコードの検証
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// 5. Content-Typeヘッダーの検証
	expectedContentType := "application/json"
	if contentType := rr.Header().Get("Content-Type"); contentType != expectedContentType {
		t.Errorf("handler returned wrong content type: got %v want %v",
			contentType, expectedContentType)
	}

	// 6. レスポンスボディの検証
	var response Response
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatalf("could not decode response: %v", err)
	}
	expectedMessage := "Hello World!"
	if response.Message != expectedMessage {
		t.Errorf("handler returned unexpected body: got %v want %v",
			response.Message, expectedMessage)
	}
}
