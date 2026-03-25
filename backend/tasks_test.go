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
	// テスト用にdbFileNameを変更
	dbFileName = testFilePath
	// テスト開始前にファイルを空にする、または削除する
	os.WriteFile(testFilePath, []byte("[]"), 0644)
}

func teardown() {
	// テスト終了後にファイルを削除する
	os.Remove(testFilePath)
	// dbFileNameを実データ用に戻す
	dbFileName = "tasks.json"
}

// ========== CreateTask のテスト ==========

// 1-1. 新規作成 (POST) のテスト - 正常系
func TestCreateTask_Success(t *testing.T) {
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
	
	// IDが割り当てられていることを確認
	var responseTask Task
	json.Unmarshal(w.Body.Bytes(), &responseTask)
	assert.Equal(t, 1, responseTask.ID)
	assert.Equal(t, 0, responseTask.Status)
	assert.Equal(t, false, responseTask.Deleted)
}

// 1-2. 新規作成 (POST) のテスト - 異常系（無効なリクエストボディ）
func TestCreateTask_InvalidBody(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	assert.Contains(t, w.Body.String(), "Invalid request body")
}

// 1-3. 複数のタスク作成 - 正常系
func TestCreateTask_Multiple(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// 最初のタスク作成
	task1 := map[string]interface{}{
		"title": "タスク1",
	}
	body1, _ := json.Marshal(task1)
	req1, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// 2番目のタスク作成
	task2 := map[string]interface{}{
		"title": "タスク2",
	}
	body2, _ := json.Marshal(task2)
	req2, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// IDが順序よく割り当てられていることを確認
	var task1Response Task
	json.Unmarshal(w1.Body.Bytes(), &task1Response)
	assert.Equal(t, 1, task1Response.ID)

	var task2Response Task
	json.Unmarshal(w2.Body.Bytes(), &task2Response)
	assert.Equal(t, 2, task2Response.ID)
}

// ========== GetTasks のテスト ==========

// 2-1. 一覧表示 (GET) のテスト - 正常系（タスクあり）
func TestGetTasks_WithTasks(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// まずタスクを作成
	task := map[string]interface{}{
		"title": "テストタスク",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// タスク一覧を取得
	req2, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Equal(t, "application/json", w2.Header().Get("Content-Type"))
	assert.Contains(t, w2.Body.String(), "テストタスク")

	// JSONとして正しくパースできることを確認
	var tasks []Task
	err := json.Unmarshal(w2.Body.Bytes(), &tasks)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(tasks))
}

// 2-2. 一覧表示 (GET) のテスト - 正常系（タスクなし）
func TestGetTasks_Empty(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "application/json", w.Header().Get("Content-Type"))
	assert.Equal(t, "[]", w.Body.String())
}

// 2-3. 一覧表示 (GET) のテスト - 削除されたタスクは含まれるが削除フラグが立っている
func TestGetTasks_IncludeDeleted(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// タスクを作成
	task := map[string]interface{}{
		"title": "削除されるタスク",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// タスクを削除（deleted=true）
	req2, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// 一覧を取得（deleted=trueでも返ってくる）
	req3, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)

	var tasks []Task
	json.Unmarshal(w3.Body.Bytes(), &tasks)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, true, tasks[0].Deleted)
}

// ========== UpdateTask のテスト ==========

// 3-1. ステータス更新 (PATCH) のテスト - 正常系
func TestUpdateTask_Success(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// まずタスクを作成
	task := map[string]interface{}{
		"title": "更新テストタスク",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// status を 1 に更新
	updatePayload := map[string]interface{}{
		"status": 1,
	}
	updateBody, _ := json.Marshal(updatePayload)
	req2, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer(updateBody))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// レスポンスのstatusが1になっていることを確認
	var updatedTask Task
	json.Unmarshal(w2.Body.Bytes(), &updatedTask)
	assert.Equal(t, 1, updatedTask.Status)
}

