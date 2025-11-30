import React, { useEffect, useMemo, useState } from "react";
import "./App.css";

const API = "http://localhost:8080";

export default function App() {
  const [displayName, setDisplayName] = useState("Sumera Malik Blockchain");
  const [chain, setChain] = useState([]);
  const [pending, setPending] = useState([]);
  const [tx, setTx] = useState("");
  const [search, setSearch] = useState("");
  const [results, setResults] = useState([]);
  const [mining, setMining] = useState(false);
  const [msg, setMsg] = useState("");
  const [apiOK, setApiOK] = useState(true);
  const [loading, setLoading] = useState(true);

  const hasBlocks = useMemo(() => (Array.isArray(chain) && chain.length > 0), [chain]);

  async function refresh() {
    try {
      const v = await fetch(`${API}/view`).then(r => r.json());
      setDisplayName(v.displayName || "Sumera Malik Blockchain");
      setChain(v.blocks || []);
      setPending(v.pendingTx || []);
      setApiOK(true);
    } catch (e) {
      setApiOK(false);
    } finally {
      setLoading(false);
    }
  }

  // Initial load + light auto-refresh every 5s
  useEffect(() => {
    refresh();
    const t = setInterval(refresh, 5000);
    return () => clearInterval(t);
  }, []);

  // Add transaction
  async function addTx() {
    const val = tx.trim();
    if (!val) return;
    try {
      const res = await fetch(`${API}/tx`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ data: val })
      }).then(r => r.json());
      setMsg(`Added: ${res.added}`);
      setTx("");
      await refresh();
    } catch {
      setMsg("Could not add transaction (server offline?)");
    }
  }

  // Mine block
  async function mine() {
    setMining(true);
    setMsg("Mining… please wait");
    try {
      const res = await fetch(`${API}/mine`, { method: "POST" }).then(r => r.json());
      if (res && res.ok && res.block) {
        setMsg(`Mined block #${res.block.index} `);
      } else {
        setMsg((res && res.error) || "Mining failed");
      }
    } catch {
      setMsg("Mining failed (server offline?)");
    } finally {
      setMining(false);
      await refresh();
    }
  }

  // Search
  async function doSearch() {
    const q = search.trim();
    if (!q) { setResults([]); return; }
    try {
      const res = await fetch(`${API}/search?q=${encodeURIComponent(q)}`).then(r => r.json());
      setResults(res.results || []);
    } catch {
      setResults([]);
    }
  }

  return (
    <div className="wrap">
      <h1> {displayName}</h1>

      {!apiOK && (
        <div className="banner error">
          Backend not reachable at <code>{API}</code>. Make sure it’s running:
          <code>go run .</code> in <b>C:\SumeraBlockchain\backend</b>
        </div>
      )}

      {loading && <div className="banner info">Loading blockchain…</div>}

      <div className="grid2">
        <section>
          <h3>Add Transaction</h3>
          <input
            value={tx}
            onChange={(e)=>setTx(e.target.value)}
            placeholder="Enter transaction text (e.g., 'Alice pays Bob 5')"
          />
          <div className="row">
            <button onClick={addTx}>Add</button>
            <div className="muted">Pending: {pending.length}</div>
          </div>
        </section>

        <section>
          <h3>Mine Block (Proof of Work)</h3>
          <button onClick={mine} disabled={mining}>{mining ? "Mining…" : "Mine"}</button>
          <div className="msg">{msg}</div>
        </section>
      </div>

      <section>
        <h3>Search Transactions</h3>
        <input
          value={search}
          onChange={(e)=>setSearch(e.target.value)}
          placeholder="Type text to find (e.g., 21i-1579)"
        />
        <button onClick={doSearch} style={{marginTop: 8}}>Search</button>
        {results.length > 0 && (
          <ul className="hits">
            {results.map((r, i) => (
              <li key={i}>
                In block <b>#{r.blockIndex}</b> — matches: {r.matches.join(", ")}
              </li>
            ))}
          </ul>
        )}
      </section>

      <section>
        <h2>Blockchain</h2>
        {!hasBlocks && <div className="muted">No blocks yet.</div>}
        <ul className="cards">
          {chain.map(b => (
            <li key={b.index} className="card">
              <div className="cardHead">
                <div><b>Block #{b.index}</b></div>
                <div className="muted">{b.timestamp}</div>
              </div>
              <div><span className="k">Hash:</span> <code>{b.hash}</code></div>
              <div><span className="k">Prev:</span> <code>{b.prevHash}</code></div>
              <div className="row">
                <div><span className="k">Nonce:</span> {b.nonce}</div>
                <div><span className="k">MerkleRoot:</span> <code>{b.merkleRoot}</code></div>
              </div>
              <div className="data">
                <b>Data:</b>
                <ul>
                  {b.data.map((d, idx) => <li key={idx}>{d}</li>)}
                </ul>
              </div>
            </li>
          ))}
        </ul>
      </section>

      <footer>
        Owner: <b>Sumera Malik</b> • First TX: <code>21i-1579</code>
      </footer>
    </div>
  );
}
