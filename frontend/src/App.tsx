import { useEffect, useState } from "react";
import { Button } from "./components/Button";
import { DraftTaskRow } from "./components/DraftTaskRow";
import { TaskRow } from "./components/TaskRow";
import { Title } from "./components/Title";
import { getTasks, postTask } from "./lib/api";
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

  useEffect(() => {
    getTasks().then(res => res.json())
              .then(data => setTasks(data));
  }, []);


  return (
    <div className="min-h-screen bg-teal-600 p-8 font-sans">

      <div className="relative flex justify-center items-center mb-12 max-w-4xl mx-auto">
        <Button className="absolute left-0 py-3" onClick={handleAddTask}>新規</Button>
        <Title>TODO APP TDD</Title>
      </div>

      <div className="space-y-6 max-w-4xl mx-auto">
        {tasks.filter(task => task.deleted === false)
              .map(task => <TaskRow key={task.id} task={task} />)}
        {draftTask !== null && 
          <DraftTaskRow draftTask={draftTask} setDraftTask={setDraftTask} handleConfirm={handleConfirm} />}
      </div>

    </div>
  );
}
