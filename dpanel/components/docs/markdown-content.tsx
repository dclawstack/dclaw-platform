"use client";

import { useMemo } from "react";

function slugify(text: string): string {
  return text
    .toLowerCase()
    .replace(/[^\w\s-]/g, "")
    .replace(/\s+/g, "-")
    .replace(/-+/g, "-");
}

function parseMarkdown(content: string): React.ReactNode[] {
  const lines = content.split("\n");
  const elements: React.ReactNode[] = [];
  let i = 0;
  let listItems: string[] = [];
  let listType: "ul" | "ol" | null = null;
  let codeBlock: string[] = [];
  let codeLang = "";
  let inCodeBlock = false;

  function flushList() {
    if (listItems.length === 0) return;
    const items = listItems.map((item, idx) => (
      <li key={idx} className="text-zinc-300 leading-relaxed">
        {renderInline(item)}
      </li>
    ));
    if (listType === "ul") {
      elements.push(
        <ul key={`ul-${i}`} className="list-disc list-inside space-y-1 my-4 text-zinc-300">
          {items}
        </ul>
      );
    } else {
      elements.push(
        <ol key={`ol-${i}`} className="list-decimal list-inside space-y-1 my-4 text-zinc-300">
          {items}
        </ol>
      );
    }
    listItems = [];
    listType = null;
  }

  function flushCodeBlock() {
    if (codeBlock.length === 0) return;
    const code = codeBlock.join("\n");
    elements.push(
      <pre
        key={`pre-${i}`}
        className="my-4 p-4 rounded-lg bg-zinc-900 border border-zinc-800 overflow-x-auto"
      >
        {codeLang && (
          <div className="text-[10px] text-zinc-600 uppercase tracking-wider mb-2">{codeLang}</div>
        )}
        <code className="text-sm text-zinc-300 font-mono whitespace-pre">{code}</code>
      </pre>
    );
    codeBlock = [];
    codeLang = "";
    inCodeBlock = false;
  }

  while (i < lines.length) {
    const line = lines[i];

    if (line.startsWith("```")) {
      if (inCodeBlock) {
        flushCodeBlock();
      } else {
        flushList();
        inCodeBlock = true;
        codeLang = line.slice(3).trim();
      }
      i++;
      continue;
    }

    if (inCodeBlock) {
      codeBlock.push(line);
      i++;
      continue;
    }

    if (line.startsWith("# ")) {
      flushList();
      const text = line.slice(2);
      const id = slugify(text);
      elements.push(
        <h1 key={`h1-${i}`} id={id} className="text-3xl font-bold text-zinc-100 mt-8 mb-4">
          {text}
        </h1>
      );
      i++;
      continue;
    }

    if (line.startsWith("## ")) {
      flushList();
      const text = line.slice(3);
      const id = slugify(text);
      elements.push(
        <h2 key={`h2-${i}`} id={id} className="text-2xl font-semibold text-zinc-200 mt-8 mb-3">
          {text}
        </h2>
      );
      i++;
      continue;
    }

    if (line.startsWith("### ")) {
      flushList();
      const text = line.slice(4);
      const id = slugify(text);
      elements.push(
        <h3 key={`h3-${i}`} id={id} className="text-xl font-medium text-zinc-300 mt-6 mb-2">
          {text}
        </h3>
      );
      i++;
      continue;
    }

    if (line.startsWith("- ") || line.startsWith("* ")) {
      if (listType === "ol") flushList();
      listType = "ul";
      listItems.push(line.slice(2));
      i++;
      continue;
    }

    if (/^\d+\.\s/.test(line)) {
      if (listType === "ul") flushList();
      listType = "ol";
      listItems.push(line.replace(/^\d+\.\s/, ""));
      i++;
      continue;
    }

    if (line.trim() === "") {
      flushList();
      i++;
      continue;
    }

    // Regular paragraph
    flushList();

    // Handle table rows
    if (line.startsWith("|")) {
      const tableLines: string[] = [line];
      i++;
      while (i < lines.length && lines[i].startsWith("|")) {
        tableLines.push(lines[i]);
        i++;
      }
      // Skip separator row (---)
      const dataRows = tableLines.filter((l) => !/^\|[-:\s|]+\|$/.test(l));
      if (dataRows.length > 0) {
        elements.push(
          <div key={`table-${i}`} className="my-4 overflow-x-auto">
            <table className="w-full text-sm border-collapse">
              <tbody>
                {dataRows.map((row, rIdx) => {
                  const cells = row
                    .split("|")
                    .filter((c) => c.trim() !== "")
                    .map((c) => c.trim());
                  return (
                    <tr
                      key={rIdx}
                      className={`border-b border-zinc-800 ${
                        rIdx === 0 ? "bg-zinc-900/50 font-medium text-zinc-200" : "text-zinc-400"
                      }`}
                    >
                      {cells.map((cell, cIdx) => (
                        <td key={cIdx} className="px-3 py-2">
                          {renderInline(cell)}
                        </td>
                      ))}
                    </tr>
                  );
                })}
              </tbody>
            </table>
          </div>
        );
      }
      continue;
    }

    // Blockquote
    if (line.startsWith("> ")) {
      const quoteLines: string[] = [line.slice(2)];
      i++;
      while (i < lines.length && lines[i].startsWith("> ")) {
        quoteLines.push(lines[i].slice(2));
        i++;
      }
      elements.push(
        <blockquote
          key={`bq-${i}`}
          className="my-4 pl-4 border-l-2 border-zinc-700 text-zinc-400 italic"
        >
          {quoteLines.join(" ")}
        </blockquote>
      );
      continue;
    }

    elements.push(
      <p key={`p-${i}`} className="text-zinc-300 leading-relaxed my-3">
        {renderInline(line)}
      </p>
    );
    i++;
  }

  flushList();
  flushCodeBlock();

  return elements;
}