// 3-2. ステータス更新 (PATCH) のテスト - 異常系（存在しないタスク）
func TestUpdateTask_NotFound(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// 存在しないタスク（ID: 999）を更新しようとする
	updatePayload := map[string]interface{}{
		"status": 1,
	}
	updateBody, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest("PATCH", "/api/v1/tasks/999", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

// 3-3. ステータス更新 (PATCH) のテスト - 異常系（無効なリクエストボディ）
func TestUpdateTask_InvalidBody(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// まずタスクを作成
	task := map[string]interface{}{
		"title": "テストタスク",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 無効なボディで更新リクエスト
	req2, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer([]byte("invalid json")))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
	assert.Contains(t, w2.Body.String(), "Invalid request body")
}

// 3-4. ステータス更新 (PATCH) のテスト - 正常系（0 → 1 → 0 の変更）
func TestUpdateTask_StatusToggle(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// タスク作成
	task := map[string]interface{}{
		"title": "トグルテスト",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// status を 1 に更新
	updatePayload1 := map[string]interface{}{
		"status": 1,
	}
	updateBody1, _ := json.Marshal(updatePayload1)
	req1, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer(updateBody1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	var updatedTask1 Task
	json.Unmarshal(w1.Body.Bytes(), &updatedTask1)
	assert.Equal(t, 1, updatedTask1.Status)

	// status を 0 に戻す
	updatePayload2 := map[string]interface{}{
		"status": 0,
	}
	updateBody2, _ := json.Marshal(updatePayload2)
	req2, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer(updateBody2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	var updatedTask2 Task
	json.Unmarshal(w2.Body.Bytes(), &updatedTask2)
	assert.Equal(t, 0, updatedTask2.Status)
}

// ========== DeleteTask のテスト ==========

// 4-1. 削除フラグ更新 (DELETE) のテスト - 正常系
func TestDeleteTask_Success(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// まずタスクを作成
	task := map[string]interface{}{
		"title": "削除テストタスク",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// タスクを削除
	req2, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)

	// レスポンスのdeletedがtrueになっていることを確認
	var deletedTask Task
	json.Unmarshal(w2.Body.Bytes(), &deletedTask)
	assert.Equal(t, true, deletedTask.Deleted)
}

// 4-2. 削除フラグ更新 (DELETE) のテスト - 異常系（存在しないタスク）
func TestDeleteTask_NotFound(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// 存在しないタスク（ID: 999）を削除しようとする
	req, _ := http.NewRequest("DELETE", "/api/v1/tasks/999", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Task not found")
}

// 4-3. 削除フラグ更新 (DELETE) のテスト - 正常系（複数のタスクがある場合）
func TestDeleteTask_MultipleTasksPartialDelete(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// タスク1を作成
	task1 := map[string]interface{}{
		"title": "タスク1",
	}
	body1, _ := json.Marshal(task1)
	req1, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body1))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)

	// タスク2を作成
	task2 := map[string]interface{}{
		"title": "タスク2",
	}
	body2, _ := json.Marshal(task2)
	req2, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body2))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// タスク1を削除
	req3, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	// 一覧を取得すると、両方のタスクが返ってくるが、タスク1はdeletedがtrue
	req4, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)

	var tasks []Task
	json.Unmarshal(w4.Body.Bytes(), &tasks)
	assert.Equal(t, 2, len(tasks))
	assert.Equal(t, true, tasks[0].Deleted)
	assert.Equal(t, false, tasks[1].Deleted)
}

