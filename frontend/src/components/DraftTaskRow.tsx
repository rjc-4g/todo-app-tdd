import { Button } from './Button';
import { Checkbox } from './Checkbox';
import { TaskText } from './TaskText';

export function DraftTaskRow({
  draftTask,
  setDraftTask,
  handleConfirm
}: {
  draftTask: string;
  setDraftTask: (value: string) => void;
  handleConfirm: () => void;
}) {
  return (
    <div className="flex items-stretch gap-4 h-16">
      <Checkbox status={0} className="bg-gray-200 disabled" />
      <TaskText>
        <input
          className="w-full text-xl text-gray-800 outline-none placeholder:text-gray-400"
          placeholder="タスクを入力してください"
          value={draftTask}
          onChange={(e) => setDraftTask(e.target.value)}
          autoFocus
        />
      </TaskText>
      <Button onClick={handleConfirm}>確定</Button>
    </div>
  );
}
