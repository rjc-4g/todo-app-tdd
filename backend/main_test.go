package main

import (
	"bytes"
	"os"
	"path/filepath"
	"encoding/json"
	"reflect"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// postTasksHandler のテスト
func TestPostTasksHandler(t *testing.T) {

	// 正常系1: 初回のPOSTリクエスト
	t.Run("FirstPostRequest", func(t *testing.T) {

		// 1. テスト用の一時ファイルを作成
		tempDir := t.TempDir() // テスト終了時に自動でクリーンアップされる一時ディレクトリ
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		newTask := Task{Name: "タスク01"}
		reqBody, _ := json.Marshal(newTask)

		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 4. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 5. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.postTasksHandler)
		handler.ServeHTTP(rr, req)

		// 6. ステータスコードの検証
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusOK)
		}

		// 7. JSONファイルの検証
		fileContent, _ := os.ReadFile(tempFilePath)
		var fileTasks []Task
		if err := json.Unmarshal(fileContent, &fileTasks); err != nil {
			t.Fatalf("ファイル内容のデコードに失敗: %v", err)
		}

		// JSONファイルのタスク数
		actualTaskCount := len(fileTasks)
		expectedTaskCount := 1
		if actualTaskCount != expectedTaskCount {
			t.Errorf("JSONファイルのタスク数が不正\n実測値: %v\n期待値: %v", actualTaskCount, expectedTaskCount)
		}

		// ファイルに追加されたタスク
		addedTask := fileTasks[0]

		// 作成日時の検証: 実行時刻と保存された時刻の差が小さいことを確認（1秒以内）
		if time.Since(addedTask.Created) > 1*time.Second {
			t.Errorf("作成日時が不正: %v", addedTask.Created)
		}

		// 更新日時の検証: 実行時刻と保存された時刻の差が小さいことを確認（1秒以内）
		if time.Since(addedTask.Updated) > 1*time.Second {
			t.Errorf("更新日時が不正: %v", addedTask.Updated)
		}

		// 追加されたタスクの期待値
		expectedTask := Task{Id: 1, Name: "タスク01", Status: 0, Deleted: false}
		expectedTask.Created = addedTask.Created // 作成日時は期待値に実測値を設定
		expectedTask.Updated = addedTask.Updated // 更新日時は期待値に実測値を設定

		// 期待値
		expectedTasks := []Task{expectedTask}

		// JSONファイルの内容と想定されるタスクの比較
		if !reflect.DeepEqual(fileTasks, expectedTasks) {
			t.Errorf("JSONファイルのテスト失敗\n実測値: %v\n期待値: %v", fileTasks, expectedTasks)
		}

		// 8. レスポンスボディの検証
		var task Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗: %v", err)
		}
		if !reflect.DeepEqual(task, expectedTask) {
			t.Errorf("レスポンスボディのテスト失敗\n実測値: %v\n期待値: %v", task, expectedTask)
		}
	})

	// 正常系2: 2回目のPOSTリクエスト
	t.Run("SecondPostRequest", func(t *testing.T) {

		// 1. テスト用のデータと一時ファイルを作成
		initialTasks := []Task{
			{Id: 1, Name: "タスク01", Status: 0, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
		}
		initialJSON, _ := json.Marshal(initialTasks)

		// t.TempDir() によりテスト終了時に自動でクリーンアップされる一時ディレクトリを作成
		tempDir := t.TempDir()
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")
		if err := os.WriteFile(tempFilePath, initialJSON, 0666); err != nil {
			t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
		}

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		newTask := Task{Name: "タスク02"}
		reqBody, _ := json.Marshal(newTask)

		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 4. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 5. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.postTasksHandler)
		handler.ServeHTTP(rr, req)

		// 6. ステータスコードの検証
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusOK)
		}

		// 7. JSONファイルの検証
		fileContent, _ := os.ReadFile(tempFilePath)
		var fileTasks []Task
		if err := json.Unmarshal(fileContent, &fileTasks); err != nil {
			t.Fatalf("ファイル内容のデコードに失敗: %v", err)
		}

		// JSONファイルのタスク数
		actualTaskCount := len(fileTasks)
		expectedTaskCount := 2
		if actualTaskCount != expectedTaskCount {
			t.Errorf("JSONファイルのタスク数が不正\n実測値: %v\n期待値: %v", actualTaskCount, expectedTaskCount)
		}

		// ファイルに追加されたタスク
		addedTask := fileTasks[1]

		// 作成日時の検証: 実行時刻と保存された時刻の差が小さいことを確認（1秒以内）
		if time.Since(addedTask.Created) > 1*time.Second {
			t.Errorf("作成日時が不正: %v", addedTask.Created)
		}

		// 更新日時の検証: 実行時刻と保存された時刻の差が小さいことを確認（1秒以内）
		if time.Since(addedTask.Updated) > 1*time.Second {
			t.Errorf("更新日時が不正: %v", addedTask.Updated)
		}

		// 追加されたタスクの期待値
		expectedTask := Task{Id: 2, Name: "タスク02", Status: 0, Deleted: false}
		expectedTask.Created = addedTask.Created // 作成日時は期待値に実測値を設定
		expectedTask.Updated = addedTask.Updated // 更新日時は期待値に実測値を設定

		// 期待値
		expectedTasks := append(initialTasks, expectedTask)

		// JSONファイルの内容と想定されるタスクの比較
		if !reflect.DeepEqual(fileTasks, expectedTasks) {
			t.Errorf("JSONファイルのテスト失敗\n実測値: %v\n期待値: %v", fileTasks, expectedTasks)
		}

		// 8. レスポンスボディの検証
		var task Task
		if err := json.NewDecoder(rr.Body).Decode(&task); err != nil {
			t.Fatalf("レスポンスボディのデコードに失敗: %v", err)
		}
		if !reflect.DeepEqual(task, expectedTask) {
			t.Errorf("レスポンスボディのテスト失敗\n実測値: %v\n期待値: %v", task, expectedTask)
		}
	})

	// 異常系: システムエラー
	t.Run("SystemError", func(t *testing.T) {

		// 1. テスト対象のサーバーを不正な値（ファイルパスでなくディレクトリ）で初期化
		tempDir := t.TempDir()
		server := NewServer(tempFilePath)

		// 2. テスト用のHTTPリクエストを作成
		newTask := Task{Name: "タスク01"}
		reqBody, _ := json.Marshal(newTask)

		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 3. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 4. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.postTasksHandler)
		handler.ServeHTTP(rr, req)

		// 5. ステータスコードの検証
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusInternalServerError)
		}
	})
}

