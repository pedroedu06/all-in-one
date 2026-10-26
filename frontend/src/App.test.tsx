import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { type Game, Source } from "../bindings/allinone/internal/models";
import App from "./App";

const listGames = vi.hoisted(() => vi.fn<(query: string, includeHidden: boolean) => Promise<Game[] | null>>());

vi.mock("../bindings/allinone/internal/library", () => ({
  Service: { ListGames: listGames },
}));

function game(fields: Pick<Game, "id" | "name">): Game {
  return {
    exePath: "",
    workDir: "",
    args: "",
    source: Source.SourceManual,
    coverPath: "",
    iconPath: "",
    favorite: false,
    hidden: false,
    tags: null,
    createdAt: "",
    updatedAt: "",
    ...fields,
  };
}

describe("App", () => {
  beforeEach(() => {
    listGames.mockReset();
  });

  // A fresh install has no games; the user must see that, not a blank screen.
  it("shows the empty state when the library has no games", async () => {
    listGames.mockResolvedValue(null);
    render(<App />);
    expect(await screen.findByText("Nenhum jogo na biblioteca")).toBeInTheDocument();
  });

  // The library view is driven by the Go backend, not local state.
  it("renders the games returned by the backend", async () => {
    listGames.mockResolvedValue([game({ id: 1, name: "Hades" }), game({ id: 2, name: "Celeste" })]);
    render(<App />);
    expect(await screen.findByText("Hades")).toBeInTheDocument();
    expect(screen.getByText("Celeste")).toBeInTheDocument();
  });

  // Search is done in SQLite; the UI must forward the typed term and keep
  // hidden games out of the default view.
  it("forwards the search term to the backend without hidden games", async () => {
    listGames.mockResolvedValue([]);
    render(<App />);
    fireEvent.change(screen.getByLabelText("Buscar jogos"), { target: { value: "knight" } });
    await waitFor(() => expect(listGames).toHaveBeenLastCalledWith("knight", false));
  });

  // Fail loud: a backend error must be visible to the user.
  it("shows backend errors instead of an empty library", async () => {
    listGames.mockRejectedValue(new Error("database is locked"));
    render(<App />);
    expect(await screen.findByRole("alert")).toHaveTextContent("database is locked");
    expect(screen.queryByText("Nenhum jogo na biblioteca")).not.toBeInTheDocument();
  });
});
