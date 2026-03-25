package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"
	"unicode/utf8"
)

// TestMain はテスト全体のセットアップと削除を行う。
// テスト用の一時JSONファイルを作成し、ID=123のシードデータを投入してから
// テストを実行し、終了後に一時ファイルを削除する。
func TestMain(m *testing.M) {
	// テスト用の一時ファイルを作成
	tmpFile, err := os.CreateTemp("", "tasks_test_*.json")
	if err != nil {
		panic(err)
	}
	tmpFile.Close()

	// テスト用ファイルパスを本番用から切り替える
	dataFile = tmpFile.Name()

	// ID=123のシードデータを書き込む（更新・削除テストが依存するため）
	now := time.Now()
	seedTasks := []Task{
		{
			ID:      123,
			Name:    "seed task",
			Status:  0,
			Created: now,
			Updated: now,
			Deleted: false,
		},
	}
	data, _ := json.Marshal(seedTasks)
	if err := os.WriteFile(dataFile, data, 0644); err != nil {
		panic(err)
	}

	// ファイルからストアをメモリに読み込む
	if err := loadStore(); err != nil {
		panic(err)
	}

	// テストを実行
	code := m.Run()

	// 一時ファイルを削除
	os.Remove(dataFile)

	os.Exit(code)
}

// POST /api/v1/tasks エンドポイント　テストケース
// 新規タスクを作成する機能をテスト
func TestCreateTask(t *testing.T) {
	// テスト用のタスクデータを準備（nameとstatusのみ指定、IDは自動生成される）
	requestData := map[string]interface{}{
		"name":   "新しいタスク",
		"status": 0,
	}

	// リクエストデータをJSON形式に変換
	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	// POSTリクエストを作成（エンドポイント: /api/v1/tasks）
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// レスポンスを記録するためのRecorderを作成
	rec := httptest.NewRecorder()

	// ハンドラー関数を呼び出してリクエストを処理
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが201 Createdであることを確認
	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	// レスポンスボディをTask構造体に変換
	var response Task
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// IDが自動生成されていることを確認（0より大きい値であることを確認）
	if response.ID <= 0 {
		t.Errorf("Expected task ID to be generated (positive integer), got %d", response.ID)
	}

	// 名前が正しく保存されていることを確認
	expectedName := "新しいタスク"
	if response.Name != expectedName {
		t.Errorf("Expected name %s, got %s", expectedName, response.Name)
	}

	// ステータスが正しく保存されていることを確認
	expectedStatus := 0
	if response.Status != expectedStatus {
		t.Errorf("Expected status %d, got %d", expectedStatus, response.Status)
	}

	// 作成日時が設定されていることを確認
	if response.Created.IsZero() {
		t.Error("Expected created timestamp to be set")
	}

	// 更新日時が設定されていることを確認
	if response.Updated.IsZero() {
		t.Error("Expected updated timestamp to be set")
	}

	// 削除フラグがfalse（有効）であることを確認
	if response.Deleted != false {
		t.Error("Expected deleted to be false")
	}
}

