import { Button } from "./components/Button";
import { DraftTaskRow } from "./components/DraftTaskRow";
import { TaskRow } from "./components/TaskRow";
import { Title } from "./components/Title";
import { useTask } from "./lib/hooks";

export default function App() {

  const {
    tasks,
    draftTask,
    setDraftTask,
    handleAddTask,
    handleConfirm,
    handleUpdateStatus,
    handleDelete
  } = useTask();

  return (
    <div className="min-h-screen bg-teal-600 p-8 font-sans">

      <div className="relative flex justify-center items-center mb-12 max-w-4xl mx-auto">
        <Button className="absolute left-0 py-3" onClick={handleAddTask}>新規</Button>
        <Title>TODO APP TDD</Title>
      </div>

      <div className="space-y-6 max-w-4xl mx-auto">
        {tasks.map(task => <TaskRow key={task.id} task={task} handleDelete={handleDelete} handleUpdateStatus={handleUpdateStatus} />)}
        {draftTask !== null && <DraftTaskRow draftTask={draftTask} setDraftTask={setDraftTask} handleConfirm={handleConfirm} />}
      </div>

    </div>
  );
}