// getTasksHandler のテスト
func TestGetTasksHandler(t *testing.T) {

	// 正常系1: 一覧取得（タスクあり）
	t.Run("GetTasksNonEmpty", func(t *testing.T) {

		// 1. テスト用のデータと一時ファイルを作成
		expectedTasks := []Task{
			{Id: 1, Name: "タスク01", Status: 0, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 2, Name: "タスク02", Status: 1, Created: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 3, Name: "タスク03", Status: 1, Created: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Deleted: true},
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
		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		req.Header.Set("Content-Type", "application/json")

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
	})

	// 正常系2: 一覧取得（タスクなし）
	t.Run("GetTasksEmpty", func(t *testing.T) {

		// 1. 空の一時ファイルを作成
		tempDir := t.TempDir()
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")
		if err := os.WriteFile(tempFilePath, []byte(""), 0666); err != nil {
			t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
		}

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		req.Header.Set("Content-Type", "application/json")

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
		expectedTasks := []Task{}
		if !reflect.DeepEqual(tasks, expectedTasks) {
			t.Errorf("レスポンスボディのテスト失敗\n実測値: %v\n期待値: %v", tasks, expectedTasks)
		}
	})

	// 異常系: システムエラー
	t.Run("SystemError", func(t *testing.T) {

		// 1. テスト対象のサーバーを不正な値（ファイルパスでなくディレクトリ）で初期化
		tempDir := t.TempDir()
		server := NewServer(tempFilePath)

		// 2. テスト用のHTTPリクエストを作成
		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		req.Header.Set("Content-Type", "application/json")

		// 3. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 4. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.getTasksHandler)
		handler.ServeHTTP(rr, req)

		// 5. ステータスコードの検証
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusInternalServerError)
		}
	})
}