// GET /api/v1/tasks エンドポイント　テストケース
// タスク一覧を取得する機能をテスト
func TestGetTasks(t *testing.T) {
	// GETリクエストを作成（エンドポイント: /api/v1/tasks）
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)

	// レスポンスを記録するためのRecorderを作成
	rec := httptest.NewRecorder()

	// ハンドラー関数を呼び出してリクエストを処理
	handler := http.HandlerFunc(getTasksHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが200 OKであることを確認
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// レスポンスボディをTaskの配列に変換
	var tasks []Task
	err := json.Unmarshal(rec.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 配列が返却されることを確認（空配列でも可、nilは不可）
	if tasks == nil {
		t.Error("Expected tasks array, got nil")
	}
}

// PATCH /api/v1/tasks/{id} エンドポイント　テストケース
// タスクのステータス（0/1）を更新する機能をテスト
func TestUpdateTaskStatus(t *testing.T) {
	// テスト用のタスクIDを設定（int型）
	taskID := 123

	// 更新データを準備（statusを1に更新）
	updateData := map[string]int{
		"status": 1,
	}

	// 更新データをJSON形式に変換
	body, err := json.Marshal(updateData)
	if err != nil {
		t.Fatalf("Failed to marshal update data: %v", err)
	}

	// PATCHリクエストを作成（エンドポイント: /api/v1/tasks/{id}）
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// レスポンスを記録するためのRecorderを作成
	rec := httptest.NewRecorder()

	// ハンドラー関数を呼び出してリクエストを処理（未実装のため失敗する見込み）
	handler := http.HandlerFunc(updateTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが200 OKであることを確認
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// レスポンスボディをTask構造体に変換
	var response Task
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// IDが正しく設定されていることを確認
	if response.ID != taskID {
		t.Errorf("Expected ID %d, got %d", taskID, response.ID)
	}

	// ステータスが1に更新されていることを確認
	if response.Status != 1 {
		t.Errorf("Expected status to be updated to 1, got %d", response.Status)
	}

	// 更新日時が更新されていることを確認
	if response.Updated.IsZero() {
		t.Error("Expected updated timestamp to be set")
	}
}

// DELETE /api/v1/tasks/{id} エンドポイント　テストケース
// タスクを削除する機能をテスト（論理削除：deletedフラグをtrueに更新）
func TestDeleteTask(t *testing.T) {
	// テスト用のタスクIDを設定（int型）
	taskID := 123

	// DELETEリクエストを作成（エンドポイント: /api/v1/tasks/{id}）
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/123", nil)

	// レスポンスを記録するためのRecorderを作成
	rec := httptest.NewRecorder()

	// ハンドラー関数を呼び出してリクエストを処理
	handler := http.HandlerFunc(deleteTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが200 OKであることを確認
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	// レスポンスボディをTask構造体に変換
	var response Task
	err := json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// IDが正しく設定されていることを確認
	if response.ID != taskID {
		t.Errorf("Expected ID %d, got %d", taskID, response.ID)
	}

	// deletedフラグがtrueに更新されていることを確認
	if response.Deleted != true {
		t.Error("Expected deleted to be true")
	}

	// 更新日時が更新されていることを確認
	if response.Updated.IsZero() {
		t.Error("Expected updated timestamp to be set")
	}
}

// POST /api/v1/tasks エンドポイント　エラーハンドリングテストケース
// 不正なJSONを送信した場合のエラーハンドリング
// 不正なJSON形式のリクエストに対して400 Bad Requestが返却されることを確認する
func TestCreateTaskInvalidJSON(t *testing.T) {
	// 不正なJSON形式のリクエストを作成
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	// レスポンスを記録するためのRecorderを作成
	rec := httptest.NewRecorder()

	// ハンドラー関数を呼び出してリクエストを処理
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, rec.Code)
	}
}

// PATCH /api/v1/tasks/{id} エンドポイント　エラーハンドリングテストケース
// 存在しないタスクを更新しようとした場合のエラーハンドリング
// 存在しないタスクIDに対して404 Not Foundが返却されることを確認する
func TestUpdateTaskNotFound(t *testing.T) {
	// 存在しないタスクIDを設定（int型）
	nonExistentID := 99999

	// 更新データを準備
	updateData := map[string]int{
		"status": 1,
	}

	// 更新データをJSON形式に変換
	body, _ := json.Marshal(updateData)

	// PATCHリクエストを作成（存在しないタスクIDを指定）
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/"+strconv.Itoa(nonExistentID), bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	// レスポンスを記録するためのRecorderを作成
	rec := httptest.NewRecorder()

	// ハンドラー関数を呼び出してリクエストを処理
	handler := http.HandlerFunc(updateTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが404 Not Foundであることを確認
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent task, got %d", http.StatusNotFound, rec.Code)
	}
}

// DELETE /api/v1/tasks/{id} エンドポイント　エラーハンドリングテストケース
// 存在しないタスクを削除しようとした場合のエラーハンドリング
// 存在しないタスクIDに対して404 Not Foundが返却されることを確認する
func TestDeleteTaskNotFound(t *testing.T) {
	// 存在しないタスクIDを設定（int型）
	nonExistentID := 99999

	// DELETEリクエストを作成（存在しないタスクIDを指定）
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+strconv.Itoa(nonExistentID), nil)

	// レスポンスを記録するためのRecorderを作成
	rec := httptest.NewRecorder()

	// ハンドラー関数を呼び出してリクエストを処理
	handler := http.HandlerFunc(deleteTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが404 Not Foundであることを確認
	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status %d for non-existent task, got %d", http.StatusNotFound, rec.Code)
	}
}

