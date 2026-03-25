import { ReactNode } from 'react';

export function Title({ children }: { children: ReactNode }) {
  return (
    <div className="bg-white border border-gray-600 py-4 px-16 shadow-sm">
      <h1 className="text-2xl text-gray-800">{children}</h1>
    </div>
  );
}