// 4-4. 削除フラグ更新 (DELETE) のテスト - 正常系（ファイル操作確認）
func TestDeleteTask_FileUpdate(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// タスク作成
	task := map[string]interface{}{
		"title": "ファイル更新テスト",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// タスク削除
	req2, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	// ファイルを直接読み込んで確認
	data, _ := os.ReadFile(testFilePath)
	var tasks []Task
	json.Unmarshal(data, &tasks)
	assert.Equal(t, 1, len(tasks))
	assert.Equal(t, true, tasks[0].Deleted)
}

// ========== CORS ミドルウェアのテスト ==========

// 5-1. CORS対応 - OPTIONSリクエスト
func TestCORS_OptionsRequest(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	req, _ := http.NewRequest("OPTIONS", "/api/v1/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "GET")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "DELETE")
}

// 5-2. CORS対応 - CORSヘッダーが含まれている
func TestCORS_HeadersPresent(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Headers"))
}

// ========== エッジケースのテスト ==========

// 6-1. 不正なJSONを含むGETレスポンス処理
func TestGetTasks_CorruptedFile(t *testing.T) {
	setup()
	defer teardown()

	// ファイルに不正なJSONを書き込む
	os.WriteFile(testFilePath, []byte("{ invalid json"), 0644)
	defer teardown()

	router := setupRouter()

	// 不正なファイルからの読み込みをテスト
	// 実装上、getTasks関数内でアンマーシャルエラーは発生しないため、エラーはない
	// ただし、データバリデーションは行っていないため、正常に返される
	req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)
	assert.Equal(t, http.StatusOK, w.Code)
}

// 6-2. CreateTask - titleが空の場合
func TestCreateTask_EmptyTitle(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	task := map[string]interface{}{
		"title": "",
	}
	body, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask Task
	json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.Equal(t, "", createdTask.Title)
}

// 6-3. CreateTask - statusが指定されている場合
func TestCreateTask_WithStatus(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	task := map[string]interface{}{
		"title":  "ステータス指定テスト",
		"status": 1,
	}
	body, _ := json.Marshal(task)

	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var createdTask Task
	json.Unmarshal(w.Body.Bytes(), &createdTask)
	assert.Equal(t, 1, createdTask.Status)
}

// 6-4. UpdateTask - ファイルが存在しない場合
func TestUpdateTask_FileNotFound(t *testing.T) {
	setup()
	os.Remove(testFilePath) // ファイルを削除
	defer teardown()

	router := setupRouter()

	updatePayload := map[string]interface{}{
		"status": 1,
	}
	updateBody, _ := json.Marshal(updatePayload)
	req, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer(updateBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Tasks file not found")
}

// 6-5. DeleteTask - ファイルが存在しない場合
func TestDeleteTask_FileNotFound(t *testing.T) {
	setup()
	os.Remove(testFilePath) // ファイルを削除
	defer teardown()

	router := setupRouter()

	req, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
	assert.Contains(t, w.Body.String(), "Tasks file not found")
}

// 6-6. UpdateTask - 不正なJSONボディ
func TestUpdateTask_MalformedJSON(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// タスク作成
	task := map[string]interface{}{
		"title": "テスト",
	}
	body, _ := json.Marshal(task)
	req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// 不正なJSONで更新
	req2, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer([]byte("{]")))
	req2.Header.Set("Content-Type", "application/json")
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusBadRequest, w2.Code)
}

// 6-7. ルーターの複数メソッド対応確認
func TestRouter_MultipleEndpoints(t *testing.T) {
	setup()
	defer teardown()

	router := setupRouter()

	// POSTでタスク作成
	task := map[string]interface{}{
		"title": "テストタスク",
	}
	body, _ := json.Marshal(task)
	req1, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
	req1.Header.Set("Content-Type", "application/json")
	w1 := httptest.NewRecorder()
	router.ServeHTTP(w1, req1)
	assert.Equal(t, http.StatusCreated, w1.Code)

	// GETで一覧取得
	req2, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)
	assert.Equal(t, http.StatusOK, w2.Code)

	// PATCHでステータス更新
	updatePayload := map[string]interface{}{
		"status": 1,
	}
	updateBody, _ := json.Marshal(updatePayload)
	req3, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer(updateBody))
	req3.Header.Set("Content-Type", "application/json")
	w3 := httptest.NewRecorder()
	router.ServeHTTP(w3, req3)
	assert.Equal(t, http.StatusOK, w3.Code)

	// DELETEで削除フラグ更新
	req4, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
	w4 := httptest.NewRecorder()
	router.ServeHTTP(w4, req4)
	assert.Equal(t, http.StatusOK, w4.Code)
}