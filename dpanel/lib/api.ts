import { apps as staticApps, type App } from "@/lib/apps";

const API_BASE = process.env.NEXT_PUBLIC_DPANEL_API_URL || "";

/** Known DClaw app repositories that publish manifests */
const REPO_APPS = [
  "flow",
  "rag",
  "agent",
  "med",
  "code",
  "chat",
  "learn",
];

export interface Manifest {
  app_id: string;
  name: string;
  description?: string;
  tagline?: string;
  version: string;
  category: string;
  icon?: string;
  color: string;
  entrypoint: string;
  api_base: string;
  permissions: string[];
  requires_auth: boolean;
  billing_plan: string;
}

export interface RegistryApp {
  appId: string;
  appName: string;
  appIcon: string;
  category: string;
  version: string;
  primaryColor: string;
  path: string;
  tier: string;
  status: string;
  url: string;
}

function manifestToApp(manifest: Manifest, staticApp?: App): App {
  return {
    id: manifest.app_id,
    name: manifest.name || staticApp?.name || manifest.app_id,
    tagline: manifest.tagline || staticApp?.tagline || "",
    description: manifest.description || staticApp?.description || "",
    features: staticApp?.features || [],
    version: manifest.version || staticApp?.version || "0.1.0",
    icon: staticApp?.icon || (() => null) as unknown as App["icon"],
    color: manifest.color || staticApp?.color || "#3B82F6",
    bgColor: staticApp?.bgColor || `rgba(${hexToRgb(manifest.color || "#3B82F6")}, 0.15)`,
    domain: manifest.api_base ? manifest.api_base.replace(/^https?:\/\//, "") : staticApp?.domain || "",
    category: manifest.category || staticApp?.category || "Platform",
    status: "live",
  };
}

function hexToRgb(hex: string): string {
  const clean = hex.replace("#", "");
  const bigint = parseInt(clean, 16);
  const r = (bigint >> 16) & 255;
  const g = (bigint >> 8) & 255;
  const b = bigint & 255;
  return `${r}, ${g}, ${b}`;
}

async function fetchManifest(appId: string): Promise<Manifest | null> {
  const url = `https://raw.githubusercontent.com/dclawstack/dclaw-${appId}/main/frontend/public/dclaw-manifest.json`;
  try {
    const res = await fetch(url, { next: { revalidate: 60 } });
    if (!res.ok) return null;
    return (await res.json()) as Manifest;
  } catch {
    return null;
  }
}

function mapRegistryToApp(reg: RegistryApp): App {
  const staticApp = staticApps.find((a) => a.id === reg.appId);

  return {
    id: reg.appId,
    name: reg.appName || staticApp?.name || reg.appId,
    tagline: staticApp?.tagline || "",
    description: staticApp?.description || "",
    features: staticApp?.features || [],
    version: reg.version || staticApp?.version || "0.1.0",
    icon: staticApp?.icon || (() => null) as unknown as App["icon"],
    color: reg.primaryColor || staticApp?.color || "#3B82F6",
    bgColor: staticApp?.bgColor || "rgba(59, 130, 246, 0.15)",
    domain: reg.url ? reg.url.replace(/^https?:\/\//, "") : staticApp?.domain || "",
    category: reg.category || staticApp?.category || "Platform",
    status:
      reg.status === "ready"
        ? "live"
        : (staticApp?.status ?? "coming_soon"),
  };
}

/** Fetch apps by trying (in order):
 *  1. dpanel-api registry
 *  2. GitHub repo manifests
 *  3. Static fallback
 */
export async function fetchApps(): Promise<App[]> {
  // Try dpanel-api first
  if (API_BASE) {
    try {
      const res = await fetch(`${API_BASE}/api/v1/apps`, {
        headers: { Accept: "application/json" },
        cache: "no-store",
      });
      if (res.ok) {
        const data = await res.json();
        const registryApps: RegistryApp[] = data.apps || [];
        if (registryApps.length > 0) {
          const registryIds = new Set(registryApps.map((r) => r.appId));
          const merged = registryApps.map(mapRegistryToApp);
          for (const staticApp of staticApps) {
            if (!registryIds.has(staticApp.id)) {
              merged.push(staticApp);
            }
          }
          return merged;
        }
      }
    } catch {
      // fall through
    }
  }

  // Try GitHub manifests
  const manifestPromises = REPO_APPS.map(async (id) => {
    const manifest = await fetchManifest(id);
    if (!manifest) return null;
    const staticApp = staticApps.find((a) => a.id === id);
    return manifestToApp(manifest, staticApp);
  });

  const manifestResults = await Promise.all(manifestPromises);
  const liveApps = manifestResults.filter((a): a is App => a !== null);

  if (liveApps.length > 0) {
    const liveIds = new Set(liveApps.map((a) => a.id));
    // Add static apps that don't have manifests yet
    for (const staticApp of staticApps) {
      if (!liveIds.has(staticApp.id)) {
        liveApps.push(staticApp);
      }
    }
    // Sort to match static order
    const order = new Map(staticApps.map((a, i) => [a.id, i]));
    liveApps.sort((a, b) => (order.get(a.id) ?? 999) - (order.get(b.id) ?? 999));
    return liveApps;
  }

  // Fallback to static
  return staticApps;
}

export async function fetchAppById(id: string): Promise<App | null> {
  // Try dpanel-api first
  if (API_BASE) {
    try {
      const res = await fetch(`${API_BASE}/api/v1/apps/${id}`, {
        headers: { Accept: "application/json" },
        cache: "no-store",
      });
      if (res.ok) {
        const data = await res.json();
        return mapRegistryToApp(data.app as RegistryApp);
      }
      if (res.status === 404) {
        return staticApps.find((a) => a.id === id) || null;
      }
    } catch {
      // fall through
    }
  }

  // Try GitHub manifest
  const manifest = await fetchManifest(id);
  if (manifest) {
    const staticApp = staticApps.find((a) => a.id === id);
    return manifestToApp(manifest, staticApp);
  }

  return staticApps.find((a) => a.id === id) || null;
}
