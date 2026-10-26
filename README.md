# AllInOne

Launcher de biblioteca de jogos para Windows (Linux depois). Backend em Go, UI em React/TypeScript via Wails v3, empacotado como um `.exe` único sem CGO.

## Stack

- **Wails v3** — `main.go` monta o app e registra os serviços.
- **SQLite** via **GORM** + `github.com/glebarez/sqlite` (puro Go, `modernc.org/sqlite`). Banco em `%APPDATA%\AllInOne\allinone.db`.
- **React 18 + TS** em `frontend/`, testes com **Vitest** + Testing Library.

## Estrutura

```
main.go                    entrada: abre o banco e registra serviços
internal/db                Open/Close do SQLite + AutoMigrate
internal/models            modelos GORM (Game, ScanFolder, Setting)
internal/library           Store (consultas) + Service (exposto ao frontend)
frontend/src               UI React
frontend/bindings          bindings TS gerados pelo Wails (não editar)
```

## Comandos

| Comando | O que faz |
|---|---|
| `wails3 task dev` | modo desenvolvimento com hot reload |
| `wails3 task test` | todos os testes (Go com `CGO_ENABLED=0` + Vitest) |
| `wails3 task test:go` | só testes Go |
| `wails3 task test:frontend` | só testes do frontend |
| `wails3 task lint` | `go vet` + checagem de tipos TS |
| `wails3 task build` | gera `bin/allinone.exe` |

Após mudar métodos de serviços Go, as bindings são regeneradas automaticamente pelo `build`/`dev`.
