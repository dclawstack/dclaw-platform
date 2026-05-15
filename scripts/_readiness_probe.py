#!/usr/bin/env python3
"""Probe dclawstack org for per-repo readiness signals.

Output: /tmp/readiness.json (consumed by refresh-cluster-status.sh).
For each dclaw-* app repo, gathers:
- last commits (human + total)
- top contributor in the last 7d
- open PR count
- last workflow run conclusion
- Claude workflow presence
- AGENTS.md placeholder check (`{APP_NAME}` etc.)

Run via `refresh-cluster-status.sh`, not directly.
"""

import base64
import concurrent.futures as cf
import datetime as dt
import json
import subprocess
import sys
from pathlib import Path

REPOS_FILE = "/tmp/dclaw_app_repos_only.txt"
OUT_JSON = "/tmp/readiness.json"
ORG = "dclawstack"
SINCE_DAYS = 7

now = dt.datetime.now(dt.timezone.utc)
since = (now - dt.timedelta(days=SINCE_DAYS)).strftime("%Y-%m-%dT%H:%M:%SZ")

BOT_AUTHORS = {
    "github-actions[bot]",
    "dependabot[bot]",
    "claude[bot]",
    None,
    "renovate[bot]",
}
SCAFFOLD_MSG_PREFIX = "chore(scaffold): sync"


def gh(args):
    try:
        r = subprocess.run(
            ["gh", "api", *args], capture_output=True, text=True, timeout=15
        )
        if r.returncode != 0:
            return None
        return json.loads(r.stdout) if r.stdout.strip() else None
    except Exception:
        return None


def fetch_repo(repo):
    info = {"repo": repo}

    commits = gh([f"repos/{ORG}/{repo}/commits?since={since}&per_page=50"]) or []
    human_commits = []
    for c in commits:
        author_login = (c.get("author") or {}).get("login")
        msg = c["commit"]["message"].split("\n")[0]
        if author_login in BOT_AUTHORS:
            continue
        if msg.startswith(SCAFFOLD_MSG_PREFIX):
            continue
        human_commits.append({
            "sha": c["sha"][:7],
            "date": c["commit"]["author"]["date"][:10],
            "author": author_login or c["commit"]["author"]["name"],
            "msg": msg[:60],
        })

    info["human_commits_7d"] = len(human_commits)
    info["last_human_commit"] = human_commits[0] if human_commits else None

    counter = {}
    for c in human_commits:
        counter[c["author"]] = counter.get(c["author"], 0) + 1
    info["contributors_7d"] = [
        {"login": a, "commits": n}
        for a, n in sorted(counter.items(), key=lambda x: -x[1])
    ]

    prs = gh([f"repos/{ORG}/{repo}/pulls?state=open&per_page=10"]) or []
    info["open_prs"] = len(prs)
    info["open_pr_titles"] = [f"#{p['number']} {p['title'][:50]}" for p in prs[:3]]

    runs = gh([f"repos/{ORG}/{repo}/actions/runs?per_page=5"]) or {}
    workflow_runs = runs.get("workflow_runs", []) if isinstance(runs, dict) else []
    if workflow_runs:
        wr = workflow_runs[0]
        info["last_workflow"] = {
            "name": wr.get("name"),
            "conclusion": wr.get("conclusion") or wr.get("status"),
            "branch": wr.get("head_branch"),
            "at": (wr.get("created_at") or "")[:10],
        }
    else:
        info["last_workflow"] = None

    wfs = gh([f"repos/{ORG}/{repo}/contents/.github/workflows"]) or []
    if isinstance(wfs, list):
        names = [w.get("name", "") for w in wfs]
        info["has_claude_yml"] = any("claude" in n.lower() for n in names)
        info["workflow_files"] = [n for n in names if n.endswith(".yml")]
    else:
        info["has_claude_yml"] = False
        info["workflow_files"] = []

    agents = gh([f"repos/{ORG}/{repo}/contents/AGENTS.md"])
    if agents and "content" in agents:
        try:
            txt = base64.b64decode(agents["content"]).decode("utf-8", "ignore")
            info["agents_md_unfilled"] = "{APP_NAME}" in txt or "{BACKEND_PORT}" in txt
        except Exception:
            info["agents_md_unfilled"] = None
    else:
        info["agents_md_unfilled"] = None

    return info


def main():
    repos = [l.strip() for l in Path(REPOS_FILE).read_text().splitlines() if l.strip()]
    print(f"  probing {len(repos)} repos…", file=sys.stderr)
    results = {}
    with cf.ThreadPoolExecutor(max_workers=8) as ex:
        for info in ex.map(fetch_repo, repos):
            results[info["repo"]] = info
            print(".", end="", flush=True, file=sys.stderr)
    print("", file=sys.stderr)
    Path(OUT_JSON).write_text(json.dumps(results, indent=2))
    print(f"  wrote {OUT_JSON}", file=sys.stderr)


if __name__ == "__main__":
    main()
