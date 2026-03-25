import { Task } from '../types';
import { Button } from './Button';
import { Checkbox } from './Checkbox';
import { TaskText } from './TaskText';

export function TaskRow({ task }: { task: Task }) {
  return (
    <div className="flex items-stretch gap-4 h-16">
      <Checkbox status={task.status} />
      <TaskText>{task.name}</TaskText>
      <Button>削除</Button>
    </div>
  );
}
