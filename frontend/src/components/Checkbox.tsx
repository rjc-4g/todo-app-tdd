import { InputHTMLAttributes } from "react";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface CheckboxProps extends Omit<InputHTMLAttributes<HTMLInputElement>, "type"> {
  // 追加のPropsがあればここに記述
}

export default function Checkbox({ className, ...props }: CheckboxProps) {
  return (
    <input
      type="checkbox"
      className={cn(
        "relative w-6 h-6 shrink-0 cursor-pointer appearance-none rounded border-2 border-gray-300 bg-white transition-all",
        "checked:border-lime-500 checked:bg-lime-500",
        "focus:outline-none focus:ring-0",
        /* レ点（チェックマーク）の精密な中央配置 */
        "after:content-[''] after:absolute after:hidden",
        "after:top-[3px] after:left-[8px] after:w-[6px] after:h-[12px]",
        "after:border-white after:border-r-2 after:border-b-2 after:rotate-45",
        "checked:after:block",
        className
      )}
      {...props}
    />
  );
}
