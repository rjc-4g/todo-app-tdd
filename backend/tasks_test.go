package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// テスト用の一時ファイルパス
const testFilePath = "tasks_test.json"

// テストの事前準備と後片付け
func setup() {
	// テスト開始前にファイルを空にする、または削除する
	os.WriteFile(testFilePath, []byte("[]"), 0644)
}

func teardown() {
	// テスト終了後にファイルを削除する
	os.Remove(testFilePath)
}

// 1. 新規作成 (POST) のテスト
func TestCreateTask(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()
	
	task := map[string]interface{}{
		"title": "Goのテストを書く",
	}
	body, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	assert.Contains(t, w.Body.String(), "Goのテストを書く")
}

// 2. 一覧表示 (GET) のテスト
func TestGetTasks(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	// 配列として返ってくることを期待
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
}

// 3. ステータス更新 (PATCH) のテスト
func TestUpdateTaskStatus(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()
	
	// ID: 1 のタスクのstatusを 1 に更新するリクエスト
	req, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// 4. 削除フラグ更新 (DELETE) のテスト
func TestDeleteTask(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// ID: 1 のタスクを削除（deleted=true）するリクエスト
	req, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}