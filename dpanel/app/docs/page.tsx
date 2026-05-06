import Link from "next/link";
import { BookOpen, ArrowRight, LayoutGrid, Wrench } from "lucide-react";
import { apps } from "@/lib/apps";

export default function DocsHomePage() {
  const liveApps = apps.filter((a) => a.status === "live");
  const comingSoon = apps.filter((a) => a.status === "coming_soon");

  return (
    <div className="max-w-4xl mx-auto px-6 py-12">
      {/* Header */}
      <div className="mb-12">
        <h1 className="text-3xl font-bold text-zinc-100 mb-3">Documentation</h1>
        <p className="text-zinc-400 text-lg">
          Guides, references, and troubleshooting for the DClaw platform and all 65 apps.
        </p>
      </div>

      {/* Ecosystem docs */}
      <div className="mb-12">
        <h2 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider mb-4 flex items-center gap-2">
          <Wrench className="w-4 h-4" />
          Ecosystem
        </h2>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-4">
          <DocCard
            href="/docs/ecosystem/getting-started"
            title="Getting Started"
            description="Install the platform, deploy your first app, and configure your environment."
          />
          <DocCard
            href="/docs/ecosystem/architecture"
            title="Architecture"
            description="System overview, network topology, CRD reference, and security model."
          />
          <DocCard
            href="/docs/ecosystem/reference"
            title="Reference"
            description="Conventions, stack details, API specs, and configuration options."
          />
          <DocCard
            href="/docs/ecosystem/troubleshooting"
            title="Troubleshooting"
            description="Common issues, FAQ, and debugging guides for platform operators."
          />
          <DocCard
            href="/docs/ecosystem/releases"
            title="Releases"
            description="Changelog, roadmap, and version history for the DClaw platform."
          />
        </div>
      </div>

      {/* App docs */}
      <div className="mb-12">
        <h2 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider mb-4 flex items-center gap-2">
          <LayoutGrid className="w-4 h-4" />
          App Documentation
        </h2>

        {liveApps.length > 0 && (
          <div className="mb-8">
            <h3 className="text-xs font-medium text-zinc-500 uppercase tracking-wider mb-3">
              Live Apps
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              {liveApps.map((app) => (
                <AppDocCard key={app.id} app={app} />
              ))}
            </div>
          </div>
        )}

        {comingSoon.length > 0 && (
          <div>
            <h3 className="text-xs font-medium text-zinc-500 uppercase tracking-wider mb-3">
              Coming Soon
            </h3>
            <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-3">
              {comingSoon.map((app) => (
                <AppDocCard key={app.id} app={app} />
              ))}
            </div>
          </div>
        )}
      </div>

      {/* Footer */}
      <div className="text-xs text-zinc-600 text-center pt-8 border-t border-zinc-800">
        DClaw Platform Documentation — {apps.length} apps documented
      </div>
    </div>
  );
}

function DocCard({
  href,
  title,
  description,
}: {
  href: string;
  title: string;
  description: string;
}) {
  return (
    <Link
      href={href}
      className="group flex flex-col p-4 rounded-xl border border-zinc-800 bg-zinc-900/40 hover:bg-zinc-800/60 hover:border-zinc-700 transition-all"
    >
      <div className="flex items-center justify-between mb-2">
        <h3 className="font-semibold text-zinc-200 group-hover:text-zinc-100 transition-colors">
          {title}
        </h3>
        <ArrowRight className="w-4 h-4 text-zinc-600 group-hover:text-zinc-400 transition-colors" />
      </div>
      <p className="text-sm text-zinc-500 leading-relaxed">{description}</p>
    </Link>
  );
}

function AppDocCard({ app }: { app: (typeof apps)[0] }) {
  const Icon = app.icon;
  return (
    <Link
      href={`/docs/apps/${app.id}`}
      className="group flex items-center gap-3 p-3 rounded-xl border border-zinc-800 bg-zinc-900/40 hover:bg-zinc-800/60 hover:border-zinc-700 transition-all"
    >
      <div
        className="w-10 h-10 rounded-lg flex items-center justify-center shrink-0"
        style={{ backgroundColor: app.bgColor }}
      >
        <Icon className="w-5 h-5" style={{ color: app.color }} />
      </div>
      <div className="min-w-0">
        <h3 className="text-sm font-medium text-zinc-200 group-hover:text-zinc-100 truncate">
          {app.name}
        </h3>
        <p className="text-xs text-zinc-500 truncate">{app.tagline}</p>
      </div>
      <ArrowRight className="w-4 h-4 text-zinc-600 group-hover:text-zinc-400 ml-auto shrink-0" />
    </Link>
  );
}