// POST /api/v1/tasks エンドポイント　エラーハンドリングテストケース
// nameが40文字を超える場合のエラーハンドリング
// varchar(40)の制約を確認します
func TestCreateTaskNameTooLong(t *testing.T) {
	// 41文字のnameを準備（制限を超える）
	longName := "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらりるれろわをんあ"
	if len(longName) <= 40 {
		t.Fatal("Test setup error: name should be longer than 40 characters")
	}

	requestData := map[string]interface{}{
		"name":   longName,
		"status": 0,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for name too long, got %d", http.StatusBadRequest, rec.Code)
	}
}

// POST /api/v1/tasks エンドポイント　エラーハンドリングテストケース
// statusが0,1以外の値の場合のエラーハンドリング
// int(1)の制約（0または1のみ）を確認します
func TestCreateTaskInvalidStatus(t *testing.T) {
	// statusが2の場合（無効な値）
	requestData := map[string]interface{}{
		"name":   "テストタスク",
		"status": 2,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid status, got %d", http.StatusBadRequest, rec.Code)
	}
}

// POST /api/v1/tasks エンドポイント　エラーハンドリングテストケース
// statusが負の数の場合のエラーハンドリング
func TestCreateTaskNegativeStatus(t *testing.T) {
	requestData := map[string]interface{}{
		"name":   "テストタスク",
		"status": -1,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for negative status, got %d", http.StatusBadRequest, rec.Code)
	}
}

// POST /api/v1/tasks エンドポイント　エラーハンドリングテストケース
// nameが空文字の場合のエラーハンドリング（not null 制約）
func TestCreateTaskEmptyName(t *testing.T) {
	requestData := map[string]interface{}{
		"name":   "",
		"status": 0,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty name, got %d", http.StatusBadRequest, rec.Code)
	}
}

// POST /api/v1/tasks エンドポイント　エラーハンドリングテストケース
// nameフィールドが欠如している場合のエラーハンドリング
func TestCreateTaskMissingName(t *testing.T) {
	// nameフィールドがないリクエストデータ
	requestData := map[string]interface{}{
		"status": 0,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for missing name, got %d", http.StatusBadRequest, rec.Code)
	}
}

// POST /api/v1/tasks エンドポイント　エラーハンドリングテストケース
// リクエストボディが空の場合のエラーハンドリング
func TestCreateTaskEmptyBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBufferString(""))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for empty body, got %d", http.StatusBadRequest, rec.Code)
	}
}

// PATCH /api/v1/tasks/{id} エンドポイント　エラーハンドリングテストケース
// 不正なID形式（文字列など）が渡された場合のエラーハンドリング
func TestUpdateTaskInvalidID(t *testing.T) {
	updateData := map[string]int{
		"status": 1,
	}

	body, _ := json.Marshal(updateData)

	// 文字列のIDを指定（数値でない）
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/invalid-id", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(updateTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid ID format, got %d", http.StatusBadRequest, rec.Code)
	}
}

