"use client";

import { useEffect, useState } from "react";
import {
  Activity,
  AlertTriangle,
  CheckCircle2,
  Clock,
  GitCommit,
  GitPullRequest,
  Server,
  Users,
  XCircle,
} from "lucide-react";

type RepoStatus = {
  repo: string;
  human_commits_7d: number;
  last_human_commit: {
    sha: string;
    date: string;
    author: string;
    msg: string;
  } | null;
  contributors_7d: { login: string; commits: number }[];
  open_prs: number;
  open_pr_titles: string[];
  last_workflow: {
    name: string;
    conclusion: string | null;
    branch: string;
    at: string;
  } | null;
  has_claude_yml: boolean;
  workflow_files: string[];
  agents_md_unfilled: boolean | null;
};

type StatusFile = {
  generated_at: string;
  repos: Record<string, RepoStatus>;
};

function formatRelative(iso: string): string {
  const d = new Date(iso);
  const sec = Math.floor((Date.now() - d.getTime()) / 1000);
  if (sec < 60) return `${sec}s ago`;
  if (sec < 3600) return `${Math.floor(sec / 60)}m ago`;
  if (sec < 86400) return `${Math.floor(sec / 3600)}h ago`;
  return `${Math.floor(sec / 86400)}d ago`;
}

function CiBadge({ conclusion }: { conclusion: string | null }) {
  if (conclusion === "success") {
    return (
      <span className="inline-flex items-center gap-1 text-emerald-400">
        <CheckCircle2 className="w-3.5 h-3.5" />
        ok
      </span>
    );
  }
  if (conclusion === "failure") {
    return (
      <span className="inline-flex items-center gap-1 text-rose-400">
        <XCircle className="w-3.5 h-3.5" />
        fail
      </span>
    );
  }
  return <span className="text-zinc-500">—</span>;
}

function YesNo({ yes }: { yes: boolean }) {
  return yes ? (
    <span className="text-emerald-400">✓</span>
  ) : (
    <span className="text-rose-400/70">✗</span>
  );
}

