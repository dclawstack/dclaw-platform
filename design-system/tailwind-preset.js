/**
 * DKube Tailwind CSS Preset (v3)
 * Usage: import in tailwind.config.ts:
 *   import dkubePreset from "../dclaw-platform/design-system/tailwind-preset";
 *   presets: [dkubePreset],
 */

/** @type {import('tailwindcss').Config} */
const dkubePreset = {
  darkMode: ["class"],
  theme: {
    extend: {
      colors: {
        dk: {
          purple: {
            DEFAULT: "#6B53A3",
            light: "#9985BF",
            deep: "#4A3A7A",
            wash: "rgba(107, 83, 163, 0.08)",
          },
          surface: {
            DEFAULT: "#0E0E10",
            raised: "#1F1F23",
            elevated: "#2A2A30",
          },
          border: "rgba(255, 255, 255, 0.08)",
          body: "#F4F2F8",
          muted: {
            DEFAULT: "#9E9AAB",
            darker: "#6E6E76",
          },
          success: "#22C55E",
          warning: "#F59E0B",
          error: "#EF4444",
          info: "#3B82F6",
        },
        border: "hsl(var(--border))",
        input: "hsl(var(--input))",
        ring: "hsl(var(--ring))",
        background: "hsl(var(--background))",
        foreground: "hsl(var(--foreground))",
        primary: {
          DEFAULT: "hsl(var(--primary))",
          foreground: "hsl(var(--primary-foreground))",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary))",
          foreground: "hsl(var(--secondary-foreground))",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive))",
          foreground: "hsl(var(--destructive-foreground))",
        },
        muted: {
          DEFAULT: "hsl(var(--muted))",
          foreground: "hsl(var(--muted-foreground))",
        },
        accent: {
          DEFAULT: "hsl(var(--accent))",
          foreground: "hsl(var(--accent-foreground))",
        },
        popover: {
          DEFAULT: "hsl(var(--popover))",
          foreground: "hsl(var(--popover-foreground))",
        },
        card: {
          DEFAULT: "hsl(var(--card))",
          foreground: "hsl(var(--card-foreground))",
        },
      },
      fontFamily: {
        sans: ["var(--font-inter)", "Inter", "system-ui", "sans-serif"],
        display: ["var(--font-manrope)", "Manrope", "system-ui", "sans-serif"],
        mono: ["var(--font-jetbrains)", "JetBrains Mono", "ui-monospace", "monospace"],
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      boxShadow: {
        "dk-sm": "0 1px 2px rgba(0, 0, 0, 0.25)",
        "dk-md": "0 4px 12px rgba(0, 0, 0, 0.35)",
        "dk-lg": "0 8px 24px rgba(0, 0, 0, 0.45)",
      },
      transitionTimingFunction: {
        "dk-fast": "cubic-bezier(0.4, 0, 0.2, 1)",
      },
      transitionDuration: {
        "dk-fast": "150ms",
        "dk-base": "250ms",
        "dk-slow": "350ms",
      },
    },
  },
};

module.exports = dkubePreset;
