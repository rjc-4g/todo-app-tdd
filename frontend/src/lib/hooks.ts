import { useCallback, useEffect, useState } from "react";
import { Task } from "./types";
import { deleteTask, getTasks, patchTask, postTask } from "./api";

export const useTask = () => {

  const [tasks, setTasks] = useState<Task[]>([]);
  const [draftTask, setDraftTask] = useState<string | null>(null);

  const handleAddTask = useCallback(() => {
    setDraftTask("");
  }, []);

  const handleConfirm = useCallback(async () => {
    if (draftTask === null || draftTask === "") {
      return;
    }
    const response = await postTask(draftTask);
    if (response?.ok) {
      const newTask = (await response.json()) as Task;
      setTasks((prev) => [...prev, newTask]);
      setDraftTask(null);
    }
  }, [draftTask]);

  const handleUpdateStatus = useCallback(async (id: number, status: number) => {
    const response = await patchTask(id, status);
    if (response?.ok) {
      setTasks((prev) => prev.map(task => (task.id === id ? { ...task, status } : task)));
    }
  }, []);

  const handleDelete = useCallback(async (id: number) => {
    const response = await deleteTask(id);
    if (response?.ok) {
      setTasks((prev) => prev.filter(task => task.id !== id));
    }
  }, []);

  useEffect(() => {
    const fetchTasks = async () => {
      const response = await getTasks();
      if (response?.ok) {
        const data = (await response.json()) as Task[];
        const activeTasks = data.filter(task => task.deleted === false);
        setTasks(activeTasks);
      }
    };
    fetchTasks();
  }, []);

  return {
    tasks,
    draftTask,
    setDraftTask,
    handleAddTask,
    handleConfirm,
    handleUpdateStatus,
    handleDelete
  };
};
