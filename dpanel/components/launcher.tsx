"use client";

import { useState } from "react";
import { Search, ExternalLink, Sparkles } from "lucide-react";
import { apps } from "@/lib/apps";

export function Launcher() {
  const [query, setQuery] = useState("");

  const filtered = apps.filter(
    (app) =>
      app.name.toLowerCase().includes(query.toLowerCase()) ||
      app.tagline.toLowerCase().includes(query.toLowerCase()) ||
      app.category.toLowerCase().includes(query.toLowerCase())
  );

  return (
    <div className="flex flex-col items-center w-full max-w-3xl mx-auto px-6 py-12">
      {/* Header */}
      <div className="flex items-center gap-3 mb-2">
        <Sparkles className="w-8 h-8 text-zinc-100" />
        <h1 className="text-3xl font-bold tracking-tight text-zinc-100">
          DPanel
        </h1>
      </div>
      <p className="text-zinc-400 mb-10 text-center">
        Your AI app store and command center
      </p>

      {/* Search */}
      <div className="relative w-full max-w-md mb-12">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-zinc-500" />
        <input
          type="text"
          placeholder="Search apps..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          className="w-full pl-10 pr-4 py-2.5 rounded-xl bg-zinc-900 border border-zinc-800 text-zinc-100 placeholder:text-zinc-600 focus:outline-none focus:ring-2 focus:ring-zinc-700 transition-all"
        />
      </div>

      {/* App Grid */}
      <div className="grid grid-cols-2 sm:grid-cols-3 gap-4 w-full">
        {filtered.map((app) => {
          const Icon = app.icon;
          const isLive = app.status === "live";
          return (
            <a
              key={app.id}
              href={isLive ? `https://${app.domain}` : undefined}
              target={isLive ? "_blank" : undefined}
              rel={isLive ? "noopener noreferrer" : undefined}
              className={`
                group relative flex flex-col items-center justify-center
                rounded-2xl p-6 border transition-all duration-200
                ${
                  isLive
                    ? "border-zinc-800 bg-zinc-900/50 hover:bg-zinc-800/60 hover:scale-[1.02] hover:shadow-lg cursor-pointer"
                    : "border-zinc-800/40 bg-zinc-900/20 opacity-60 cursor-not-allowed"
                }
              `}
            >
              {/* Status badge */}
              <span
                className={`absolute top-3 right-3 text-[10px] font-medium px-2 py-0.5 rounded-full ${
                  isLive
                    ? "bg-emerald-500/15 text-emerald-400"
                    : "bg-zinc-700/40 text-zinc-500"
                }`}
              >
                {isLive ? "Live" : "Soon"}
              </span>

              {/* Icon */}
              <div
                className="w-14 h-14 rounded-xl flex items-center justify-center mb-4 transition-transform group-hover:scale-110"
                style={{ backgroundColor: app.bgColor }}
              >
                <Icon
                  className="w-7 h-7"
                  style={{ color: app.color }}
                />
              </div>

              {/* Name */}
              <span className="text-sm font-semibold text-zinc-100 text-center">
                {app.name}
              </span>

              {/* Tagline */}
              <span className="text-xs text-zinc-500 text-center mt-1 line-clamp-1">
                {app.tagline}
              </span>

              {/* External link indicator */}
              {isLive && (
                <ExternalLink className="absolute bottom-3 right-3 w-3.5 h-3.5 text-zinc-600 opacity-0 group-hover:opacity-100 transition-opacity" />
              )}
            </a>
          );
        })}
      </div>

      {/* Empty state */}
      {filtered.length === 0 && (
        <div className="text-zinc-500 text-sm mt-8">
          No apps match &ldquo;{query}&rdquo;
        </div>
      )}

      {/* Footer */}
      <div className="mt-12 text-xs text-zinc-600">
        DClaw Platform v0.1.0 — {apps.filter((a) => a.status === "live").length}{" "}
        of {apps.length} apps live
      </div>
    </div>
  );
}
