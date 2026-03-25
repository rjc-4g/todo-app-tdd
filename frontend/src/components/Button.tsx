import { ReactNode } from 'react';

export function Button({ children, className = "" }: { children: ReactNode; className?: string }) {
  return (
    <button className={`bg-white border border-black text-blue-400 text-lg font-bold px-8 rounded-2xl shadow-md hover:bg-gray-50 transition-colors ${className}`}>
      {children}
    </button>
  );
}
