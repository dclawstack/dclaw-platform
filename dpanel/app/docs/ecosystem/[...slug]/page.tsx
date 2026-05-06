import { notFound } from "next/navigation";
import Link from "next/link";
import { ArrowLeft } from "lucide-react";
import { readEcosystemDoc, listEcosystemDocs } from "@/lib/docs/content";
import { MarkdownContent } from "@/components/docs/markdown-content";
import { DocsToc } from "@/components/docs/docs-toc";

interface EcosystemDocPageProps {
  params: Promise<{ slug: string[] }>;
}

export default async function EcosystemDocPage({ params }: EcosystemDocPageProps) {
  const { slug } = await params;
  const path = slug.join("/");

  const content = readEcosystemDoc(path);
  if (!content) {
    notFound();
  }

  const docs = listEcosystemDocs();

  return (
    <div className="flex h-full">
      {/* Docs sidebar */}
      <aside className="hidden lg:block w-64 border-r border-zinc-800 bg-zinc-900/30 overflow-y-auto">
        <div className="px-4 py-6">
          <Link
            href="/docs"
            className="inline-flex items-center gap-2 text-sm text-zinc-500 hover:text-zinc-300 transition-colors mb-6"
          >
            <ArrowLeft className="w-4 h-4" />
            Back to Docs
          </Link>

          <nav className="space-y-6">
            {docs.map((section) => (
              <div key={section.path}>
                <Link
                  href={`/docs/ecosystem/${section.path}`}
                  className={`
                    block text-sm font-medium mb-2 transition-colors
                    ${path === section.path
                      ? "text-zinc-100"
                      : "text-zinc-400 hover:text-zinc-200"
                    }
                  `}
                >
                  {section.title}
                </Link>
                {section.children && section.children.length > 0 && (
                  <div className="ml-3 space-y-1 border-l border-zinc-800 pl-3">
                    {section.children.map((child) => (
                      <Link
                        key={child.path}
                        href={`/docs/ecosystem/${child.path}`}
                        className={`
                          block text-xs transition-colors
                          ${path === child.path
                            ? "text-zinc-100 font-medium"
                            : "text-zinc-500 hover:text-zinc-300"
                          }
                        `}
                      >
                        {child.title}
                      </Link>
                    ))}
                  </div>
                )}
              </div>
            ))}
          </nav>
        </div>
      </aside>

      {/* Content */}
      <div className="flex-1 overflow-y-auto">
        <article className="max-w-3xl mx-auto px-6 py-12">
          <MarkdownContent content={content} />
        </article>
      </div>

      {/* TOC sidebar */}
      <aside className="hidden xl:block w-56 border-l border-zinc-800 bg-zinc-900/30 overflow-y-auto">
        <div className="px-4 py-6">
          <DocsToc content={content} />
        </div>
      </aside>
    </div>
  );
}
