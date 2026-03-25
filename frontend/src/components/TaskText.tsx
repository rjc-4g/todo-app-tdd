import { ReactNode } from 'react';

export function TaskText({ children }: { children: ReactNode }) {
  return (
    <div className="flex-1 bg-white border border-black flex items-center px-4 text-xl text-gray-800">
      {children}
    </div>
  );
}
