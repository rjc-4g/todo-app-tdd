package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// テルパー: テスト環境のセットアップとクリーンアップ
func setupTest(t *testing.T) func() {
	// テスト用のファイルパスに差し替える
	originalPath := jsonFilePath
	jsonFilePath = "tasks_test.json"

	if _, err := os.Stat(jsonFilePath); err == nil {
		os.Remove(jsonFilePath)
	}

	return func() {
		os.Remove(jsonFilePath)
		jsonFilePath = originalPath
	}
}

func TestTasksAPI(t *testing.T) {
	teardown := setupTest(t)
	defer teardown()

	t.Run("初期状態は空の配列を返す", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		rr := httptest.NewRecorder()
		getTasksHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}
		var tasks []Task
		json.NewDecoder(rr.Body).Decode(&tasks)
		if len(tasks) != 0 {
			t.Errorf("Expected 0 tasks, got %d", len(tasks))
		}
	})

	t.Run("新規タスクを作成できる", func(t *testing.T) {
		taskBody := map[string]string{"name": "Learn TDD"}
		body, _ := json.Marshal(taskBody)
		req, _ := http.NewRequest("POST", "/api/v1/tasks", bytes.NewBuffer(body))
		rr := httptest.NewRecorder()
		createTasksHandler(rr, req)

		if rr.Code != http.StatusCreated {
			t.Errorf("Expected status 201, got %d", rr.Code)
		}

		var createdTask Task
		json.NewDecoder(rr.Body).Decode(&createdTask)
		if createdTask.ID != 1 {
			t.Errorf("Expected ID 1, got %d", createdTask.ID)
		}
		if createdTask.Name != "Learn TDD" {
			t.Errorf("Expected Name 'Learn TDD', got %s", createdTask.Name)
		}
		if createdTask.Status != 0 {
			t.Errorf("Expected Status 0, got %d", createdTask.Status)
		}
		if createdTask.Deleted != false {
			t.Error("Expected Deleted false, got true")
		}
	})

	t.Run("作成したタスクを一覧で取得できる", func(t *testing.T) {
		req, _ := http.NewRequest("GET", "/api/v1/tasks", nil)
		rr := httptest.NewRecorder()
		getTasksHandler(rr, req)

		var tasks []Task
		json.NewDecoder(rr.Body).Decode(&tasks)
		if len(tasks) != 1 {
			t.Fatalf("Expected 1 task, got %d", len(tasks))
		}
		if tasks[0].Name != "Learn TDD" {
			t.Errorf("Expected task name 'Learn TDD', got %s", tasks[0].Name)
		}
	})

	t.Run("タスクのステータスを更新できる (PATCH)", func(t *testing.T) {
		// ID 1 のタスクを更新 (Status 0 -> 1)
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/1", nil)
		rr := httptest.NewRecorder()
		updateTaskHandler(rr, req)

		if rr.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", rr.Code)
		}

		// 更新後の確認
		tasks, _ := readTasks()
		if tasks[0].Status != 1 {
			t.Errorf("Expected status 1, got %d", tasks[0].Status)
		}
	})

	t.Run("タスクを論理削除できる (DELETE)", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/v1/tasks/1", nil)
		rr := httptest.NewRecorder()
		deleteTaskHandler(rr, req)

		if rr.Code != http.StatusNoContent {
			t.Errorf("Expected status 204, got %d", rr.Code)
		}

		// 更新後の確認
		tasks, _ := readTasks()
		if tasks[0].Deleted != true {
			t.Error("Expected Deleted true, got false")
		}
	})

	t.Run("存在しないIDの更新は404を返す", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/999", nil)
		rr := httptest.NewRecorder()
		updateTaskHandler(rr, req)

		if rr.Code != http.StatusNotFound {
			t.Errorf("Expected status 404, got %d", rr.Code)
		}
	})

	t.Run("不正なID形式は400を返す", func(t *testing.T) {
		req, _ := http.NewRequest("PATCH", "/api/v1/tasks/invalid", nil)
		rr := httptest.NewRecorder()
		updateTaskHandler(rr, req)

		if rr.Code != http.StatusBadRequest {
			t.Errorf("Expected status 400, got %d", rr.Code)
		}
	})
}
