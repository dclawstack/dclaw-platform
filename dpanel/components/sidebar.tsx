"use client";

import Link from "next/link";
import { usePathname } from "next/navigation";
import {
  Sparkles,
  Home,
  BookOpen,
  ChevronDown,
  ChevronRight,
  Search,
  LayoutGrid,
  Wrench,
  ArrowLeft,
} from "lucide-react";
import { useState } from "react";
import { apps } from "@/lib/apps";

interface SidebarProps {
  onNavigate?: () => void;
}

function groupAppsByCategory() {
  const groups: Record<string, typeof apps> = {};
  for (const app of apps) {
    const cat = app.category;
    if (!groups[cat]) groups[cat] = [];
    groups[cat].push(app);
  }
  return groups;
}

function NavItem({
  href,
  icon: Icon,
  label,
  active,
  onClick,
}: {
  href: string;
  icon: React.ElementType;
  label: string;
  active?: boolean;
  onClick?: () => void;
}) {
  return (
    <Link
      href={href}
      onClick={onClick}
      className={`
        flex items-center gap-3 px-3 py-2 rounded-lg text-sm transition-colors
        ${active
          ? "bg-zinc-800 text-zinc-100 font-medium"
          : "text-zinc-400 hover:text-zinc-100 hover:bg-zinc-800/50"
        }
      `}
    >
      <Icon className="w-4 h-4 shrink-0" />
      <span className="truncate">{label}</span>
    </Link>
  );
}

export function Sidebar({ onNavigate }: SidebarProps) {
  const pathname = usePathname();
  const [appsExpanded, setAppsExpanded] = useState(true);
  const [ecosystemExpanded, setEcosystemExpanded] = useState(true);
  const grouped = groupAppsByCategory();

  const isDocsRoute = pathname.startsWith("/docs");
  const isAppDocs = pathname.startsWith("/docs/apps/");
  const isEcosystemDocs = pathname.startsWith("/docs/ecosystem/");

  return (
    <div className="flex flex-col h-full">
      {/* Header */}
      <div className="flex items-center gap-3 px-4 py-4 border-b border-zinc-800">
        <div className="w-8 h-8 rounded-lg bg-zinc-800 flex items-center justify-center">
          <Sparkles className="w-4 h-4 text-zinc-100" />
        </div>
        <div>
          <h1 className="text-sm font-bold text-zinc-100">DPanel</h1>
          <p className="text-[10px] text-zinc-500">DClaw Platform</p>
        </div>
      </div>

      {/* Scrollable nav */}
      <nav className="flex-1 overflow-y-auto px-3 py-4 space-y-6">
        {/* Main */}
        <div className="space-y-1">
          <NavItem
            href="/"
            icon={Home}
            label="Home"
            active={pathname === "/"}
            onClick={onNavigate}
          />
          <NavItem
            href="/docs"
            icon={BookOpen}
            label="Documentation"
            active={isDocsRoute && !isAppDocs && !isEcosystemDocs}
            onClick={onNavigate}
          />
          <NavItem
            href="/docs/search"
            icon={Search}
            label="Search"
            active={pathname === "/docs/search"}
            onClick={onNavigate}
          />
        </div>

        {/* Apps */}
        <div>
          <button
            onClick={() => setAppsExpanded(!appsExpanded)}
            className="flex items-center gap-2 px-3 py-2 text-xs font-semibold text-zinc-500 uppercase tracking-wider w-full hover:text-zinc-300 transition-colors"
          >
            <LayoutGrid className="w-3.5 h-3.5" />
            <span className="flex-1 text-left">Apps</span>
            {appsExpanded ? (
              <ChevronDown className="w-3.5 h-3.5" />
            ) : (
              <ChevronRight className="w-3.5 h-3.5" />
            )}
          </button>

          {appsExpanded && (
            <div className="mt-1 space-y-4">
              {Object.entries(grouped).map(([category, categoryApps]) => (
                <div key={category}>
                  <div className="px-3 py-1 text-[10px] font-medium text-zinc-600 uppercase tracking-wider">
                    {category}
                  </div>
                  <div className="space-y-0.5 mt-1">
                    {categoryApps.map((app) => {
                      const appDocsPath = `/docs/apps/${app.id}`;
                      const isActive = pathname === appDocsPath || pathname.startsWith(`/docs/apps/${app.id}/`);
                      const Icon = app.icon;
                      return (
                        <Link
                          key={app.id}
                          href={appDocsPath}
                          onClick={onNavigate}
                          className={`
                            flex items-center gap-2.5 px-3 py-1.5 rounded-md text-xs transition-colors
                            ${isActive
                              ? "bg-zinc-800 text-zinc-100 font-medium"
                              : "text-zinc-500 hover:text-zinc-300 hover:bg-zinc-800/30"
                            }
                          `}
                        >
                          <Icon className="w-3.5 h-3.5 shrink-0" style={{ color: app.color }} />
                          <span className="truncate">{app.name.replace("DClaw ", "")}</span>
                          {app.status === "live" && (
                            <span className="ml-auto w-1.5 h-1.5 rounded-full bg-emerald-500/60" />
                          )}
                        </Link>
                      );
                    })}
                  </div>
                </div>
              ))}
            </div>
          )}
        </div>

        {/* Ecosystem */}
        <div>
          <button
            onClick={() => setEcosystemExpanded(!ecosystemExpanded)}
            className="flex items-center gap-2 px-3 py-2 text-xs font-semibold text-zinc-500 uppercase tracking-wider w-full hover:text-zinc-300 transition-colors"
          >
            <Wrench className="w-3.5 h-3.5" />
            <span className="flex-1 text-left">Ecosystem</span>
            {ecosystemExpanded ? (
              <ChevronDown className="w-3.5 h-3.5" />
            ) : (
              <ChevronRight className="w-3.5 h-3.5" />
            )}
          </button>

          {ecosystemExpanded && (
            <div className="mt-1 space-y-0.5">
              <NavItem
                href="/docs/ecosystem/getting-started"
                icon={ArrowLeft}
                label="Getting Started"
                active={isEcosystemDocs && pathname.includes("getting-started")}
                onClick={onNavigate}
              />
              <NavItem
                href="/docs/ecosystem/architecture"
                icon={ArrowLeft}
                label="Architecture"
                active={isEcosystemDocs && pathname.includes("architecture")}
                onClick={onNavigate}
              />
              <NavItem
                href="/docs/ecosystem/reference"
                icon={ArrowLeft}
                label="Reference"
                active={isEcosystemDocs && pathname.includes("reference")}
                onClick={onNavigate}
              />
              <NavItem
                href="/docs/ecosystem/troubleshooting"
                icon={ArrowLeft}
                label="Troubleshooting"
                active={isEcosystemDocs && pathname.includes("troubleshooting")}
                onClick={onNavigate}
              />
              <NavItem
                href="/docs/ecosystem/releases"
                icon={ArrowLeft}
                label="Releases"
                active={isEcosystemDocs && pathname.includes("releases")}
                onClick={onNavigate}
              />
            </div>
          )}
        </div>
      </nav>

      {/* Footer */}
      <div className="px-4 py-3 border-t border-zinc-800 text-[10px] text-zinc-600 text-center">
        DClaw Platform © 2026
      </div>
    </div>
  );
}
