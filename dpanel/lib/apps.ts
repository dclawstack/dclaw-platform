import {
  MessageSquare,
  Workflow,
  HeartPulse,
  BookOpen,
  TrendingUp,
  Palette,
  Code2,
  Bot,
  Search,
  type LucideIcon,
} from "lucide-react";

export interface App {
  id: string;
  name: string;
  tagline: string;
  icon: LucideIcon;
  color: string;
  bgColor: string;
  domain: string;
  category: string;
  status: "live" | "beta" | "coming_soon";
}

export const apps: App[] = [
  {
    id: "chat",
    name: "DClaw Chat",
    tagline: "AI conversations that remember",
    icon: MessageSquare,
    color: "#3B82F6",
    bgColor: "rgba(59, 130, 246, 0.15)",
    domain: "chat.dclawstack.io",
    category: "Communication",
    status: "live",
  },
  {
    id: "flow",
    name: "DClaw Flow",
    tagline: "Connect anything, automate everything",
    icon: Workflow,
    color: "#10B981",
    bgColor: "rgba(16, 185, 129, 0.15)",
    domain: "flow.dclawstack.io",
    category: "Automation",
    status: "coming_soon",
  },
  {
    id: "med",
    name: "DClaw Med",
    tagline: "Clinical intelligence at your fingertips",
    icon: HeartPulse,
    color: "#EF4444",
    bgColor: "rgba(239, 68, 68, 0.15)",
    domain: "med.dclawstack.io",
    category: "Healthcare",
    status: "coming_soon",
  },
  {
    id: "learn",
    name: "DClaw Learn",
    tagline: "Adaptive learning that works",
    icon: BookOpen,
    color: "#3B82F6",
    bgColor: "rgba(59, 130, 246, 0.15)",
    domain: "learn.dclawstack.io",
    category: "Education",
    status: "coming_soon",
  },
  {
    id: "seo",
    name: "DClaw SEO",
    tagline: "Rank higher with AI",
    icon: TrendingUp,
    color: "#10B981",
    bgColor: "rgba(16, 185, 129, 0.15)",
    domain: "seo.dclawstack.io",
    category: "Marketing",
    status: "coming_soon",
  },
  {
    id: "create",
    name: "DClaw Create",
    tagline: "Generate anything",
    icon: Palette,
    color: "#EC4899",
    bgColor: "rgba(236, 72, 153, 0.15)",
    domain: "create.dclawstack.io",
    category: "Media",
    status: "coming_soon",
  },
  {
    id: "code",
    name: "DClaw Code",
    tagline: "AI-native IDE inside your desktop",
    icon: Code2,
    color: "#1F2937",
    bgColor: "rgba(31, 41, 55, 0.15)",
    domain: "code.dclawstack.io",
    category: "Development",
    status: "coming_soon",
  },
  {
    id: "agent",
    name: "DClaw Agent",
    tagline: "Build, share, and sell AI agents",
    icon: Bot,
    color: "#8B5CF6",
    bgColor: "rgba(139, 92, 246, 0.15)",
    domain: "agent.dclawstack.io",
    category: "Platform",
    status: "coming_soon",
  },
  {
    id: "rag",
    name: "DClaw RAG",
    tagline: "Universal knowledge retrieval",
    icon: Search,
    color: "#F59E0B",
    bgColor: "rgba(245, 158, 11, 0.15)",
    domain: "rag.dclawstack.io",
    category: "Platform",
    status: "coming_soon",
  },
];
