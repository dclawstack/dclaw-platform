"use client";

import { useMemo } from "react";

function extractHeadings(content: string): { level: number; text: string; id: string }[] {
  const headings: { level: number; text: string; id: string }[] = [];
  const lines = content.split("\n");

  for (const line of lines) {
    const match = line.match(/^(#{2,3})\s+(.+)$/);
    if (match) {
      const level = match[1].length;
      const text = match[2].trim();
      const id = text
        .toLowerCase()
        .replace(/[^\w\s-]/g, "")
        .replace(/\s+/g, "-")
        .replace(/-+/g, "-");
      headings.push({ level, text, id });
    }
  }

  return headings;
}

export function DocsToc({ content }: { content: string }) {
  const headings = useMemo(() => extractHeadings(content), [content]);

  if (headings.length === 0) {
    return (
      <div className="text-xs text-zinc-600">
        No sections on this page.
      </div>
    );
  }

  return (
    <div>
      <h4 className="text-xs font-semibold text-zinc-500 uppercase tracking-wider mb-3">
        On this page
      </h4>
      <nav className="space-y-1">
        {headings.map((h, i) => (
          <a
            key={i}
            href={`#${h.id}`}
            className={`
              block text-xs transition-colors hover:text-zinc-200
              ${h.level === 2 ? "text-zinc-400 font-medium" : "text-zinc-600 pl-3"}
            `}
          >
            {h.text}
          </a>
        ))}
      </nav>
    </div>
  );
}
