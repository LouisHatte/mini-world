package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"mini-world-go/internal/commandrun"
	"mini-world-go/internal/commands"
	"mini-world-go/internal/world"
)

type commandRequest struct {
	Command string `json:"command"`
}

type commandResponse struct {
	Output string `json:"output"`
	Error  string `json:"error,omitempty"`
	World  any    `json:"world,omitempty"`
}

func Run(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /", handleIndex)
	mux.HandleFunc("GET /api/world", handleWorld)
	mux.HandleFunc("POST /api/command", handleCommand)

	fmt.Printf("mini-web listening on http://localhost%s\n", addr)
	return http.ListenAndServe(addr, mux)
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = w.Write([]byte(indexHTML))
}

func handleWorld(w http.ResponseWriter, r *http.Request) {
	currentWorld, err := world.Load()
	if err != nil {
		writeJSON(w, http.StatusOK, map[string]any{
			"error": err.Error(),
		})
		return
	}

	writeJSON(w, http.StatusOK, currentWorld)
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	var request commandRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, commandResponse{Error: err.Error()})
		return
	}

	commandLine := strings.TrimSpace(request.Command)
	if commandLine == "" {
		writeJSON(w, http.StatusBadRequest, commandResponse{Error: "command is required"})
		return
	}

	output, err := commandrun.CaptureStdout(func() error {
		return commands.ExecuteArgs(strings.Fields(commandLine))
	})

	response := commandResponse{
		Output: strings.TrimSpace(output),
	}
	if err != nil {
		response.Error = err.Error()
	}

	currentWorld, loadErr := world.Load()
	if loadErr == nil {
		response.World = currentWorld
	} else if response.Error == "" {
		response.Error = loadErr.Error()
	}

	writeJSON(w, http.StatusOK, response)
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(value); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}

