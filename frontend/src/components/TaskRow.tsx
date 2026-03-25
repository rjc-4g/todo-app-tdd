import { Task } from '../lib/types';
import { Button } from './Button';
import { Checkbox } from './Checkbox';
import { TaskText } from './TaskText';

export function TaskRow({
  task,
  handleDelete,
  handleUpdateStatus
}: {
  task: Task;
  handleDelete: (id: number) => void;
  handleUpdateStatus: (id: number, status: number) => void;
}) {
  return (
    <div className="flex items-stretch gap-4 h-16">
      <Checkbox status={task.status} className="bg-white" onClick={() => handleUpdateStatus(task.id, task.status === 0 ? 1 : 0)} />
      <TaskText>{task.name}</TaskText>
      <Button onClick={() => handleDelete(task.id)}>削除</Button>
    </div>
  );
}