// DELETE /api/v1/tasks/{id} エンドポイント　エラーハンドリングテストケース
// 不正なID形式が渡された場合のエラーハンドリング
func TestDeleteTaskInvalidID(t *testing.T) {
	// 文字列のIDを指定（数値でない）
	req := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/invalid-id", nil)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(deleteTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid ID format, got %d", http.StatusBadRequest, rec.Code)
	}
}

// PATCH /api/v1/tasks/{id} エンドポイント　エラーハンドリングテストケース
// 更新時にstatusが0,1以外の値の場合のエラーハンドリング
func TestUpdateTaskInvalidStatus(t *testing.T) {
	updateData := map[string]int{
		"status": 2, // 無効な値
	}

	body, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(updateTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid status, got %d", http.StatusBadRequest, rec.Code)
	}
}

// PATCH /api/v1/tasks/{id} エンドポイント　エラーハンドリングテストケース
// 更新時にstatusフィールドが欠如している場合のエラーハンドリング
func TestUpdateTaskMissingStatus(t *testing.T) {
	// statusフィールドがない更新データ
	updateData := map[string]interface{}{}

	body, _ := json.Marshal(updateData)

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(updateTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for missing status, got %d", http.StatusBadRequest, rec.Code)
	}
}

// PATCH /api/v1/tasks/{id} エンドポイント　エラーハンドリングテストケース
// 更新時に不正なJSONが送信された場合のエラーハンドリング
func TestUpdateTaskInvalidJSON(t *testing.T) {
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/123", bytes.NewBufferString("{invalid json}"))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(updateTaskHandler)
	handler.ServeHTTP(rec, req)

	// レスポンスのステータスコードが400 Bad Requestであることを確認
	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, rec.Code)
	}
}

