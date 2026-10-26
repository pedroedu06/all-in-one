import { useEffect } from "react";
import { CheckCircle2, XCircle } from "lucide-react";

export type ToastMessage = { kind: "ok" | "error"; text: string };

type Props = {
  toast: ToastMessage | null;
  onClose: () => void;
};

const DURATION_MS = 4000;

export function Toast({ toast, onClose }: Props) {
  useEffect(() => {
    if (!toast) return;
    const timer = setTimeout(onClose, DURATION_MS);
    return () => clearTimeout(timer);
  }, [toast, onClose]);

  if (!toast) return null;
  const Icon = toast.kind === "ok" ? CheckCircle2 : XCircle;
  return (
    <div className={`toast ${toast.kind}`} role={toast.kind === "error" ? "alert" : "status"}>
      <Icon size={18} aria-hidden />
      {toast.text}
    </div>
  );
}
