import { useEffect, useState } from "react";
import { Button } from "./components/Button";
import { TaskRow } from "./components/TaskRow";
import { Title } from "./components/Title";
import { API_BASE_URL } from "./lib/constants";
import { Task } from "./lib/types";

export default function App() {

  const [tasks, setTasks] = useState<Task[]>([]);

  useEffect(() => {
    fetch(`${API_BASE_URL}/api/v1/tasks`)
      .then((res) => res.json())
      .then((data) => setTasks(data));
  }, []);

  return (
    <div className="min-h-screen bg-teal-600 p-8 font-sans">

      <div className="relative flex justify-center items-center mb-12 max-w-4xl mx-auto">
        <Button className="absolute left-0 py-3">新規</Button>
        <Title>TODO APP TDD</Title>
      </div>

      <div className="space-y-6 max-w-4xl mx-auto">
        {tasks.filter(task => task.deleted === false)
              .map(task => <TaskRow key={task.id} task={task} />)}
      </div>

    </div>
  );
}
