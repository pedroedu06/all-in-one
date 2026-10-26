import { fireEvent, render, screen, waitFor } from "@testing-library/react";
import { beforeEach, describe, expect, it, vi } from "vitest";
import type { ScanFolder } from "../../bindings/allinone/internal/models";
import { Settings } from "./Settings";

const svc = vi.hoisted(() => ({
  ListFolders: vi.fn<() => Promise<ScanFolder[] | null>>(),
  PickFolder: vi.fn<() => Promise<string>>(),
  AddFolder: vi.fn<(path: string) => Promise<ScanFolder>>(),
  RemoveFolder: vi.fn<(id: number) => Promise<void>>(),
  GetSteamGridDBKey: vi.fn<() => Promise<string>>(),
  SetSteamGridDBKey: vi.fn<(key: string) => Promise<void>>(),
  TestSteamGridDBKey: vi.fn<(key: string) => Promise<void>>(),
}));

vi.mock("../../bindings/allinone/internal/settings", () => ({ Service: svc }));
vi.mock("@wailsio/runtime", () => ({ Browser: { OpenURL: vi.fn(() => Promise.resolve()) } }));

const folder = (id: number, path: string): ScanFolder => ({ id, path, createdAt: "" });

function typeKey(value: string) {
  fireEvent.change(screen.getByLabelText("Chave da API"), { target: { value } });
}

describe("Settings", () => {
  beforeEach(() => {
    for (const fn of Object.values(svc)) fn.mockReset();
    svc.ListFolders.mockResolvedValue([]);
    svc.GetSteamGridDBKey.mockResolvedValue("");
  });

  describe("game folders", () => {
    // Without folders the scanner has nothing to do; the user must be told and
    // offered the action right there.
    it("shows the empty state with an add action", async () => {
      render(<Settings />);
      expect(await screen.findByText("Nenhuma pasta configurada")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: /adicionar pasta/i })).toBeInTheDocument();
    });

    // Folders come from the native picker so the stored path is always real.
    it("adds the folder chosen in the native picker", async () => {
      svc.PickFolder.mockResolvedValue("D:\\Jogos");
      svc.AddFolder.mockResolvedValue(folder(1, "D:\\Jogos"));
      render(<Settings />);

      svc.ListFolders.mockResolvedValue([folder(1, "D:\\Jogos")]);
      fireEvent.click(await screen.findByRole("button", { name: /adicionar pasta/i }));

      expect(await screen.findByText("D:\\Jogos")).toBeInTheDocument();
      expect(svc.AddFolder).toHaveBeenCalledWith("D:\\Jogos");
    });

    // Closing the picker is not an error and must not add an empty path.
    it("does nothing when the picker is cancelled", async () => {
      svc.PickFolder.mockResolvedValue("");
      render(<Settings />);
      fireEvent.click(await screen.findByRole("button", { name: /adicionar pasta/i }));

      await waitFor(() => expect(svc.PickFolder).toHaveBeenCalled());
      expect(svc.AddFolder).not.toHaveBeenCalled();
    });

    it("removes a folder", async () => {
      svc.ListFolders.mockResolvedValue([folder(7, "C:\\Games")]);
      svc.RemoveFolder.mockResolvedValue();
      render(<Settings />);

      fireEvent.click(await screen.findByRole("button", { name: "Remover C:\\Games" }));

      await waitFor(() => expect(screen.queryByText("C:\\Games")).not.toBeInTheDocument());
      expect(svc.RemoveFolder).toHaveBeenCalledWith(7);
    });

    // Fail loud: e.g. adding the same folder twice must tell the user why.
    it("shows backend errors", async () => {
      svc.PickFolder.mockResolvedValue("D:\\Jogos");
      svc.AddFolder.mockRejectedValue(new Error("Essa pasta já foi adicionada"));
      render(<Settings />);
      fireEvent.click(await screen.findByRole("button", { name: /adicionar pasta/i }));

      expect(await screen.findByRole("alert")).toHaveTextContent("Essa pasta já foi adicionada");
    });
  });

  describe("SteamGridDB key", () => {
    // The saved key must be shown so the user knows one is configured.
    it("loads the saved key", async () => {
      svc.GetSteamGridDBKey.mockResolvedValue("saved-key");
      render(<Settings />);
      await waitFor(() => expect(screen.getByLabelText("Chave da API")).toHaveValue("saved-key"));
    });

    // The point of the test button: confirm the key works before relying on it.
    it("reports a valid key", async () => {
      svc.TestSteamGridDBKey.mockResolvedValue();
      render(<Settings />);
      typeKey("good");
      fireEvent.click(screen.getByRole("button", { name: "Testar chave" }));

      expect(await screen.findByText("Chave válida")).toBeInTheDocument();
      expect(svc.TestSteamGridDBKey).toHaveBeenCalledWith("good");
    });

    it("shows why the key was rejected", async () => {
      svc.TestSteamGridDBKey.mockRejectedValue(new Error("Chave inválida"));
      render(<Settings />);
      typeKey("bad");
      fireEvent.click(screen.getByRole("button", { name: "Testar chave" }));

      expect(await screen.findByRole("alert")).toHaveTextContent("Chave inválida");
    });

    it("saves the typed key", async () => {
      svc.SetSteamGridDBKey.mockResolvedValue();
      render(<Settings />);
      typeKey("new-key");
      fireEvent.click(screen.getByRole("button", { name: "Salvar" }));

      await waitFor(() => expect(svc.SetSteamGridDBKey).toHaveBeenCalledWith("new-key"));
      expect(await screen.findByText("Chave salva")).toBeInTheDocument();
    });
  });
});
