import { apps as staticApps, type App } from "@/lib/apps";

const API_BASE = process.env.NEXT_PUBLIC_DPANEL_API_URL || "";

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

function mapRegistryToApp(reg: RegistryApp): App {
  // Map registry app to frontend App shape
  // We keep static metadata (description, features, icon) and overlay registry data
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

export async function fetchApps(): Promise<App[]> {
  if (!API_BASE) {
    return staticApps;
  }

  try {
    const res = await fetch(`${API_BASE}/api/v1/apps`, {
      headers: { Accept: "application/json" },
      // Avoid caching so we get fresh registry data
      cache: "no-store",
    });

    if (!res.ok) {
      console.warn("dpanel-api returned non-OK, falling back to static", res.status);
      return staticApps;
    }

    const data = await res.json();
    const registryApps: RegistryApp[] = data.apps || [];

    if (registryApps.length === 0) {
      return staticApps;
    }

    // Build merged list: registry apps + any static apps not in registry
    const registryIds = new Set(registryApps.map((r) => r.appId));
    const merged = registryApps.map(mapRegistryToApp);

    // Append static-only apps (not yet registered by operator)
    for (const staticApp of staticApps) {
      if (!registryIds.has(staticApp.id)) {
        merged.push(staticApp);
      }
    }

    return merged;
  } catch (err) {
    console.warn("dpanel-api unreachable, falling back to static", err);
    return staticApps;
  }
}

export async function fetchAppById(id: string): Promise<App | null> {
  if (!API_BASE) {
    return staticApps.find((a) => a.id === id) || null;
  }

  try {
    const res = await fetch(`${API_BASE}/api/v1/apps/${id}`, {
      headers: { Accept: "application/json" },
      cache: "no-store",
    });

    if (!res.ok) {
      if (res.status === 404) {
        return staticApps.find((a) => a.id === id) || null;
      }
      console.warn("dpanel-api returned non-OK, falling back to static", res.status);
      return staticApps.find((a) => a.id === id) || null;
    }

    const data = await res.json();
    const reg: RegistryApp = data.app;
    return mapRegistryToApp(reg);
  } catch (err) {
    console.warn("dpanel-api unreachable, falling back to static", err);
    return staticApps.find((a) => a.id === id) || null;
  }
}
