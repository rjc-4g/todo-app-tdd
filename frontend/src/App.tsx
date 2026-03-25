import { Button } from "./components/Button";
import { TaskRow } from "./components/TaskRow";
import { Title } from "./components/Title";
import { Task } from "./types";

export default function App() {

  const tasks: Task[] = [
    { id: 1, name: "タスク1", status: 0, created: new Date(), updated: new Date(), deleted: false },
    { id: 2, name: "タスク2", status: 1, created: new Date(), updated: new Date(), deleted: false },
    { id: 3, name: "タスクを入力してください", status: 0, created: new Date(), updated: new Date(), deleted: false },
  ];

  return (
    <div className="min-h-screen bg-teal-600 p-8 font-sans">

      <div className="relative flex justify-center items-center mb-12 max-w-4xl mx-auto">
        <Button className="absolute left-0 py-3">新規</Button>
        <Title>TODO APP TDD</Title>
      </div>

      <div className="space-y-6 max-w-4xl mx-auto">
        {tasks.map(task => <TaskRow task={task} />)}
      </div>

    </div>
  );
}
