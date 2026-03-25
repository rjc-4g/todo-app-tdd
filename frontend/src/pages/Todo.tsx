import { useState, useEffect } from "react";
import Button from "@/components/Button";
import { Task, ApiTask } from "@/types";

export default function Todo() {
  const [tasks, setTasks] = useState<Task[]>([
    {
      id: Date.now().toString(),
      name: "",
      completed: false,
    }
  ]);

  // コンポーネントマウント時にバックエンドからタスクを取得
  useEffect(() => {
    const fetchTasks = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/tasks`);
        if (response.ok) {
          const data: ApiTask[] = await response.json();
          // バックエンドのレスポンスをフロントエンド形式に変換
          const convertedTasks: Task[] = data
            .filter((task) => !task.deleted)
            .map((task) => ({
              id: task.id.toString(),
              name: task.title,
              completed: task.status === 1,
            }));

          // 常に最後に空のタスク行を追加
          const emptyTask: Task = {
            id: Date.now().toString(),
            name: "",
            completed: false,
          };

          // データがある場合は表示、ない場合は空のタスク行のみ
          if (convertedTasks.length > 0) {
            setTasks([...convertedTasks, emptyTask]);
          } else {
            setTasks([emptyTask]);
          }
        }
      } catch (error) {
        console.error("Failed to fetch tasks:", error);
      }
    };

    fetchTasks();
  }, []);

  // Create: 新しいタスクを追加
  const addTask = async () => {
    // 最後のタスク以外のタスクから、空でないものを取得
    const tasksToSave = tasks.slice(0, -1);
    const lastTask = tasks[tasks.length - 1];

    // 最後のタスクに入力がある場合、APIで登録
    if (lastTask && lastTask.name.trim()) {
      try {
        const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/tasks`, {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
          },
          body: JSON.stringify({
            title: lastTask.name,
            status: lastTask.completed ? 1 : 0,
            deleted: false,
          }),
        });

        if (response.ok) {
          const createdTask: ApiTask = await response.json();

          // 最後のタスクをAPIから返されたタスクに置き換え、新しい空のタスク行を追加
          const newEmptyTask: Task = {
            id: Date.now().toString(),
            name: "",
            completed: false,
          };

          setTasks([
            ...tasksToSave,
            {
              id: createdTask.id.toString(),
              name: createdTask.title,
              completed: createdTask.status === 1,
            },
            newEmptyTask,
          ]);
        }
      } catch (error) {
        console.error("Failed to create task:", error);
      }
    }
  };

  // Delete: タスクを削除
  const deleteTask = async (id: string) => {
    try {
      const taskId = parseInt(id);
      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/tasks/${taskId}`, {
        method: "DELETE",
        headers: {
          "Content-Type": "application/json",
        },
      });

      if (response.ok) {
        setTasks(tasks.filter((task) => task.id !== id));
      } else {
        console.error("Failed to delete task:", response.statusText);
      }
    } catch (error) {
      console.error("Failed to delete task:", error);
    }
  };

  // Update: タスク名を更新
  const updateTaskName = (id: string, name: string) => {
    setTasks(
      tasks.map((task) =>
        task.id === id ? { ...task, name } : task
      )
    );
  };

  // Update: チェック状態を切り替え
  const toggleComplete = async (id: string) => {
    try {
      const taskId = parseInt(id);
      const task = tasks.find((t) => t.id === id);
      if (!task) return;

      // 新しいステータス（0又は1）
      const newStatus = task.completed ? 0 : 1;

      const response = await fetch(`${import.meta.env.VITE_API_BASE_URL}/tasks/${taskId}`, {
        method: "PATCH",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify({
          status: newStatus,
        }),
      });

      if (response.ok) {
        setTasks(
          tasks.map((t) =>
            t.id === id ? { ...t, completed: !t.completed } : t
          )
        );
      } else {
        console.error("Failed to update task:", response.statusText);
      }
    } catch (error) {
      console.error("Failed to update task:", error);
    }
  };

  return (
    <div className="min-h-screen bg-teal-500 p-8">
      <div className="max-w-2xl mx-auto">
        {/* タイトルと新規ボタン */}
        <div className="flex justify-between items-center mb-8">
          <button
            onClick={addTask}
            className="bg-white text-black font-bold py-3 px-6 rounded-lg hover:bg-gray-100 transition"
          >
            新規
          </button>
          <div className="bg-white rounded-lg px-10 py-4">
            <h1 className="text-2xl font-bold text-gray-800">TODO APP TDD</h1>
          </div>
          <div className="w-32" />
        </div>

        {/* タスク一覧 */}
        <div className="space-y-4">
          {tasks.map((task) => (
            <div
              key={task.id}
              className="bg-white rounded-lg p-4 flex items-center gap-4 shadow-md"
            >
              {/* チェックボックス */}
              <input
                type="checkbox"
                checked={task.completed}
                onChange={() => toggleComplete(task.id)}
                className="w-6 h-6 text-teal-500 rounded cursor-pointer accent-teal-500"
              />

              {/* タスク名入力フィールド */}
              <input
                type="text"
                value={task.name}
                onChange={(e) => updateTaskName(task.id, e.target.value)}
                placeholder="タスクを入力してください"
                className="flex-1 px-4 py-2 border-2 border-gray-300 rounded-lg focus:outline-none focus:border-cyan-500 text-black"
              />

              {/* 削除ボタン */}
              <button
                onClick={() => deleteTask(task.id)}
                className="bg-cyan-500 text-white font-bold py-2 px-4 rounded-lg hover:bg-cyan-600 transition"
              >
                削除
              </button>
            </div>
          ))}
        </div>

        {/* タスクが空の場合のメッセージ */}
        {tasks.length === 0 && (
          <div className="text-center mt-12">
            <p className="text-white text-lg">タスクを追加してください</p>
          </div>
        )}
      </div>
    </div>
  );
}
