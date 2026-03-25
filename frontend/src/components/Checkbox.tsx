export function Checkbox({ status, className = "", onClick }: { status: number; className?: string; onClick?: () => void }) {
  return (
    <div onClick={onClick} className={`w-16 border border-black flex items-center justify-center shrink-0 transition-colors ${onClick ? "cursor-pointer hover:bg-green-100" : ""} ${className}`}>
      {status === 1 && <CheckMark/>}
    </div>
  );
}

function CheckMark() {
  return (
    <div className="w-8 h-8 bg-green-400 rounded-sm flex items-center justify-center">
      <svg xmlns="http://www.w3.org/2000/svg" className="h-6 w-6 text-white" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={3}>
        <path strokeLinecap="round" strokeLinejoin="round" d="M5 13l4 4L19 7" />
      </svg>
    </div>
  );
}
