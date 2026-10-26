import { useEffect, useState } from "react";
import { Service as LibraryService } from "../bindings/allinone/internal/library";
import type { Game } from "../bindings/allinone/internal/models";

function App() {
  const [query, setQuery] = useState("");
  const [games, setGames] = useState<Game[]>([]);
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
        if (cancelled) return;
        setError(err instanceof Error ? err.message : String(err));
      });
    return () => {
      cancelled = true;
    };
  }, [query]);

  return (
    <main className="app">
      <header className="toolbar">
        <h1>AllInOne</h1>
        <input
          type="search"
          aria-label="Buscar jogos"
          placeholder="Buscar jogos..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
        />
      </header>

      {error && (
        <p role="alert" className="error">
          Erro ao carregar a biblioteca: {error}
        </p>
      )}

      {!error && games.length === 0 ? (
        <p className="empty">Nenhum jogo na biblioteca</p>
      ) : (
        <ul className="grid">
          {games.map((g) => (
            <li key={g.id} className="card">
              {g.favorite && <span aria-label="Favorito">★ </span>}
              {g.name}
            </li>
          ))}
        </ul>
      )}
    </main>
  );
}

export default App;
