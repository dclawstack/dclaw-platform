import Link from "next/link";
import { notFound } from "next/navigation";
import { ArrowLeft, BookOpen, Wrench, FileText, Tag, GitBranch } from "lucide-react";
import { apps } from "@/lib/apps";

interface AppDocsPageProps {
  params: Promise<{ id: string }>;
}

export default async function AppDocsPage({ params }: AppDocsPageProps) {
  const { id } = await params;
  const app = apps.find((a) => a.id === id);

  if (!app) {
    notFound();
  }

  const Icon = app.icon;

  const docSections = [
    {
      group: "Getting Started",
      pages: [
        { title: "Overview", path: `getting-started/index` },
        { title: "Installation", path: `getting-started/installation` },
        { title: "Quickstart", path: `getting-started/quickstart` },
        { title: "Configuration", path: `getting-started/configuration` },
      ],
    },
    {
      group: "Guides",
      pages: [
        { title: "Overview", path: `guides/index` },
        { title: "Use Cases", path: `guides/use-cases` },
        { title: "Best Practices", path: `guides/best-practices` },
      ],
    },
    {
      group: "Reference",
      pages: [
        { title: "Overview", path: `reference/index` },
        { title: "Architecture", path: `reference/architecture` },
        { title: "Stack", path: `reference/stack` },
        { title: "API", path: `reference/api` },
      ],
    },
    {
      group: "Troubleshooting",
      pages: [
        { title: "Overview", path: `troubleshooting/index` },
        { title: "Common Issues", path: `troubleshooting/common-issues` },
        { title: "FAQ", path: `troubleshooting/faq` },
      ],
    },
    {
      group: "Releases",
      pages: [
        { title: "Overview", path: `releases/index` },
        { title: "Changelog", path: `releases/changelog` },
        { title: "Roadmap", path: `releases/roadmap` },
      ],
    },
  ];

  return (
    <div className="max-w-4xl mx-auto px-6 py-12">
      {/* Back + Hero */}
      <div className="mb-10">
        <Link
          href="/docs"
          className="inline-flex items-center gap-2 text-sm text-zinc-500 hover:text-zinc-300 transition-colors mb-6"
        >
          <ArrowLeft className="w-4 h-4" />
          Back to Docs
        </Link>

        <div className="flex items-center gap-5">
          <div
            className="w-16 h-16 rounded-2xl flex items-center justify-center shrink-0"
            style={{ backgroundColor: app.bgColor }}
          >
            <Icon className="w-8 h-8" style={{ color: app.color }} />
          </div>
          <div>
            <h1 className="text-2xl font-bold text-zinc-100">{app.name}</h1>
            <p className="text-zinc-400">{app.tagline}</p>
          </div>
        </div>
      </div>

      {/* Metadata */}
      <div className="flex flex-wrap gap-3 mb-10">
        <span className="flex items-center gap-2 text-xs text-zinc-500 bg-zinc-900/60 px-3 py-1.5 rounded-lg border border-zinc-800/50">
          <Tag className="w-3.5 h-3.5" />
          {app.category}
        </span>
        <span className="flex items-center gap-2 text-xs text-zinc-500 bg-zinc-900/60 px-3 py-1.5 rounded-lg border border-zinc-800/50">
          <GitBranch className="w-3.5 h-3.5" />
          v{app.version}
        </span>
        <span
          className={`text-xs font-medium px-3 py-1.5 rounded-lg border ${
            app.status === "live"
              ? "bg-emerald-500/10 text-emerald-400 border-emerald-500/20"
              : "bg-zinc-800/40 text-zinc-500 border-zinc-700/40"
          }`}
        >
          {app.status === "live" ? "Live" : "Coming Soon"}
        </span>
      </div>

      {/* Description */}
      <div className="mb-10">
        <p className="text-zinc-400 leading-relaxed">{app.description}</p>
      </div>

      {/* Doc sections */}
      <div className="space-y-8">
        {docSections.map((section) => (
          <div key={section.group}>
            <h2 className="text-sm font-semibold text-zinc-300 uppercase tracking-wider mb-3 flex items-center gap-2">
              {section.group === "Getting Started" && <BookOpen className="w-4 h-4" />}
              {section.group === "Guides" && <FileText className="w-4 h-4" />}
              {section.group === "Reference" && <Tag className="w-4 h-4" />}
              {section.group === "Troubleshooting" && <Wrench className="w-4 h-4" />}
              {section.group === "Releases" && <GitBranch className="w-4 h-4" />}
              {section.group}
            </h2>
            <div className="grid grid-cols-1 sm:grid-cols-2 gap-2">
              {section.pages.map((page) => (
                <Link
                  key={page.path}
                  href={`https://github.com/dclawstack/dclaw-${app.id}/tree/main/docs/${page.path}.md`}
                  target="_blank"
                  rel="noopener noreferrer"
                  className="flex items-center gap-2 px-4 py-3 rounded-lg border border-zinc-800 bg-zinc-900/30 hover:bg-zinc-800/50 hover:border-zinc-700 transition-all text-sm text-zinc-400 hover:text-zinc-200"
                >
                  <FileText className="w-4 h-4 text-zinc-600 shrink-0" />
                  <span>{page.title}</span>
                </Link>
              ))}
            </div>
          </div>
        ))}
      </div>

      {/* Placeholder notice */}
      <div className="mt-12 p-4 rounded-xl border border-dashed border-zinc-700 bg-zinc-900/20">
        <p className="text-sm text-zinc-500 text-center">
          Documentation content is being migrated from each app&apos;s repository.
          <br />
          Links above open the source markdown files on GitHub.
        </p>
      </div>
    </div>
  );
}
