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
  description: string;
  features: string[];
  version: string;
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
    description:
      "A multi-model AI chat interface with persistent memory, voice input, and local or cloud LLM support. Switch between Ollama, OpenRouter, and proprietary models on the fly.",
    features: [
      "Multi-model chat (Ollama, OpenRouter, Claude, GPT-4)",
      "Persistent conversation history",
      "Voice input and speech-to-text",
      "Desktop app via Tauri",
      "End-to-end encrypted messaging",
    ],
    version: "0.2.0",
    icon: MessageSquare,
    color: "#3B82F6",
    bgColor: "rgba(59, 130, 246, 0.15)",
    domain: "chat.dclawstack.io",
    category: "Communication",
    status: "coming_soon",
  },
  {
    id: "flow",
    name: "DClaw Flow",
    tagline: "Connect anything, automate everything",
    description:
      "Visual workflow builder with AI-powered node generation. Connect APIs, databases, and services with drag-and-drop simplicity.",
    features: [
      "Visual node editor",
      "AI-generated workflow suggestions",
      "Webhook triggers and scheduled jobs",
      "200+ native integrations",
      "Real-time execution monitoring",
    ],
    version: "0.1.0",
    icon: Workflow,
    color: "#10B981",
    bgColor: "rgba(16, 185, 129, 0.15)",
    domain: "flow.dclawstack.io",
    category: "Automation",
    status: "live",
  },
  {
    id: "med",
    name: "DClaw Med",
    tagline: "Clinical intelligence at your fingertips",
    description:
      "HIPAA-compliant AI assistant for healthcare professionals. Clinical decision support, medical literature search, and patient note generation.",
    features: [
      "HIPAA-compliant data handling",
      "Clinical decision support",
      "Medical literature RAG search",
      "Patient note auto-generation",
      "Drug interaction checking",
    ],
    version: "0.1.0",
    icon: HeartPulse,
    color: "#EF4444",
    bgColor: "rgba(239, 68, 68, 0.15)",
    domain: "med.dclawstack.io",
    category: "Healthcare",
    status: "live",
  },
  {
    id: "learn",
    name: "DClaw Learn",
    tagline: "Adaptive learning that works",
    description:
      "AI tutor that adapts to your learning style. Upload any material and get personalized study plans, quizzes, and explanations.",
    features: [
      "Adaptive learning paths",
      "Upload PDFs, videos, or notes",
      "Auto-generated quizzes",
      "Spaced repetition scheduling",
      "Progress analytics dashboard",
    ],
    version: "0.1.0",
    icon: BookOpen,
    color: "#3B82F6",
    bgColor: "rgba(59, 130, 246, 0.15)",
    domain: "learn.dclawstack.io",
    category: "Education",
    status: "live",
  },
  {
    id: "seo",
    name: "DClaw SEO",
    tagline: "Rank higher with AI",
    description:
      "Autonomous SEO agent that audits your site, optimizes content, and tracks rankings. Complete technical SEO and content strategy in one tool.",
    features: [
      "Full site technical audits",
      "AI content optimization",
      "Competitor keyword analysis",
      "Rank tracking and alerts",
      "Automated backlink outreach",
    ],
    version: "0.1.0",
    icon: TrendingUp,
    color: "#10B981",
    bgColor: "rgba(16, 185, 129, 0.15)",
    domain: "seo.dclawstack.io",
    category: "Marketing",
    status: "live",
  },
  {
    id: "create",
    name: "DClaw Create",
    tagline: "Generate anything",
    description:
      "Multimodal AI content studio. Generate images, video, audio, 3D models, and code from natural language prompts.",
    features: [
      "Text-to-image generation",
      "Text-to-video synthesis",
      "Voice cloning and TTS",
      "3D model generation",
      "Batch processing pipeline",
    ],
    version: "0.1.0",
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
    description:
      "A desktop IDE built for the AI era. Inline copilot, multi-file refactoring, and autonomous agent coding with full project context.",
    features: [
      "Inline AI copilot",
      "Multi-file refactoring",
      "Autonomous coding agents",
      "Local and remote dev environments",
      "Built-in terminal and debugger",
    ],
    version: "0.1.0",
    icon: Code2,
    color: "#1F2937",
    bgColor: "rgba(31, 41, 55, 0.15)",
    domain: "code.dclawstack.io",
    category: "Development",
    status: "live",
  },
  {
    id: "agent",
    name: "DClaw Agent",
    tagline: "Build, share, and sell AI agents",
    description:
      "Agent marketplace and builder. Create custom AI agents with no-code tools, share them with the community, or sell them in the marketplace.",
    features: [
      "No-code agent builder",
      "Agent marketplace",
      "Revenue sharing for creators",
      "Agent-to-agent communication",
      "Usage analytics and billing",
    ],
    version: "0.1.0",
    icon: Bot,
    color: "#8B5CF6",
    bgColor: "rgba(139, 92, 246, 0.15)",
    domain: "agent.dclawstack.io",
    category: "Platform",
    status: "live",
  },
  {
    id: "rag",
    name: "DClaw RAG",
    tagline: "Universal knowledge retrieval",
    description:
      "Connect any data source — documents, databases, APIs — and get accurate, cited answers with full provenance. Enterprise-grade RAG without the complexity.",
    features: [
      "Multi-source data connectors",
      "Vector search with citations",
      "Real-time sync from APIs",
      "Custom embedding models",
      "Access control and audit logs",
    ],
    version: "0.1.0",
    icon: Search,
    color: "#F59E0B",
    bgColor: "rgba(245, 158, 11, 0.15)",
    domain: "rag.dclawstack.io",
    category: "Platform",
    status: "live",
  },
];
