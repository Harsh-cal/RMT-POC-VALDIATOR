const config = {
  HIGH:   { bg: "bg-red-50",    text: "text-red-600",    border: "border-red-200",   dot: "bg-red-500",    label: "HIGH RISK" },
  MEDIUM: { bg: "bg-blue-50",   text: "text-blue-600",   border: "border-blue-200",  dot: "bg-blue-500",   label: "MEDIUM RISK" },
  LOW:    { bg: "bg-green-50",  text: "text-green-600",  border: "border-green-200", dot: "bg-green-500",  label: "LOW RISK" },
  SAFE:   { bg: "bg-green-50",  text: "text-green-600",  border: "border-green-200", dot: "bg-green-500",  label: "SAFE" },
  PASS:   { bg: "bg-green-50",  text: "text-green-600",  border: "border-green-200", dot: "bg-green-500",  label: "GO" },
  FAILED: { bg: "bg-red-50",    text: "text-red-600",    border: "border-red-200",   dot: "bg-red-500",    label: "NO-GO" },
};

export default function RiskBadge({ risk, large = false }) {
  const c = config[risk] || config.SAFE;
  if (large) {
    return (
      <div className={`inline-flex items-center gap-2 px-4 py-2 rounded-lg border ${c.bg} ${c.border}`}>
        <span className={`w-2.5 h-2.5 rounded-full ${c.dot} animate-pulse`} />
        <span className={`font-bold text-2xl tracking-wide ${c.text}`}>{risk}</span>
      </div>
    );
  }
  return (
    <span className={`inline-flex items-center gap-1.5 px-2.5 py-1 rounded-full text-xs font-semibold border ${c.bg} ${c.border} ${c.text}`}>
      <span className={`w-1.5 h-1.5 rounded-full ${c.dot}`} />
      {c.label}
    </span>
  );
}