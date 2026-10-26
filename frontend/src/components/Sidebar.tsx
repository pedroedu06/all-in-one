import { Gamepad2, LibraryBig, Settings, type LucideIcon } from "lucide-react";

export type Page = "library" | "settings";

const items: { page: Page; label: string; icon: LucideIcon }[] = [
  { page: "library", label: "Biblioteca", icon: LibraryBig },
  { page: "settings", label: "Configurações", icon: Settings },
];

type Props = {
  current: Page;
  onNavigate: (page: Page) => void;
};

export function Sidebar({ current, onNavigate }: Props) {
  return (
    <nav className="sidebar" aria-label="Navegação principal">
      <div className="brand">
        <Gamepad2 size={22} aria-hidden />
        AllInOne
      </div>
      {items.map(({ page, label, icon: Icon }) => (
        <button
          key={page}
          type="button"
          className="nav-item"
          aria-current={current === page ? "page" : undefined}
          onClick={() => onNavigate(page)}
        >
          <Icon size={18} aria-hidden />
          {label}
        </button>
      ))}
    </nav>
  );
}
