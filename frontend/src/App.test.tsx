import { fireEvent, render, screen } from "@testing-library/react";
import { describe, expect, it, vi } from "vitest";
import App from "./App";

vi.mock("./pages/Library", () => ({ Library: () => <p>library page</p> }));
vi.mock("./pages/Settings", () => ({ Settings: () => <p>settings page</p> }));

describe("App", () => {
  // The library is what the user opens the launcher for; settings is one click away.
  it("starts on the library and navigates to settings", () => {
    render(<App />);
    expect(screen.getByText("library page")).toBeInTheDocument();

    fireEvent.click(screen.getByRole("button", { name: /configurações/i }));

    expect(screen.getByText("settings page")).toBeInTheDocument();
    expect(screen.getByRole("button", { name: /configurações/i })).toHaveAttribute("aria-current", "page");
  });
});