export default function ClusterStatusPage() {
  const [data, setData] = useState<StatusFile | null>(null);
  const [error, setError] = useState<string | null>(null);
  const [filter, setFilter] = useState<"all" | "active" | "idle" | "issues">("all");

  useEffect(() => {
    fetch("/cluster-status.json", { cache: "no-store" })
      .then((r) => {
        if (!r.ok) throw new Error(`status ${r.status}`);
        return r.json();
      })
      .then((d: StatusFile) => setData(d))
      .catch((e) => setError(String(e)));
  }, []);

  if (error) {
    return (
      <main className="p-8 max-w-6xl mx-auto">
        <h1 className="text-2xl font-semibold mb-4 text-zinc-100">Cluster Status</h1>
        <div className="bg-rose-950/40 border border-rose-900 text-rose-200 rounded-md p-4">
          Failed to load <code>/cluster-status.json</code>: {error}
        </div>
      </main>
    );
  }

  if (!data) {
    return (
      <main className="p-8 max-w-6xl mx-auto">
        <div className="text-zinc-500">Loading cluster status…</div>
      </main>
    );
  }

  const repos = Object.values(data.repos).sort((a, b) =>
    b.human_commits_7d - a.human_commits_7d
  );

  const active = repos.filter((r) => r.human_commits_7d > 0);
  const idle = repos.filter((r) => r.human_commits_7d === 0);
  const failingCi = repos.filter(
    (r) => r.last_workflow?.conclusion === "failure"
  );
  const noClaude = repos.filter((r) => !r.has_claude_yml);
  const unfilled = repos.filter((r) => r.agents_md_unfilled === true);
  const withPrs = repos.filter((r) => r.open_prs > 0);

  // Aggregate contributors
  const contribMap = new Map<string, { commits: number; repos: Set<string> }>();
  for (const r of repos) {
    for (const c of r.contributors_7d) {
      if (!contribMap.has(c.login)) {
        contribMap.set(c.login, { commits: 0, repos: new Set() });
      }
      const entry = contribMap.get(c.login)!;
      entry.commits += c.commits;
      entry.repos.add(r.repo);
    }
  }
  const contributors = Array.from(contribMap.entries())
    .map(([login, v]) => ({ login, commits: v.commits, repos: v.repos.size }))
    .sort((a, b) => b.commits - a.commits);

  const visible =
    filter === "active"
      ? active
      : filter === "idle"
      ? idle
      : filter === "issues"
      ? repos.filter(
          (r) =>
            r.last_workflow?.conclusion === "failure" ||
            !r.has_claude_yml ||
            r.agents_md_unfilled === true
        )
      : repos;

  return (
    <main className="p-6 md:p-8 max-w-7xl mx-auto">
      <header className="mb-6">
        <div className="flex items-center gap-3 mb-2">
          <Server className="w-6 h-6 text-violet-400" />
          <h1 className="text-2xl font-semibold text-zinc-100">Cluster Status</h1>
        </div>
        <p className="text-sm text-zinc-500">
          Last 7 days of dev activity across {repos.length} dclaw-* app repos. Snapshot
          generated <span className="text-zinc-300">{formatRelative(data.generated_at)}</span>{" "}
          (
          {new Date(data.generated_at).toLocaleString(undefined, {
            year: "numeric",
            month: "short",
            day: "numeric",
            hour: "2-digit",
            minute: "2-digit",
          })}
          ).
        </p>
      </header>

      {/* Summary tiles */}
      <section className="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-6 gap-3 mb-8">
        <Tile
          icon={<Activity className="w-4 h-4 text-emerald-400" />}
          label="Active (7d)"
          value={active.length}
          sub={`${Math.round((active.length / repos.length) * 100)}% of ${repos.length}`}
          onClick={() => setFilter("active")}
        />
        <Tile
          icon={<Clock className="w-4 h-4 text-zinc-500" />}
          label="Idle"
          value={idle.length}
          sub="no human commit in 7d"
          onClick={() => setFilter("idle")}
        />
        <Tile
          icon={<GitPullRequest className="w-4 h-4 text-violet-400" />}
          label="Open PRs"
          value={withPrs.length}
          sub={`${repos.reduce((s, r) => s + r.open_prs, 0)} total`}
        />
        <Tile
          icon={<XCircle className="w-4 h-4 text-rose-400" />}
          label="CI failing"
          value={failingCi.length}
          sub="latest workflow run"
        />
        <Tile
          icon={<AlertTriangle className="w-4 h-4 text-amber-400" />}
          label="No Claude"
          value={noClaude.length}
          sub="missing workflow"
        />
        <Tile
          icon={<AlertTriangle className="w-4 h-4 text-amber-400" />}
          label="Unfilled"
          value={unfilled.length}
          sub="AGENTS.md placeholders"
        />
      </section>

      {/* Filter pills */}
      <div className="flex flex-wrap gap-2 mb-4">
        {(["all", "active", "idle", "issues"] as const).map((f) => (
          <button
            key={f}
            onClick={() => setFilter(f)}
            className={`text-xs px-3 py-1.5 rounded-full border transition-colors ${
              filter === f
                ? "bg-violet-600/20 border-violet-500/50 text-violet-200"
                : "border-zinc-800 text-zinc-400 hover:text-zinc-100 hover:border-zinc-700"
            }`}
          >
            {f}
          </button>
        ))}
        <span className="text-xs text-zinc-500 self-center ml-auto">
          showing {visible.length} of {repos.length}
        </span>
      </div>

      {/* Repos table */}
      <section className="bg-zinc-950 border border-zinc-800 rounded-lg overflow-hidden mb-8">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="bg-zinc-900/60 text-xs uppercase tracking-wider text-zinc-500">
              <tr>
                <th className="text-left py-2 px-3">Repo</th>
                <th className="text-left py-2 px-3">Owner (last 7d)</th>
                <th className="text-right py-2 px-3">Commits</th>
                <th className="text-left py-2 px-3">Last human commit</th>
                <th className="text-right py-2 px-3">PRs</th>
                <th className="text-left py-2 px-3">CI</th>
                <th className="text-center py-2 px-3">Claude</th>
                <th className="text-center py-2 px-3">AGENTS</th>
              </tr>
            </thead>
            <tbody>
              {visible.map((r) => {
                const owner = r.contributors_7d[0]?.login || "—";
                const lhc = r.last_human_commit;
                return (
                  <tr
                    key={r.repo}
                    className="border-t border-zinc-900 hover:bg-zinc-900/40"
                  >
                    <td className="py-2 px-3 font-mono text-xs text-zinc-200">
                      <a
                        href={`https://github.com/dclawstack/${r.repo}`}
                        target="_blank"
                        rel="noreferrer"
                        className="hover:text-violet-300"
                      >
                        {r.repo}
                      </a>
                    </td>
                    <td className="py-2 px-3">
                      {owner === "—" ? (
                        <span className="text-zinc-600">—</span>
                      ) : (
                        <span className="text-zinc-200 font-mono text-xs">
                          @{owner}
                        </span>
                      )}
                    </td>
                    <td className="py-2 px-3 text-right text-zinc-300 tabular-nums">
                      {r.human_commits_7d}
                    </td>
                    <td className="py-2 px-3 text-zinc-400 text-xs">
                      {lhc ? (
                        <>
                          <span className="text-zinc-500">{lhc.date}</span>{" "}
                          <span className="text-zinc-300">{lhc.msg.slice(0, 40)}</span>
                        </>
                      ) : (
                        <span className="text-zinc-600">—</span>
                      )}
                    </td>
                    <td className="py-2 px-3 text-right tabular-nums">
                      {r.open_prs > 0 ? (
                        <span className="text-violet-300">{r.open_prs}</span>
                      ) : (
                        <span className="text-zinc-700">0</span>
                      )}
                    </td>
                    <td className="py-2 px-3">
                      <CiBadge conclusion={r.last_workflow?.conclusion || null} />
                    </td>
                    <td className="py-2 px-3 text-center">
                      <YesNo yes={r.has_claude_yml} />
                    </td>
                    <td className="py-2 px-3 text-center">
                      {r.agents_md_unfilled === true ? (
                        <span
                          title="AGENTS.md still has {APP_NAME} placeholders"
                          className="text-amber-400"
                        >
                          !
                        </span>
                      ) : r.agents_md_unfilled === false ? (
                        <span className="text-emerald-400">✓</span>
                      ) : (
                        <span className="text-zinc-700">?</span>
                      )}
                    </td>
                  </tr>
                );
              })}
            </tbody>
          </table>
        </div>
      </section>

      {/* Contributors */}
      <section className="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        <div className="bg-zinc-950 border border-zinc-800 rounded-lg p-4">
          <h3 className="text-sm font-semibold text-zinc-200 mb-3 flex items-center gap-2">
            <Users className="w-4 h-4 text-violet-400" />
            Contributors (last 7d)
          </h3>
          <ul className="space-y-2">
            {contributors.slice(0, 15).map((c) => (
              <li key={c.login} className="flex items-center justify-between text-xs">
                <span className="font-mono text-zinc-200">@{c.login}</span>
                <span className="text-zinc-500">
                  <span className="text-zinc-300 tabular-nums">{c.commits}</span> commits
                  in {c.repos} {c.repos === 1 ? "repo" : "repos"}
                </span>
              </li>
            ))}
          </ul>
        </div>

        <div className="bg-zinc-950 border border-zinc-800 rounded-lg p-4">
          <h3 className="text-sm font-semibold text-zinc-200 mb-3 flex items-center gap-2">
            <GitCommit className="w-4 h-4 text-violet-400" />
            Tomorrow&apos;s readiness
          </h3>
          <ul className="space-y-2 text-xs text-zinc-400">
            <li className="flex items-start gap-2">
              <span className="text-emerald-400 mt-0.5">●</span>
              <span>
                <span className="text-zinc-200">{active.length}/{repos.length}</span> repos saw
                human commits in last 7d — strong engagement.
              </span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-amber-400 mt-0.5">●</span>
              <span>
                <span className="text-zinc-200">{noClaude.length}</span> repos missing Claude
                workflow — those devs can&apos;t use Claude Action tomorrow.
              </span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-amber-400 mt-0.5">●</span>
              <span>
                <span className="text-zinc-200">{unfilled.length}</span> repos still have
                <code className="px-1 mx-0.5 bg-zinc-800 rounded text-zinc-300">
                  {"{APP_NAME}"}
                </code>
                placeholders in <code>AGENTS.md</code>.
              </span>
            </li>
            <li className="flex items-start gap-2">
              <span className="text-rose-400 mt-0.5">●</span>
              <span>
                <span className="text-zinc-200">{failingCi.length}</span> repos have a
                failing latest CI run — likely caused by today&apos;s scaffold sync;
                investigate before tomorrow.
              </span>
            </li>
          </ul>
          <div className="mt-4 pt-3 border-t border-zinc-800 text-xs text-zinc-500">
            <p>
              Snapshot file: <code className="text-zinc-400">/cluster-status.json</code>.
              Regenerate with{" "}
              <code className="text-zinc-400">
                dclaw-platform/scripts/refresh-cluster-status.sh
              </code>
              .
            </p>
          </div>
        </div>
      </section>
    </main>
  );
}

function Tile({
  icon,
  label,
  value,
  sub,
  onClick,
}: {
  icon: React.ReactNode;
  label: string;
  value: number | string;
  sub: string;
  onClick?: () => void;
}) {
  const Wrapper = onClick ? "button" : "div";
  return (
    <Wrapper
      onClick={onClick}
      className={`bg-zinc-950 border border-zinc-800 rounded-lg p-3 text-left ${
        onClick ? "hover:border-zinc-700 transition-colors" : ""
      }`}
    >
      <div className="flex items-center gap-2 text-xs text-zinc-400 mb-1">
        {icon}
        {label}
      </div>
      <div className="text-2xl font-semibold text-zinc-100 tabular-nums">{value}</div>
      <div className="text-xs text-zinc-500 mt-0.5">{sub}</div>
    </Wrapper>
  );
}
