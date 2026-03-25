import { useEffect, useState } from "react";
import { Button } from "./components/Button";
import { DraftTaskRow } from "./components/DraftTaskRow";
import { TaskRow } from "./components/TaskRow";
import { Title } from "./components/Title";
import { deleteTask, getTasks, patchTask, postTask } from "./lib/api";
import { Task } from "./lib/types";

export default function App() {

  const [tasks, setTasks] = useState<Task[]>([]);
  const [draftTask, setDraftTask] = useState<string | null>(null);

  const handleAddTask = () => {
    setDraftTask("");
  };

  const handleConfirm = async () => {
    if (draftTask === null || draftTask === "") return;

    const response = await postTask(draftTask);

    if (response.ok) {
      const newTask = await response.json();
      setTasks([...tasks, newTask]);
      setDraftTask(null);
    }
  };

  const handleUpdateStatus = async (id: number, status: number) => {

    const response = await patchTask(id, status);

    if (response.ok) {
      setTasks(tasks.map(task => (task.id === id ? { ...task, status } : task)));
    }
  };

  const handleDelete = async (id: number) => {

    const response = await deleteTask(id);

    if (response.ok) {
      setTasks(tasks.filter(task => task.id !== id));
    }
  };

  useEffect(() => {
    getTasks().then(res => res.json() as Promise<Task[]>)
              .then(tasks => tasks.filter(task => task.deleted === false))
              .then(tasks => setTasks(tasks));
  }, []);


  return (
    <div className="min-h-screen bg-teal-600 p-8 font-sans">

      <div className="relative flex justify-center items-center mb-12 max-w-4xl mx-auto">
        <Button className="absolute left-0 py-3" onClick={handleAddTask}>新規</Button>
        <Title>TODO APP TDD</Title>
      </div>

      <div className="space-y-6 max-w-4xl mx-auto">
        {tasks.map(task => <TaskRow key={task.id} task={task} handleDelete={handleDelete} handleUpdateStatus={handleUpdateStatus} />)}
        {draftTask !== null && 
          <DraftTaskRow draftTask={draftTask} setDraftTask={setDraftTask} handleConfirm={handleConfirm} />}
      </div>

    </div>
  );
}
