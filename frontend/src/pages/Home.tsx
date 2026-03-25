import { useState, useEffect } from "react";
import axios from "axios";
import TaskForm from "../components/TaskForm";
import Button from "../components/Button";
import Checkbox from "../components/Checkbox";

// タスクの型定義
interface Task {
  id: number;
  name: string;
  status: number;
  created: string;
  updated: string;
  deleted: boolean;
}

export default function Home() {
  const [tasks, setTasks] = useState<Task[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false); // フォームの表示状態

  const fetchTasks = async () => {
    try {
      const response = await axios.get("/api/v1/tasks");
      setTasks(response.data.filter((task: Task) => !task.deleted));
    } catch (error) {
      console.error("Failed to fetch tasks:", error);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, []);

  // 新規タスクの作成ロジック
  const handleCreateTask = async (name: string) => {
    try {
      await axios.post("/api/v1/tasks", { name });
      await fetchTasks();
    } finally {
      setShowForm(false); // 作成の成否に関わらず、フォームを閉じる
    }
  };

  const handleToggleStatus = async (id: number) => {
    try {
      await axios.patch(`/api/v1/tasks/${id}`);
      fetchTasks();
    } catch (error) {
      console.error("Failed to update task:", error);
    }
  };

  const handleDelete = async (id: number) => {
    if (!window.confirm("このタスクを削除しますか？")) return;
    try {
      await axios.delete(`/api/v1/tasks/${id}`);
      fetchTasks();
    } catch (error) {
      console.error("Failed to delete task:", error);
    }
  };

  if (loading) return <div className="p-8 text-gray-900 bg-gray-50 min-h-screen text-center">読み込み中...</div>;

  return (
    <div className="min-h-screen bg-gray-50 py-12 px-4 sm:px-6 lg:px-8 text-gray-900 font-sans">
      <div className="max-w-2xl mx-auto">
        
        {/* ヘッダーエリア */}
        <div className="flex items-center justify-center mb-8 relative">
          <Button
            onMouseDown={(e) => e.preventDefault()}
            onClick={() => setShowForm(!showForm)}
            variant={showForm ? "secondary" : "primary"}
            size="lg"
            className="absolute left-0"
          >
            {showForm ? "閉じる" : "新規作成"}
          </Button>
          <h1 className="text-3xl font-extrabold text-gray-900">TODO APP TDD</h1>
        </div>

        {/* タスク一覧 */}
        <div className="space-y-3 mb-6">
          {tasks.length === 0 && !showForm ? (
            <div className="text-center py-10 bg-white rounded-lg border border-dashed border-gray-300 shadow-sm">
              <p className="text-gray-500">タスクがありません。「新規作成」から追加しましょう！</p>
            </div>
          ) : (
            tasks.map((task) => (
              <div
                key={task.id}
                className="flex items-center justify-between p-4 bg-white border border-gray-200 rounded-lg shadow-sm hover:shadow-md transition-shadow"
              >
                <div className="flex items-center gap-4 flex-1">
                  <Checkbox
                    checked={task.status === 1}
                    onChange={() => handleToggleStatus(task.id)}
                  />
                  <span
                    className={`text-lg font-medium transition-all ${
                      task.status === 1 ? "line-through text-gray-400" : "text-gray-900"
                    }`}
                  >
                    {task.name}
                  </span>
                </div>
                <Button
                  onClick={() => handleDelete(task.id)}
                  variant="outline-danger"
                  size="sm"
                >
                  削除
                </Button>
              </div>
            ))
          )}
        </div>

        {/* 一覧の一番下に表示される新規作成フォーム */}
        {showForm && (
          <div className="mt-6 animate-in fade-in slide-in-from-bottom-4 duration-300">
            <TaskForm 
              onSubmit={handleCreateTask} 
              onBlur={() => setShowForm(false)} 
              autoFocus 
            />
          </div>
        )}

      </div>
    </div>
  );
}
