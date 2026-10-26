import type { CSSProperties } from "react";

// Index for the .stagger cascade (120ms between items).
export function staggerStyle(i: number): CSSProperties {
  return { "--i": i } as CSSProperties;
}
