import fs from "fs";
import path from "path";

const DOCS_DIR = path.join(process.cwd(), "docs-content", "ecosystem");

export interface DocNode {
  path: string;
  title: string;
  children?: DocNode[];
}

function extractTitle(content: string): string {
  const match = content.match(/^#\s+(.+)$/m);
  return match ? match[1].trim() : "Untitled";
}

export function readEcosystemDoc(slugPath: string): string | null {
  const filePath = path.join(DOCS_DIR, slugPath + ".md");
  try {
    return fs.readFileSync(filePath, "utf-8");
  } catch {
    // Try index.md inside directory
    const indexPath = path.join(DOCS_DIR, slugPath, "index.md");
    try {
      return fs.readFileSync(indexPath, "utf-8");
    } catch {
      return null;
    }
  }
}

export function listEcosystemDocs(): DocNode[] {
  const sections: DocNode[] = [];

  try {
    const entries = fs.readdirSync(DOCS_DIR, { withFileTypes: true });

    for (const entry of entries) {
      if (entry.isDirectory()) {
        const sectionPath = entry.name;
        const indexPath = path.join(DOCS_DIR, sectionPath, "index.md");
        let title = sectionPath.replace(/-/g, " ");
        title = title.charAt(0).toUpperCase() + title.slice(1);

        try {
          const content = fs.readFileSync(indexPath, "utf-8");
          title = extractTitle(content);
        } catch {
          // use default title
        }

        const children: DocNode[] = [];
        const childEntries = fs.readdirSync(path.join(DOCS_DIR, sectionPath), {
          withFileTypes: true,
        });

        for (const child of childEntries) {
          if (child.isFile() && child.name.endsWith(".md") && child.name !== "index.md") {
            const childPath = `${sectionPath}/${child.name.replace(".md", "")}`;
            const childFilePath = path.join(DOCS_DIR, sectionPath, child.name);
            let childTitle = child.name.replace(".md", "").replace(/-/g, " ");
            childTitle = childTitle.charAt(0).toUpperCase() + childTitle.slice(1);

            try {
              const childContent = fs.readFileSync(childFilePath, "utf-8");
              childTitle = extractTitle(childContent);
            } catch {
              // use default
            }

            children.push({ path: childPath, title: childTitle });
          }
        }

        sections.push({
          path: sectionPath,
          title,
          children: children.length > 0 ? children : undefined,
        });
      }
    }
  } catch {
    // docs-content/ecosystem may not exist yet
  }

  // Sort in a logical order
  const order = ["getting-started", "architecture", "reference", "troubleshooting", "releases"];
  sections.sort((a, b) => {
    const aIndex = order.indexOf(a.path);
    const bIndex = order.indexOf(b.path);
    if (aIndex !== -1 && bIndex !== -1) return aIndex - bIndex;
    if (aIndex !== -1) return -1;
    if (bIndex !== -1) return 1;
    return a.path.localeCompare(b.path);
  });

  return sections;
}