const indexHTML = `<!doctype html>
<html lang="en">
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, initial-scale=1">
  <title>mini-world</title>
  <style>
    :root {
      color-scheme: light;
      --bg: #f7f7f4;
      --panel: #ffffff;
      --line: #d9d9d2;
      --text: #1f2328;
      --muted: #667085;
      --accent: #0f766e;
      --danger: #b42318;
      --code: #f1f3f1;
    }
    * { box-sizing: border-box; }
    body {
      margin: 0;
      background: var(--bg);
      color: var(--text);
      font: 14px/1.4 ui-sans-serif, system-ui, -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;
    }
    header {
      position: sticky;
      top: 0;
      z-index: 2;
      border-bottom: 1px solid var(--line);
      background: rgba(247, 247, 244, 0.96);
      backdrop-filter: blur(8px);
      padding: 12px 16px;
    }
    h1 {
      margin: 0 0 10px;
      font-size: 18px;
      font-weight: 650;
      letter-spacing: 0;
    }
    form {
      display: grid;
      grid-template-columns: 1fr auto;
      gap: 8px;
    }
    input, button {
      height: 38px;
      border: 1px solid var(--line);
      border-radius: 6px;
      font: inherit;
    }
    input {
      padding: 0 10px;
      background: var(--panel);
      color: var(--text);
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    }
    button {
      padding: 0 14px;
      background: var(--accent);
      border-color: var(--accent);
      color: white;
      cursor: pointer;
    }
    main {
      display: grid;
      grid-template-columns: minmax(280px, 420px) minmax(0, 1fr);
      gap: 12px;
      padding: 12px;
    }
    section {
      border: 1px solid var(--line);
      background: var(--panel);
      border-radius: 8px;
      min-width: 0;
    }
    .side {
      display: grid;
      gap: 12px;
      align-content: start;
    }
    .panel-title {
      margin: 0;
      padding: 10px 12px;
      border-bottom: 1px solid var(--line);
      font-size: 13px;
      font-weight: 650;
      color: var(--muted);
      text-transform: uppercase;
    }
    pre {
      margin: 0;
      padding: 12px;
      overflow: auto;
      max-height: 38vh;
      background: var(--code);
      font: 12px/1.45 ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      white-space: pre-wrap;
      word-break: break-word;
    }
    .error {
      color: var(--danger);
    }
    .summary {
      display: grid;
      grid-template-columns: repeat(2, minmax(0, 1fr));
      gap: 8px;
      padding: 12px;
    }
    .metric {
      border: 1px solid var(--line);
      border-radius: 6px;
      padding: 8px;
      background: #fbfbf9;
    }
    .metric strong {
      display: block;
      font-size: 20px;
      line-height: 1.1;
    }
    .metric span {
      color: var(--muted);
      font-size: 12px;
    }
    .world {
      overflow: hidden;
    }
    .main-content {
      display: grid;
      gap: 12px;
      min-width: 0;
      align-content: start;
    }
    .tables {
      overflow: hidden;
    }
    .table-wrap {
      padding: 0 12px 12px;
      overflow: auto;
      max-height: 42vh;
    }
    table {
      width: 100%;
      border-collapse: collapse;
      min-width: 720px;
      font-size: 12px;
    }
    th, td {
      border-bottom: 1px solid var(--line);
      padding: 6px 8px;
      text-align: left;
      vertical-align: top;
      white-space: nowrap;
    }
    th {
      position: sticky;
      top: 0;
      background: #fbfbf9;
      z-index: 1;
      color: var(--muted);
      font-weight: 650;
    }
    td {
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    }
    td.complex {
      white-space: pre;
      max-width: 320px;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .empty {
      padding: 0 12px 12px;
      color: var(--muted);
    }
    details {
      border-top: 1px solid var(--line);
    }
    details:first-of-type {
      border-top: 0;
    }
    summary {
      cursor: pointer;
      padding: 10px 12px;
      font-weight: 650;
      list-style-position: inside;
    }
    .json-tree {
      padding: 0 12px 12px;
      overflow: auto;
      max-height: 72vh;
    }
    .json-tree details {
      border: 0;
      margin-left: 14px;
    }
    .json-tree summary {
      padding: 2px 0;
      font-weight: 500;
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
    }
    .leaf {
      margin-left: 28px;
      font-family: ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
      font-size: 12px;
      padding: 1px 0;
    }
    .key { color: #344054; }
    .value { color: #0f766e; }
    .null { color: #667085; }
    .status {
      margin-top: 8px;
      color: var(--muted);
      font-size: 12px;
    }
    @media (max-width: 900px) {
      main { grid-template-columns: 1fr; }
      pre { max-height: 28vh; }
    }
  </style>
</head>
<body>
  <header>
    <h1>mini-world</h1>
    <form id="command-form">
      <input id="command-input" autocomplete="off" spellcheck="false" placeholder="init, create-bank bank1, deposit-cash alice bank1 EUR 100">
      <button type="submit">Run</button>
    </form>
    <div id="status" class="status">Polling mini_world.json...</div>
  </header>

  <main>
    <div class="side">
      <section>
        <h2 class="panel-title">Summary</h2>
        <div id="summary" class="summary"></div>
      </section>
      <section>
        <h2 class="panel-title">Last Output</h2>
        <pre id="output"></pre>
      </section>
      <section>
        <h2 class="panel-title">Error</h2>
        <pre id="error" class="error"></pre>
      </section>
    </div>

    <div class="main-content">
      <section class="tables">
        <h2 class="panel-title">Entity Tables</h2>
        <div id="tables"></div>
      </section>
      <section class="world">
        <h2 class="panel-title">World JSON</h2>
        <div id="world" class="json-tree"></div>
      </section>
    </div>
  </main>

  <script>
    const form = document.querySelector("#command-form");
    const input = document.querySelector("#command-input");
    const statusEl = document.querySelector("#status");
    const summaryEl = document.querySelector("#summary");
    const outputEl = document.querySelector("#output");
    const errorEl = document.querySelector("#error");
    const tablesEl = document.querySelector("#tables");
    const worldEl = document.querySelector("#world");

    let lastWorldText = "";
    const openSections = new Set(["central_banks", "banks", "humans", "accounts"]);
    const tableKeys = [
      "central_banks",
      "banks",
      "humans",
      "accounts",
      "currencies",
      "assets",
      "reserve_loans",
      "customer_loans",
      "payment_instructions",
      "messages",
      "settlements",
      "correspondent_accounts",
      "fx_markets",
      "command_history"
    ];

    form.addEventListener("submit", async (event) => {
      event.preventDefault();
      const command = input.value.trim();
      if (!command) return;

      const response = await fetch("/api/command", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ command })
      });
      const result = await response.json();
      outputEl.textContent = result.output || "";
      errorEl.textContent = result.error || "";
      input.value = "";
      if (result.world) renderWorld(result.world);
    });

    async function pollWorld() {
      try {
        const response = await fetch("/api/world", { cache: "no-store" });
        const text = await response.text();
        if (text !== lastWorldText) {
          lastWorldText = text;
          renderWorld(JSON.parse(text));
        }
        statusEl.textContent = "Live: " + new Date().toLocaleTimeString();
      } catch (error) {
        statusEl.textContent = "Polling error: " + error.message;
      }
    }

    function renderWorld(world) {
      if (world.error) {
        summaryEl.innerHTML = "";
        tablesEl.textContent = world.error;
        worldEl.textContent = world.error;
        return;
      }

      renderSummary(world);
      renderTables(world);
      worldEl.replaceChildren(renderObject("world", world, true));
    }

    function renderSummary(world) {
      summaryEl.replaceChildren(...tableKeys.map((key) => {
        const value = world[key];
        const count = Array.isArray(value) ? value.length : value && typeof value === "object" ? Object.keys(value).length : 0;
        const item = document.createElement("div");
        item.className = "metric";
        item.innerHTML = "<strong>" + count + "</strong><span>" + key + "</span>";
        return item;
      }));
    }

    function renderTables(world) {
      const sections = tableKeys.map((key) => renderTableSection(key, world[key]));
      tablesEl.replaceChildren(...sections);
    }

    function renderTableSection(name, value) {
      const details = document.createElement("details");
      details.open = openSections.has(name);
      details.addEventListener("toggle", () => {
        if (details.open) openSections.add(name);
        else openSections.delete(name);
      });

      const rows = tableRows(value);
      const summary = document.createElement("summary");
      summary.textContent = name + " (" + rows.length + ")";
      details.appendChild(summary);

      if (rows.length === 0) {
        const empty = document.createElement("div");
        empty.className = "empty";
        empty.textContent = "No rows";
        details.appendChild(empty);
        return details;
      }

      const columns = tableColumns(rows);
      const wrap = document.createElement("div");
      wrap.className = "table-wrap";
      const table = document.createElement("table");
      const thead = document.createElement("thead");
      const headRow = document.createElement("tr");
      for (const column of columns) {
        const th = document.createElement("th");
        th.textContent = column;
        headRow.appendChild(th);
      }
      thead.appendChild(headRow);
      table.appendChild(thead);

      const tbody = document.createElement("tbody");
      for (const row of rows) {
        const tr = document.createElement("tr");
        for (const column of columns) {
          const td = document.createElement("td");
          const cellValue = row[column];
          td.textContent = tableCellValue(cellValue);
          if (cellValue && typeof cellValue === "object") td.className = "complex";
          tr.appendChild(td);
        }
        tbody.appendChild(tr);
      }
      table.appendChild(tbody);
      wrap.appendChild(table);
      details.appendChild(wrap);
      return details;
    }

    function tableRows(value) {
      if (Array.isArray(value)) {
        return value.map((item, index) => normalizeRow(index, item));
      }
      if (value && typeof value === "object") {
        return Object.entries(value).map(([key, item]) => normalizeRow(key, item));
      }
      return [];
    }

    function normalizeRow(key, value) {
      if (value && typeof value === "object" && !Array.isArray(value)) {
        return { id: key, ...value };
      }
      return { id: key, value };
    }

    function tableColumns(rows) {
      const columns = ["id"];
      for (const row of rows) {
        for (const key of Object.keys(row)) {
          if (!columns.includes(key)) columns.push(key);
        }
      }
      return columns;
    }

    function tableCellValue(value) {
      if (value === undefined) return "";
      if (value === null) return "null";
      if (typeof value === "object") return JSON.stringify(value, null, 2);
      return String(value);
    }

    function renderObject(name, value, open) {
      if (value === null || typeof value !== "object") {
        const leaf = document.createElement("div");
        leaf.className = "leaf";
        leaf.innerHTML = '<span class="key">' + escapeHTML(name) + '</span>: ' + formatValue(value);
        return leaf;
      }

      const details = document.createElement("details");
      details.open = open;
      const summary = document.createElement("summary");
      const count = Array.isArray(value) ? value.length : Object.keys(value).length;
      summary.textContent = name + " (" + count + ")";
      details.appendChild(summary);

      const entries = Array.isArray(value) ? value.entries() : Object.entries(value);
      for (const [key, child] of entries) {
        details.appendChild(renderObject(String(key), child, false));
      }

      return details;
    }

    function formatValue(value) {
      if (value === null) return '<span class="null">null</span>';
      return '<span class="value">' + escapeHTML(JSON.stringify(value)) + '</span>';
    }

    function escapeHTML(value) {
      return String(value)
        .replaceAll("&", "&amp;")
        .replaceAll("<", "&lt;")
        .replaceAll(">", "&gt;")
        .replaceAll('"', "&quot;");
    }

    pollWorld();
    setInterval(pollWorld, 400);
  </script>
</body>
</html>`
