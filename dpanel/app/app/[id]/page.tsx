"use client";

import { useState, useEffect } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import {
  ArrowLeft,
  Download,
  Trash2,
  ExternalLink,
  Check,
  Sparkles,
  Tag,
  GitBranch,
  Loader2,
} from "lucide-react";
import { type App } from "@/lib/apps";
import { fetchAppById } from "@/lib/api";
import { useInstall } from "@/components/install-context";

export default function AppDetailPage() {
  const params = useParams();
  const { install, uninstall, isInstalled } = useInstall();

  const appId = params.id as string;
  const [app, setApp] = useState<App | null>(null);
  const [loading, setLoading] = useState(true);
  const [notFound, setNotFound] = useState(false);

  useEffect(() => {
    let mounted = true;
    fetchAppById(appId).then((data) => {
      if (!mounted) return;
      if (data) {
        setApp(data);
      } else {
        setNotFound(true);
      }
      setLoading(false);
    });
    return () => {
      mounted = false;
    };
  }, [appId]);

  if (loading) {
    return (
      <main className="flex flex-1 flex-col items-center justify-center py-12 px-6">
        <Loader2 className="w-8 h-8 text-zinc-500 animate-spin" />
        <p className="text-zinc-500 mt-3 text-sm">Loading app details...</p>
      </main>
    );
  }

  if (notFound || !app) {
    return (
      <main className="flex flex-1 flex-col items-center justify-center py-12 px-6">
        <div className="text-center">
          <h1 className="text-2xl font-bold text-zinc-100 mb-2">App not found</h1>
          <p className="text-zinc-500 mb-6">
            The app you are looking for does not exist.
          </p>
          <Link
            href="/"
            className="inline-flex items-center gap-2 px-4 py-2 rounded-lg bg-zinc-800 text-zinc-100 hover:bg-zinc-700 transition-colors"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to DPanel
          </Link>
        </div>
      </main>
    );
  }

  const Icon = app.icon;
  const isLive = app.status === "live";
  const installed = isInstalled(app.id);

  return (
    <main className="flex flex-1 flex-col items-center py-12 px-6">
      <div className="w-full max-w-2xl">
        {/* Back link */}
        <Link
          href="/"
          className="inline-flex items-center gap-2 text-sm text-zinc-500 hover:text-zinc-300 transition-colors mb-8"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to DPanel
        </Link>

        {/* Hero */}
        <div className="flex flex-col sm:flex-row items-start sm:items-center gap-6 mb-10">
          <div
            className="w-20 h-20 rounded-2xl flex items-center justify-center shrink-0"
            style={{ backgroundColor: app.bgColor }}
          >
            {Icon && <Icon className="w-10 h-10" style={{ color: app.color }} />}
          </div>

          <div className="flex-1 min-w-0">
            <div className="flex items-center gap-3 mb-1 flex-wrap">
              <h1 className="text-2xl font-bold text-zinc-100">{app.name}</h1>
              <span
                className={`text-[10px] font-medium px-2.5 py-0.5 rounded-full ${
                  isLive
                    ? "bg-emerald-500/15 text-emerald-400"
                    : "bg-zinc-700/40 text-zinc-500"
                }`}
              >
                {isLive ? "Live" : "Coming Soon"}
              </span>
              {installed && (
                <span className="text-[10px] font-medium px-2.5 py-0.5 rounded-full bg-blue-500/15 text-blue-400 flex items-center gap-1">
                  <Download className="w-3 h-3" />
                  Installed
                </span>
              )}
            </div>
            <p className="text-zinc-400">{app.tagline}</p>
          </div>
        </div>

        {/* Actions */}
        <div className="flex flex-wrap gap-3 mb-10">
          {isLive ? (
            installed ? (
              <>
                <a
                  href={`https://${app.domain}`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-emerald-600 text-white font-medium hover:bg-emerald-500 transition-colors"
                >
                  <ExternalLink className="w-4 h-4" />
                  Open App
                </a>
                <button
                  onClick={() => uninstall(app.id)}
                  className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-zinc-800 text-zinc-300 font-medium hover:bg-red-900/30 hover:text-red-400 transition-colors"
                >
                  <Trash2 className="w-4 h-4" />
                  Uninstall
                </button>
              </>
            ) : (
              <button
                onClick={() => install(app.id)}
                className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-zinc-100 text-zinc-950 font-medium hover:bg-white transition-colors"
              >
                <Download className="w-4 h-4" />
                Install
              </button>
            )
          ) : (
            <button
              disabled
              className="inline-flex items-center gap-2 px-5 py-2.5 rounded-xl bg-zinc-800/50 text-zinc-600 font-medium cursor-not-allowed"
            >
              <Sparkles className="w-4 h-4" />
              Coming Soon
            </button>
          )}
        </div>

        {/* Description */}
        <div className="mb-10">
          <h2 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider mb-3">
            About
          </h2>
          <p className="text-zinc-400 leading-relaxed">{app.description}</p>
        </div>

        {/* Features */}
        <div className="mb-10">
          <h2 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider mb-3">
            Features
          </h2>
          <ul className="space-y-2.5">
            {app.features.map((feature, i) => (
              <li
                key={i}
                className="flex items-start gap-3 text-zinc-400 text-sm"
              >
                <span className="mt-0.5 shrink-0 w-4 h-4 rounded-full bg-emerald-500/15 flex items-center justify-center">
                  <Check className="w-3 h-3 text-emerald-400" />
                </span>
                {feature}
              </li>
            ))}
          </ul>
        </div>

        {/* Metadata */}
        <div className="flex flex-wrap gap-4">
          <div className="flex items-center gap-2 text-xs text-zinc-500 bg-zinc-900/60 px-3 py-1.5 rounded-lg border border-zinc-800/50">
            <Tag className="w-3.5 h-3.5" />
            {app.category}
          </div>
          <div className="flex items-center gap-2 text-xs text-zinc-500 bg-zinc-900/60 px-3 py-1.5 rounded-lg border border-zinc-800/50">
            <GitBranch className="w-3.5 h-3.5" />
            v{app.version}
          </div>
        </div>
      </div>
    </main>
  );
}
