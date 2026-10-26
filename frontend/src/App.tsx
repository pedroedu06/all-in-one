import { useState } from "react";
import { Sidebar, type Page } from "./components/Sidebar";
import { Library } from "./pages/Library";
import { Settings } from "./pages/Settings";

function App() {
  const [page, setPage] = useState<Page>("library");

  return (
    <div className="shell">
      <Sidebar current={page} onNavigate={setPage} />
      <main className="content">
        {page === "library" ? <Library onOpenSettings={() => setPage("settings")} /> : <Settings />}
      </main>
    </div>
  );
}

export default App;
