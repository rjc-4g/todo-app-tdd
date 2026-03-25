import { useState } from "react";

interface TaskFormProps {
  onSubmit: (name: string) => Promise<void>;
  onBlur?: () => void; // フォーカスアウト時のクローズ用
  placeholder?: string;
  autoFocus?: boolean;
}

export default function TaskForm({
  onSubmit,
  onBlur,
  placeholder = "タスク名を入力...",
  autoFocus = false,
}: TaskFormProps) {
  const [taskName, setTaskName] = useState("");
  const [isSubmitting, setIsSubmitting] = useState(false);

  // 保存処理の実行
  const commitChange = async () => {
    const trimmed = taskName.trim();
    if (!trimmed) {
      if (onBlur) onBlur(); // 空文字ならそのまま閉じる
      return;
    }

    if (isSubmitting) return;

    setIsSubmitting(true);
    try {
      await onSubmit(trimmed);
      // 成功時、親コンポーネント側で showForm(false) されることを想定
    } catch (error) {
      console.error("Failed to submit task:", error);
      setIsSubmitting(false);
    }
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    commitChange(); // Enterキー押下時
  };

  return (
    <form
      onSubmit={handleSubmit}
      className="flex gap-2 mb-10 bg-white p-2 rounded-lg shadow-md border border-gray-300"
    >
      <input
        type="text"
        name="task-name-input"
        value={taskName}
        onChange={(e) => setTaskName(e.target.value)}
        onBlur={commitChange} // フォーカスアウト時に保存を実行
        placeholder={placeholder}
        autoComplete="off"
        autoFocus={autoFocus}
        disabled={isSubmitting}
        className="flex-1 px-4 py-3 bg-white text-gray-900 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 border-none appearance-none disabled:opacity-50"
      />
    </form>
  );
}
