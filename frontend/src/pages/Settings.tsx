import { useCallback, useEffect, useState } from "react";
import { CheckCircle2, Eye, EyeOff, Folder, FolderPlus, FolderSearch, Trash2, XCircle } from "lucide-react";
import { Browser } from "@wailsio/runtime";
import { Service as SettingsService } from "../../bindings/allinone/internal/settings";
import type { ScanFolder } from "../../bindings/allinone/internal/models";
import { Toast, type ToastMessage } from "../components/Toast";
import { errorMessage } from "../lib/errors";
import { staggerStyle } from "../lib/stagger";

const API_KEY_URL = "https://www.steamgriddb.com/profile/preferences/api";

type Notify = (toast: ToastMessage) => void;

export function Settings() {
  const [toast, setToast] = useState<ToastMessage | null>(null);
  const closeToast = useCallback(() => setToast(null), []);

  return (
    <section className="enter">
      <header className="page-header">
        <h1>Configurações</h1>
      </header>
      <FoldersSection notify={setToast} />
      <SteamGridDBSection notify={setToast} />
      <Toast toast={toast} onClose={closeToast} />
    </section>
  );
}

function FoldersSection({ notify }: { notify: Notify }) {
  const [folders, setFolders] = useState<ScanFolder[] | null>(null);

  const load = useCallback(() => {
    SettingsService.ListFolders()
      .then((result) => setFolders(result ?? []))
      .catch((err: unknown) => notify({ kind: "error", text: errorMessage(err) }));
  }, [notify]);

  useEffect(load, [load]);

  const add = async () => {
    try {
      const path = await SettingsService.PickFolder();
      if (!path) return; // picker cancelled
      await SettingsService.AddFolder(path);
      notify({ kind: "ok", text: "Pasta adicionada" });
      load();
    } catch (err: unknown) {
      notify({ kind: "error", text: errorMessage(err) });
    }
  };

  const remove = async (folder: ScanFolder) => {
    try {
      await SettingsService.RemoveFolder(folder.id);
      setFolders((prev) => (prev ?? []).filter((f) => f.id !== folder.id));
    } catch (err: unknown) {
      notify({ kind: "error", text: errorMessage(err) });
    }
  };

  const addButton = (
    <button type="button" className="btn btn-primary" onClick={add}>
      <FolderPlus size={18} aria-hidden />
      Adicionar pasta
    </button>
  );

  return (
    <div className="card section">
      <div className="section-head">
        <div>
          <h2>Pastas de jogos</h2>
          <p className="muted">O AllInOne procura executáveis de jogos dentro destas pastas.</p>
        </div>
        {folders && folders.length > 0 && addButton}
      </div>

      {folders === null ? (
        <div className="folder-list" aria-hidden>
          <div className="skeleton" />
          <div className="skeleton" />
        </div>
      ) : folders.length === 0 ? (
        <div className="empty">
          <div className="empty-icon">
            <FolderSearch size={28} aria-hidden />
          </div>
          <h3>Nenhuma pasta configurada</h3>
          <p className="muted">Escolha onde seus jogos estão instalados, por exemplo D:\Jogos.</p>
          {addButton}
        </div>
      ) : (
        <ul className="folder-list stagger">
          {folders.map((f, i) => (
            <li key={f.id} className="folder-item" style={staggerStyle(i)}>
              <Folder size={18} aria-hidden />
              <span className="mono" title={f.path}>
                {f.path}
              </span>
              <button
                type="button"
                className="icon-btn danger"
                aria-label={`Remover ${f.path}`}
                onClick={() => remove(f)}
              >
                <Trash2 size={18} aria-hidden />
              </button>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
}

type TestResult = { kind: "ok" } | { kind: "error"; message: string } | null;

function SteamGridDBSection({ notify }: { notify: Notify }) {
  const [key, setKey] = useState("");
  const [visible, setVisible] = useState(false);
  const [testing, setTesting] = useState(false);
  const [saving, setSaving] = useState(false);
  const [result, setResult] = useState<TestResult>(null);

  useEffect(() => {
    SettingsService.GetSteamGridDBKey()
      .then(setKey)
      .catch((err: unknown) => notify({ kind: "error", text: errorMessage(err) }));
  }, [notify]);

  const test = async () => {
    setTesting(true);
    setResult(null);
    try {
      await SettingsService.TestSteamGridDBKey(key);
      setResult({ kind: "ok" });
    } catch (err: unknown) {
      setResult({ kind: "error", message: errorMessage(err) });
    } finally {
      setTesting(false);
    }
  };

  const save = async () => {
    setSaving(true);
    try {
      await SettingsService.SetSteamGridDBKey(key);
      notify({ kind: "ok", text: key.trim() ? "Chave salva" : "Chave removida" });
    } catch (err: unknown) {
      notify({ kind: "error", text: errorMessage(err) });
    } finally {
      setSaving(false);
    }
  };

  const openApiPage = () => {
    Browser.OpenURL(API_KEY_URL).catch((err: unknown) => notify({ kind: "error", text: errorMessage(err) }));
  };

  return (
    <div className="card section">
      <div className="section-head">
        <div>
          <h2>SteamGridDB</h2>
          <p className="muted">
            Opcional. Usada para baixar as capas dos jogos. Sem ela, o ícone do executável é usado. Gere sua chave
            gratuita em{" "}
            <button type="button" className="link" onClick={openApiPage}>
              steamgriddb.com
            </button>
            .
          </p>
        </div>
      </div>

      <div className="field" style={{ maxWidth: 560 }}>
        <label className="label" htmlFor="sgdb-key">
          Chave da API
        </label>
        <div className="input-wrap">
          <input
            id="sgdb-key"
            className="input mono"
            type={visible ? "text" : "password"}
            autoComplete="off"
            spellCheck={false}
            value={key}
            onChange={(e) => {
              setKey(e.target.value);
              setResult(null);
            }}
          />
          <button
            type="button"
            className="icon-btn"
            aria-label={visible ? "Ocultar chave" : "Mostrar chave"}
            onClick={() => setVisible((v) => !v)}
          >
            {visible ? <EyeOff size={18} aria-hidden /> : <Eye size={18} aria-hidden />}
          </button>
        </div>

        {result?.kind === "ok" && (
          <p className="status ok" role="status">
            <CheckCircle2 size={16} aria-hidden />
            Chave válida
          </p>
        )}
        {result?.kind === "error" && (
          <p className="status error" role="alert">
            <XCircle size={16} aria-hidden />
            {result.message}
          </p>
        )}

        <div className="row" style={{ marginTop: 8 }}>
          <button type="button" className="btn btn-ghost" onClick={test} disabled={testing || !key.trim()}>
            {testing ? "Testando..." : "Testar chave"}
          </button>
          <button type="button" className="btn btn-primary" onClick={save} disabled={saving}>
            {saving ? "Salvando..." : "Salvar"}
          </button>
        </div>
      </div>
    </div>
  );
}
