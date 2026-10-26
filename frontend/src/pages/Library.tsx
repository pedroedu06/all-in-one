import { useEffect, useState } from "react";
import { FolderPlus, Gamepad2 } from "lucide-react";
import { Service as LibraryService } from "../../bindings/allinone/internal/library";
import type { Game } from "../../bindings/allinone/internal/models";
import { errorMessage } from "../lib/errors";
import { staggerStyle } from "../lib/stagger";

type Props = {
  onOpenSettings: () => void;
};

export function Library({ onOpenSettings }: Props) {
  const [query, setQuery] = useState("");
  const [games, setGames] = useState<Game[] | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;
    LibraryService.ListGames(query, false)
      .then((result) => {
        if (cancelled) return;
        // Go nil slices arrive as null.
        setGames(result ?? []);
        setError(null);
      })
      .catch((err: unknown) => {
        if (!cancelled) setError(errorMessage(err));
      });
    return () => {
      cancelled = true;
    };
  }, [query]);

  return (
    <section className="enter">
      <header className="page-header">
        <h1>Biblioteca</h1>
        <input
          type="search"
          className="input"
          style={{ maxWidth: 360 }}
          aria-label="Buscar jogos"
          placeholder="Buscar jogos..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </header>

      {error ? (
        <p role="alert" className="status error">
          Erro ao carregar a biblioteca: {error}
        </p>
      ) : games === null ? (
        <div className="game-grid" aria-hidden>
          {[0, 1, 2, 3].map((i) => (
            <div key={i} className="skeleton" style={{ height: "auto", aspectRatio: "2 / 3" }} />
          ))}
        </div>
      ) : games.length === 0 ? (
        <div className="card empty">
          <div className="empty-icon">
            <Gamepad2 size={28} aria-hidden />
          </div>
          <h2>Nenhum jogo na biblioteca</h2>
          <p className="muted">Adicione as pastas onde seus jogos estão instalados para começar.</p>
          <button type="button" className="btn btn-primary" onClick={onOpenSettings}>
            <FolderPlus size={18} aria-hidden />
            Configurar pastas
          </button>
        </div>
      ) : (
        <ul className="game-grid stagger">
          {games.map((g, i) => (
            <li key={g.id} className="game-card" style={staggerStyle(i)}>
              {g.name}
            </li>
          ))}
        </ul>
      )}
    </section>
  );
}
