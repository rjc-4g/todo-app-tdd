import { useEffect, useRef, useState } from "react";
import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import {
  createTask,
  deleteTask,
  fetchTasks,
  updateTaskStatus,
} from "../api/tasks";
import type { Task } from "../types/task";

function CheckIcon() {
  return (
    <svg
      className="w-5 h-5 text-white"
      fill="none"
      viewBox="0 0 24 24"
      stroke="currentColor"
      strokeWidth={3}
    >
      <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
    </svg>
  );
}

interface TaskRowProps {
  task: Task;
  onToggle: (id: number, current: number) => void;
  onDelete: (id: number) => void;
}

function TaskRow({ task, onToggle, onDelete }: TaskRowProps) {
  const checked = task.status === 1;
  return (
    <div className="flex items-center gap-3">
      <button
        onClick={() => onToggle(task.id, task.status)}
        className={`w-9 h-9 rounded-lg border-2 flex items-center justify-center flex-shrink-0 transition-colors ${
          checked
            ? "bg-green-500 border-green-500"
            : "bg-white border-gray-300 hover:border-gray-400"
        }`}
        aria-label={checked ? "完了を取り消す" : "完了にする"}
      >
        {checked && <CheckIcon />}
      </button>

      <input
        type="text"
        value={task.name}
        readOnly
        className="flex-1 bg-white rounded-xl px-4 py-2 text-gray-800 outline-none select-none"
      />

      <button
        onClick={() => onDelete(task.id)}
        className="bg-white text-blue-500 font-medium px-4 py-2 rounded-xl shadow-sm hover:bg-blue-50 transition-colors"
      >
        削除
      </button>
    </div>
  );
}

interface NewTaskRowProps {
  onSave: (name: string) => void;
  onCancel: () => void;
}

function NewTaskRow({ onSave, onCancel }: NewTaskRowProps) {
  const [name, setName] = useState("");
  const inputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    inputRef.current?.focus();
  }, []);

  const handleKeyDown = (e: React.KeyboardEvent<HTMLInputElement>) => {
    if (e.key === "Enter") {
      const trimmed = name.trim();
      if (trimmed) onSave(trimmed);
    } else if (e.key === "Escape") {
      onCancel();
    }
  };

  return (
    <div className="flex items-center gap-3">
      <div className="w-9 h-9 rounded-lg border-2 border-gray-300 bg-white flex-shrink-0" />
      <input
        ref={inputRef}
        type="text"
        value={name}
        onChange={(e) => setName(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="タスクを入力してください"
        className="flex-1 bg-white rounded-xl px-4 py-2 text-gray-800 outline-none placeholder:text-gray-400"
      />
      <button
        onClick={onCancel}
        className="bg-white text-blue-500 font-medium px-4 py-2 rounded-xl shadow-sm hover:bg-blue-50 transition-colors"
      >
        削除
      </button>
    </div>
  );
}

export default function TodoPage() {
  const [isAdding, setIsAdding] = useState(false);
  const queryClient = useQueryClient();

  const { data: tasks = [] } = useQuery({
    queryKey: ["tasks"],
    queryFn: fetchTasks,
  });

  const createMutation = useMutation({
    mutationFn: (name: string) => createTask(name),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
      setIsAdding(false);
    },
  });

  const updateMutation = useMutation({
    mutationFn: ({ id, status }: { id: number; status: number }) =>
      updateTaskStatus(id, status),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: number) => deleteTask(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["tasks"] });
    },
  });

  const handleToggle = (id: number, current: number) => {
    updateMutation.mutate({ id, status: current === 0 ? 1 : 0 });
  };

  return (
    <div className="min-h-screen bg-teal-500 p-6">
      {/* Header */}
      <div className="flex items-center gap-4 mb-6">
        <button
          onClick={() => setIsAdding(true)}
          disabled={isAdding}
          className="bg-white text-teal-600 font-semibold px-5 py-2 rounded-xl shadow-sm hover:bg-teal-50 transition-colors disabled:opacity-50 disabled:cursor-not-allowed"
        >
          新規
        </button>
        <div className="mx-auto border-2 border-gray-800 bg-white px-8 py-2 text-gray-800 font-bold text-lg tracking-widest">
          TODO APP TDD
        </div>
      </div>

      {/* Task list */}
      <div className="flex flex-col gap-3">
        {tasks.map((task) => (
          <TaskRow
            key={task.id}
            task={task}
            onToggle={handleToggle}
            onDelete={(id) => deleteMutation.mutate(id)}
          />
        ))}

        {/* 新規ボタンで追加される入力行 */}
        {isAdding && (
          <NewTaskRow
            onSave={(name) => createMutation.mutate(name)}
            onCancel={() => setIsAdding(false)}
          />
        )}
      </div>
    </div>
  );
}
