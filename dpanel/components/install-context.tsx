"use client";

import {
  createContext,
  useContext,
  useState,
  useEffect,
  type ReactNode,
} from "react";

interface InstallState {
  installed: Set<string>;
  install: (appId: string) => void;
  uninstall: (appId: string) => void;
  isInstalled: (appId: string) => boolean;
}

const InstallContext = createContext<InstallState | null>(null);

const STORAGE_KEY = "dpanel-installed-apps";

export function InstallProvider({ children }: { children: ReactNode }) {
  const [installed, setInstalled] = useState<Set<string>>(new Set());
  const [hydrated, setHydrated] = useState(false);

  useEffect(() => {
    try {
      const raw = localStorage.getItem(STORAGE_KEY);
      if (raw) {
        const parsed = JSON.parse(raw) as string[];
        setInstalled(new Set(parsed));
      }
    } catch {
      // ignore corrupt storage
    }
    setHydrated(true);
  }, []);

  useEffect(() => {
    if (hydrated) {
      localStorage.setItem(STORAGE_KEY, JSON.stringify(Array.from(installed)));
    }
  }, [installed, hydrated]);

  const install = (appId: string) => {
    setInstalled((prev) => {
      const next = new Set(prev);
      next.add(appId);
      return next;
    });
  };

  const uninstall = (appId: string) => {
    setInstalled((prev) => {
      const next = new Set(prev);
      next.delete(appId);
      return next;
    });
  };

  const isInstalled = (appId: string) => installed.has(appId);

  return (
    <InstallContext.Provider value={{ installed, install, uninstall, isInstalled }}>
      {children}
    </InstallContext.Provider>
  );
}

export function useInstall() {
  const ctx = useContext(InstallContext);
  if (!ctx) {
    throw new Error("useInstall must be used within an InstallProvider");
  }
  return ctx;
}