// POST /api/v1/tasks エンドポイント　テストケース
// statusを指定せずにタスクを作成するテスト
// デフォルト値0が設定されることを確認
func TestCreateTaskWithoutStatus(t *testing.T) {
	// statusを指定しないリクエストデータ
	requestData := map[string]interface{}{
		"name": "ステータス未指定のタスク",
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response Task
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// デフォルト値0が設定されていることを確認
	if response.Status != 0 {
		t.Errorf("Expected default status 0, got %d", response.Status)
	}
}

// POST /api/v1/tasks エンドポイント　テストケース
// nameが最大長（40文字）の場合のテスト
func TestCreateTaskWithMaxLengthName(t *testing.T) {
	// 40文字のnameを準備（制限内）
	maxLengthName := "あいうえおかきくけこさしすせそたちつてとなにぬねのはひふへほまみむめもやゆよらり"
	if utf8.RuneCountInString(maxLengthName) != 40 {
		t.Fatalf("Test setup error: name should be exactly 40 characters, got %d", utf8.RuneCountInString(maxLengthName))
	}

	requestData := map[string]interface{}{
		"name":   maxLengthName,
		"status": 0,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response Task
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// nameが正しく保存されていることを確認
	if response.Name != maxLengthName {
		t.Errorf("Expected name %s, got %s", maxLengthName, response.Name)
	}
}

// POST /api/v1/tasks エンドポイント　テストケース
// nameが最小長（1文字）の場合のテスト（not null の境界）
func TestCreateTaskWithMinLengthName(t *testing.T) {
	oneCharName := "a"
	requestData := map[string]interface{}{
		"name":   oneCharName,
		"status": 0,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var response Task
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.Name != oneCharName {
		t.Errorf("Expected name %q, got %q", oneCharName, response.Name)
	}
}

// PATCH /api/v1/tasks/{id} エンドポイント　テストケース
// ステータスを1→0に更新するテスト（逆方向の更新）
func TestUpdateTaskStatusFrom1To0(t *testing.T) {
	taskID := 123

	updateData := map[string]int{
		"status": 0, // 1から0に戻す
	}

	body, err := json.Marshal(updateData)
	if err != nil {
		t.Fatalf("Failed to marshal update data: %v", err)
	}

	req := httptest.NewRequest(http.MethodPatch, "/api/v1/tasks/123", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(updateTaskHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var response Task
	err = json.Unmarshal(rec.Body.Bytes(), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if response.ID != taskID {
		t.Errorf("Expected ID %d, got %d", taskID, response.ID)
	}

	// ステータスが0に更新されていることを確認
	if response.Status != 0 {
		t.Errorf("Expected status to be updated to 0, got %d", response.Status)
	}
}

// GET /api/v1/tasks エンドポイント　テストケース
// 空のタスク一覧を取得するテスト
func TestGetTasksEmpty(t *testing.T) {
	// 他のテストで追加されたタスクをクリアして空の状態を確保する
	resetStore()

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(getTasksHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var tasks []Task
	err := json.Unmarshal(rec.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 空配列が返却されることを確認
	if tasks == nil {
		t.Error("Expected empty array, got nil")
	}

	if len(tasks) != 0 {
		t.Errorf("Expected empty array, got array with %d items", len(tasks))
	}
}

// POST /api/v1/tasks と GET /api/v1/tasks エンドポイント　テストケース
// 複数のタスクを作成して一覧取得するテスト
func TestCreateMultipleTasks(t *testing.T) {
	// 複数のタスクを作成
	taskNames := []string{"タスク1", "タスク2", "タスク3"}
	var createdTasks []Task

	for _, name := range taskNames {
		requestData := map[string]interface{}{
			"name":   name,
			"status": 0,
		}

		body, err := json.Marshal(requestData)
		if err != nil {
			t.Fatalf("Failed to marshal task: %v", err)
		}

		req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		rec := httptest.NewRecorder()
		handler := http.HandlerFunc(createTaskHandler)
		handler.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
		}

		var task Task
		err = json.Unmarshal(rec.Body.Bytes(), &task)
		if err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		createdTasks = append(createdTasks, task)
	}

	// 一覧を取得
	req := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(getTasksHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, rec.Code)
	}

	var tasks []Task
	err := json.Unmarshal(rec.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 作成したタスク数以上が存在することを確認（削除されていないタスクのみ）
	if len(tasks) < len(createdTasks) {
		t.Errorf("Expected at least %d tasks, got %d", len(createdTasks), len(tasks))
	}
}

// POST /api/v1/tasks, DELETE /api/v1/tasks/{id}, GET /api/v1/tasks エンドポイント　テストケース
// 削除されたタスクが一覧に含まれないことを確認するテスト
func TestGetTasksExcludesDeleted(t *testing.T) {
	// タスクを作成
	requestData := map[string]interface{}{
		"name":   "削除されるタスク",
		"status": 0,
	}

	body, err := json.Marshal(requestData)
	if err != nil {
		t.Fatalf("Failed to marshal task: %v", err)
	}

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	handler := http.HandlerFunc(createTaskHandler)
	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status %d, got %d", http.StatusCreated, rec.Code)
	}

	var createdTask Task
	err = json.Unmarshal(rec.Body.Bytes(), &createdTask)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// タスクを削除
	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/tasks/"+strconv.Itoa(createdTask.ID), nil)
	deleteRec := httptest.NewRecorder()
	deleteHandler := http.HandlerFunc(deleteTaskHandler)
	deleteHandler.ServeHTTP(deleteRec, deleteReq)

	// 一覧を取得
	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/tasks", nil)
	getRec := httptest.NewRecorder()
	getHandler := http.HandlerFunc(getTasksHandler)
	getHandler.ServeHTTP(getRec, getReq)

	if getRec.Code != http.StatusOK {
		t.Errorf("Expected status %d, got %d", http.StatusOK, getRec.Code)
	}

	var tasks []Task
	err = json.Unmarshal(getRec.Body.Bytes(), &tasks)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 削除されたタスクが一覧に含まれていないことを確認
	for _, task := range tasks {
		if task.ID == createdTask.ID {
			t.Errorf("Deleted task (ID: %d) should not appear in the list", createdTask.ID)
		}
	}
}