function renderInline(text: string): React.ReactNode {
  const parts: React.ReactNode[] = [];
  let remaining = text;
  let key = 0;

  // Links [text](url)
  const linkRegex = /\[([^\]]+)\]\(([^)]+)\)/g;
  let lastIndex = 0;
  let match;

  while ((match = linkRegex.exec(text)) !== null) {
    if (match.index > lastIndex) {
      parts.push(renderInlineStyles(text.slice(lastIndex, match.index), key++));
    }
    parts.push(
      <a
        key={key++}
        href={match[2]}
        target="_blank"
        rel="noopener noreferrer"
        className="text-blue-400 hover:text-blue-300 underline"
      >
        {match[1]}
      </a>
    );
    lastIndex = linkRegex.lastIndex;
  }

  if (lastIndex < text.length) {
    parts.push(renderInlineStyles(text.slice(lastIndex), key++));
  }

  if (parts.length === 0) {
    return renderInlineStyles(text, 0);
  }

  return parts;
}

function renderInlineStyles(text: string, key: number): React.ReactNode {
  // Bold **text**
  const boldRegex = /\*\*([^*]+)\*\*/g;
  const italicRegex = /\*([^*]+)\*/g;
  const codeRegex = /`([^`]+)`/g;

  // Simple approach: split by all inline patterns
  const segments: { type: "text" | "bold" | "italic" | "code"; content: string }[] = [];
  let remaining = text;

  // This is a simplified inline parser
  // For P0, we handle basic patterns sequentially
  const patterns = [
    { regex: /\*\*([^*]+)\*\*/g, type: "bold" as const },
    { regex: /`([^`]+)`/g, type: "code" as const },
    { regex: /\*([^*]+)\*/g, type: "italic" as const },
  ];

  // A more robust approach: use a single pass with all patterns
  const allMatches: { index: number; end: number; type: string; content: string }[] = [];

  for (const pattern of patterns) {
    const regex = new RegExp(pattern.regex.source, pattern.regex.flags);
    let m;
    while ((m = regex.exec(text)) !== null) {
      allMatches.push({
        index: m.index,
        end: m.index + m[0].length,
        type: pattern.type,
        content: m[1],
      });
    }
  }

  // Sort by index and remove overlaps
  allMatches.sort((a, b) => a.index - b.index);
  const filtered: typeof allMatches = [];
  let lastEnd = -1;
  for (const match of allMatches) {
    if (match.index >= lastEnd) {
      filtered.push(match);
      lastEnd = match.end;
    }
  }

  const result: React.ReactNode[] = [];
  let pos = 0;
  let k = key * 1000;

  for (const match of filtered) {
    if (match.index > pos) {
      result.push(<span key={k++}>{text.slice(pos, match.index)}</span>);
    }
    if (match.type === "bold") {
      result.push(
        <strong key={k++} className="font-semibold text-zinc-200">
          {match.content}
        </strong>
      );
    } else if (match.type === "italic") {
      result.push(
        <em key={k++} className="italic text-zinc-300">
          {match.content}
        </em>
      );
    } else if (match.type === "code") {
      result.push(
        <code
          key={k++}
          className="px-1.5 py-0.5 rounded bg-zinc-800 text-zinc-300 text-sm font-mono"
        >
          {match.content}
        </code>
      );
    }
    pos = match.end;
  }

  if (pos < text.length) {
    result.push(<span key={k++}>{text.slice(pos)}</span>);
  }

  return result.length > 0 ? result : text;
}

export function MarkdownContent({ content }: { content: string }) {
  const elements = useMemo(() => parseMarkdown(content), [content]);
  return <div className="markdown-content">{elements}</div>;
}
