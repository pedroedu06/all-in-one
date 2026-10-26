import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { type Game, Source } from "../../bindings/allinone/internal/models";
import { Library } from "./Library";

const listGames = vi.hoisted(() => vi.fn<(query: string, includeHidden: boolean) => Promise<Game[] | null>>());

vi.mock("../../bindings/allinone/internal/library", () => ({
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

describe("Library", () => {
  beforeEach(() => {
    listGames.mockReset();
  });

  // A fresh install has no games; the user must see that and how to fix it.
  it("shows the empty state with a shortcut to configure folders", async () => {
    listGames.mockResolvedValue(null);
    const onOpenSettings = vi.fn();
    render(<Library onOpenSettings={onOpenSettings} />);

    expect(await screen.findByText("Nenhum jogo na biblioteca")).toBeInTheDocument();
    fireEvent.click(screen.getByRole("button", { name: /configurar pastas/i }));
    expect(onOpenSettings).toHaveBeenCalled();
  });

  // The library view is driven by the Go backend, not local state.
  it("renders the games returned by the backend", async () => {
    listGames.mockResolvedValue([game({ id: 1, name: "Hades" }), game({ id: 2, name: "Celeste" })]);
    render(<Library onOpenSettings={vi.fn()} />);
    expect(await screen.findByText("Hades")).toBeInTheDocument();
    expect(screen.getByText("Celeste")).toBeInTheDocument();
  });

  // Search is done in SQLite; the UI must forward the typed term and keep
  // hidden games out of the default view.
  it("forwards the search term to the backend without hidden games", async () => {
    listGames.mockResolvedValue([]);
    render(<Library onOpenSettings={vi.fn()} />);
    fireEvent.change(screen.getByLabelText("Buscar jogos"), { target: { value: "knight" } });
    await waitFor(() => expect(listGames).toHaveBeenLastCalledWith("knight", false));
  });

  // Fail loud: a backend error must be visible to the user.
  it("shows backend errors instead of an empty library", async () => {
    listGames.mockRejectedValue(new Error("database is locked"));
    render(<Library onOpenSettings={vi.fn()} />);
    expect(await screen.findByRole("alert")).toHaveTextContent("database is locked");
    expect(screen.queryByText("Nenhum jogo na biblioteca")).not.toBeInTheDocument();
  });
});