// patchTasksHandler のテスト
func TestPatchTasksHandler(t *testing.T) {

	// 正常系1: 0から1へ更新
	t.Run("UpdateStatusFrom0To1", func(t *testing.T) {

		// 1. テスト用のデータと一時ファイルを作成
		initialTasks := []Task{
			{Id: 1, Name: "タスク01", Status: 0, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 2, Name: "タスク02", Status: 1, Created: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 3, Name: "タスク03", Status: 1, Created: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Deleted: true},
		}
		initialJSON, err := json.Marshal(initialTasks)
		if err != nil {
			t.Fatalf("テストデータのJSON変換に失敗: %v", err)
		}

		// t.TempDir() によりテスト終了時に自動でクリーンアップされる一時ディレクトリを作成
		tempDir := t.TempDir()
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")
		if err := os.WriteFile(tempFilePath, initialJSON, 0666); err != nil {
			t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
		}

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		updatedStatus := 1 // 更新後のステータス
		reqTask := Task{Status: updatedStatus}
		reqBody, _ := json.Marshal(reqTask)

		taskId := 1 // 更新対象のId
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/" + strconv.Itoa(taskId), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 4. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 5. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.patchTasksHandler)
		handler.ServeHTTP(rr, req)

		// 6. ステータスコードの検証
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusOK)
		}

		// 7. JSONファイルの検証
		fileContent, _ := os.ReadFile(tempFilePath)
		var fileTasks []Task
		if err := json.Unmarshal(fileContent, &fileTasks); err != nil {
			t.Fatalf("ファイル内容のデコードに失敗: %v", err)
		}

		// 更新対象のタスク
		updatedTask := fileTasks[taskId - 1]

		// 更新日時の検証: 実行時刻と保存された時刻の差が小さいことを確認（1秒以内）
		if time.Since(updatedTask.Updated) > 1*time.Second {
			t.Errorf("更新日時が不正: %v", updatedTask.Updated)
		}

		// 更新されたタスクの期待値
		expectedTask := initialTasks[taskId - 1]
		expectedTask.Status = updatedStatus
		expectedTask.Updated = updatedTask.Updated // 更新日時は期待値に実測値を設定

		// JSONファイルの内容と想定されるタスクの比較
		if !reflect.DeepEqual(fileTasks, expectedTasks) {
			t.Errorf("JSONファイルのテスト失敗\n実測値: %v\n期待値: %v", fileTasks, expectedTasks)
		}
	})

	// 正常系2: 1から0へ更新
	t.Run("UpdateStatusFrom1To0", func(t *testing.T) {

		// 1. テスト用のデータと一時ファイルを作成
		initialTasks := []Task{
			{Id: 1, Name: "タスク01", Status: 0, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 2, Name: "タスク02", Status: 1, Created: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 3, Name: "タスク03", Status: 1, Created: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Deleted: true},
		}
		initialJSON, err := json.Marshal(initialTasks)
		if err != nil {
			t.Fatalf("テストデータのJSON変換に失敗: %v", err)
		}

		// t.TempDir() によりテスト終了時に自動でクリーンアップされる一時ディレクトリを作成
		tempDir := t.TempDir()
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")
		if err := os.WriteFile(tempFilePath, initialJSON, 0666); err != nil {
			t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
		}

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		updatedStatus := 0 // 更新後のステータス
		reqTask := Task{Status: updatedStatus}
		reqBody, _ := json.Marshal(reqTask)

		taskId := 2 // 更新対象のId
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/" + strconv.Itoa(taskId), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 4. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 5. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.patchTasksHandler)
		handler.ServeHTTP(rr, req)

		// 6. ステータスコードの検証
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusOK)
		}

		// 7. JSONファイルの検証
		fileContent, _ := os.ReadFile(tempFilePath)
		var fileTasks []Task
		if err := json.Unmarshal(fileContent, &fileTasks); err != nil {
			t.Fatalf("ファイル内容のデコードに失敗: %v", err)
		}

		// 更新対象のタスク
		updatedTask := fileTasks[taskId - 1]

		// 更新日時の検証: 実行時刻と保存された時刻の差が小さいことを確認（1秒以内）
		if time.Since(updatedTask.Updated) > 1*time.Second {
			t.Errorf("更新日時が不正: %v", updatedTask.Updated)
		}

		// 更新されたタスクの期待値
		expectedTask := initialTasks[taskId - 1]
		expectedTask.Status = updatedStatus
		expectedTask.Updated = updatedTask.Updated // 更新日時は期待値に実測値を設定

		// JSONファイルの内容と想定されるタスクの比較
		if !reflect.DeepEqual(fileTasks, expectedTasks) {
			t.Errorf("JSONファイルのテスト失敗\n実測値: %v\n期待値: %v", fileTasks, expectedTasks)
		}
	})

	// 異常系: 存在しないIdの指定
	t.Run("InvalidTaskId", func(t *testing.T) {

		// 1. テスト用のデータと一時ファイルを作成
		initialTasks := []Task{
			{Id: 1, Name: "タスク01", Status: 0, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 2, Name: "タスク02", Status: 1, Created: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 3, Name: "タスク03", Status: 1, Created: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Deleted: true},
		}
		initialJSON, err := json.Marshal(initialTasks)
		if err != nil {
			t.Fatalf("テストデータのJSON変換に失敗: %v", err)
		}

		// t.TempDir() によりテスト終了時に自動でクリーンアップされる一時ディレクトリを作成
		tempDir := t.TempDir()
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")
		if err := os.WriteFile(tempFilePath, initialJSON, 0666); err != nil {
			t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
		}

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		updatedStatus := 1 // 更新後のステータス
		reqTask := Task{Status: updatedStatus}
		reqBody, _ := json.Marshal(reqTask)

		taskId := 4 // 存在しないId
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/" + strconv.Itoa(taskId), bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 4. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 5. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.patchTasksHandler)
		handler.ServeHTTP(rr, req)

		// 6. ステータスコードの検証
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusNotFound)
		}

		// 7. JSONファイルの検証
		fileContent, _ := os.ReadFile(tempFilePath)
		var fileTasks []Task
		if err := json.Unmarshal(fileContent, &fileTasks); err != nil {
			t.Fatalf("ファイル内容のデコードに失敗: %v", err)
		}

		// JSONファイルの内容が更新されていないことを確認
		if !reflect.DeepEqual(fileTasks, initialTasks) {
			t.Errorf("JSONファイルのテスト失敗\n実測値: %v\n期待値: %v", fileTasks, initialTasks)
		}
	})

	// 異常系: システムエラー
	t.Run("SystemError", func(t *testing.T) {

		// 1. テスト対象のサーバーを不正な値（ファイルパスでなくディレクトリ）で初期化
		tempDir := t.TempDir()
		server := NewServer(tempFilePath)

		// 2. テスト用のHTTPリクエストを作成
		reqTask := Task{Status: 1}
		reqBody, _ := json.Marshal(reqTask)
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", bytes.NewBuffer(reqBody))
		req.Header.Set("Content-Type", "application/json")

		// 3. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 4. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.patchTasksHandler)
		handler.ServeHTTP(rr, req)

		// 5. ステータスコードの検証
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusInternalServerError)
		}
	})
}

