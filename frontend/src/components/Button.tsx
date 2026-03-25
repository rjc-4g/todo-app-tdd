import { ReactNode } from 'react';

export function Button({ children, className = "", onClick }: { children: ReactNode; className?: string; onClick?: () => void }) {
  return (
    <button onClick={onClick} className={`bg-white border border-black text-blue-400 text-lg font-bold px-8 rounded-2xl shadow-md hover:bg-gray-50 transition-colors ${className}`}>
      {children}
    </button>
  );
}
