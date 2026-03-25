import { ReactNode, ButtonHTMLAttributes } from "react";
import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

/**
 * クラス名のマージヘルパー（Tailwindの競合を解決）
 */
function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

interface ButtonProps extends ButtonHTMLAttributes<HTMLButtonElement> {
  children: ReactNode;
  variant?: "primary" | "secondary" | "danger" | "outline-danger";
  size?: "sm" | "md" | "lg";
}

export default function Button({
  children,
  className,
  variant = "primary",
  size = "md",
  ...props
}: ButtonProps) {
  // バリアントごとのスタイル定義
  const variants = {
    primary: "bg-blue-600 text-white hover:bg-blue-700 shadow-sm",
    secondary: "bg-gray-500 text-white hover:bg-gray-600 shadow-sm",
    danger: "bg-red-500 text-white hover:bg-red-600 shadow-sm",
    "outline-danger": "text-red-500 border border-red-200 hover:text-white hover:bg-red-500",
  };

  // サイズごとのスタイル定義
  const sizes = {
    sm: "px-3 py-1 text-sm",
    md: "px-4 py-2",
    lg: "px-8 py-3 font-bold",
  };

  return (
    <button
      className={cn(
        "rounded-md transition-all font-medium disabled:opacity-50 disabled:cursor-not-allowed",
        variants[variant],
        sizes[size],
        className
      )}
      {...props}
    >
      {children}
    </button>
  );
}