// deleteTasksHandler のテスト
func TestDeleteTasksHandler(t *testing.T) {

	// 正常系: タスクの論理削除
	t.Run("DeleteTask", func(t *testing.T) {

		// 1. テスト用のデータと一時ファイルを作成
		initialTasks := []Task{
			{Id: 1, Name: "タスク01", Status: 0, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 2, Name: "タスク02", Status: 1, Created: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 3, Name: "タスク03", Status: 1, Created: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Deleted: true},
		}
		initialJSON, err := json.Marshal(initialTasks)
		if err != nil {
			t.Fatalf("テストデータのJSON変換に失敗: %v", err)
		}

		// t.TempDir() によりテスト終了時に自動でクリーンアップされる一時ディレクトリを作成
		tempDir := t.TempDir()
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")
		if err := os.WriteFile(tempFilePath, initialJSON, 0666); err != nil {
			t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
		}

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		taskId := 1 // 削除対象のId
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/" + strconv.Itoa(taskId), nil)
		req.Header.Set("Content-Type", "application/json")

		// 4. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 5. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.deleteTasksHandler)
		handler.ServeHTTP(rr, req)

		// 6. ステータスコードの検証
		if status := rr.Code; status != http.StatusOK {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusOK)
		}

		// 7. JSONファイルの検証
		fileContent, _ := os.ReadFile(tempFilePath)
		var fileTasks []Task
		if err := json.Unmarshal(fileContent, &fileTasks); err != nil {
			t.Fatalf("ファイル内容のデコードに失敗: %v", err)
		}

		// 削除対象のタスク
		deletedTask := fileTasks[taskId - 1]

		// 更新日時の検証: 実行時刻と保存された時刻の差が小さいことを確認（1秒以内）
		if time.Since(deletedTask.Updated) > 1*time.Second {
			t.Errorf("更新日時が不正: %v", deletedTask.Updated)
		}

		// 削除されたタスクの期待値
		expectedTask := initialTasks[taskId - 1]
		expectedTask.Deleted = true
		expectedTask.Updated = updatedTask.Updated // 更新日時は期待値に実測値を設定

		// JSONファイルの内容と想定されるタスクの比較
		if !reflect.DeepEqual(fileTasks, expectedTasks) {
			t.Errorf("JSONファイルのテスト失敗\n実測値: %v\n期待値: %v", fileTasks, expectedTasks)
		}
	})

	// 異常系: 存在しないIdの指定
	t.Run("InvalidTaskId", func(t *testing.T) {

		// 1. テスト用のデータと一時ファイルを作成
		initialTasks := []Task{
			{Id: 1, Name: "タスク01", Status: 0, Created: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 1, 10, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 2, Name: "タスク02", Status: 1, Created: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 2, 11, 0, 0, 0, time.UTC), Deleted: false},
			{Id: 3, Name: "タスク03", Status: 1, Created: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Updated: time.Date(2025, 12, 3, 12, 0, 0, 0, time.UTC), Deleted: true},
		}
		initialJSON, err := json.Marshal(initialTasks)
		if err != nil {
			t.Fatalf("テストデータのJSON変換に失敗: %v", err)
		}

		// t.TempDir() によりテスト終了時に自動でクリーンアップされる一時ディレクトリを作成
		tempDir := t.TempDir()
		tempFilePath := filepath.Join(tempDir, "test_tasks.json")
		if err := os.WriteFile(tempFilePath, initialJSON, 0666); err != nil {
			t.Fatalf("一時ファイルへの書き込みに失敗: %v", err)
		}

		// 2. テスト対象のサーバーを一時ファイルパスで初期化
		server := NewServer(tempFilePath)

		// 3. テスト用のHTTPリクエストを作成
		taskId := 4 // 存在しないId
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/" + strconv.Itoa(taskId), nil)
		req.Header.Set("Content-Type", "application/json")

		// 4. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 5. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.deleteTasksHandler)
		handler.ServeHTTP(rr, req)

		// 6. ステータスコードの検証
		if status := rr.Code; status != http.StatusNotFound {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusNotFound)
		}

		// 7. JSONファイルの検証
		fileContent, _ := os.ReadFile(tempFilePath)
		var fileTasks []Task
		if err := json.Unmarshal(fileContent, &fileTasks); err != nil {
			t.Fatalf("ファイル内容のデコードに失敗: %v", err)
		}

		// JSONファイルの内容が更新されていないことを確認
		if !reflect.DeepEqual(fileTasks, initialTasks) {
			t.Errorf("JSONファイルのテスト失敗\n実測値: %v\n期待値: %v", fileTasks, initialTasks)
		}
	})

	// 異常系: システムエラー
	t.Run("SystemError", func(t *testing.T) {

		// 1. テスト対象のサーバーを不正な値（ファイルパスでなくディレクトリ）で初期化
		tempDir := t.TempDir()
		server := NewServer(tempFilePath)

		// 2. テスト用のHTTPリクエストを作成
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", nil)
		req.Header.Set("Content-Type", "application/json")

		// 3. レスポンスを記録するためのRecorderを作成
		rr := httptest.NewRecorder()

		// 4. テスト対象のハンドラーを実行
		handler := http.HandlerFunc(server.deleteTasksHandler)
		handler.ServeHTTP(rr, req)

		// 5. ステータスコードの検証
		if status := rr.Code; status != http.StatusInternalServerError {
			t.Errorf("ステータスコードのテスト失敗\n実測値: %v\n期待値: %v", status, http.StatusInternalServerError)
		}
	})
}
