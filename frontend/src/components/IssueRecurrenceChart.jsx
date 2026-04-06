const issueTypeLabel = {
  version_mismatch: "Version Mismatch",
  missing_dependency: "Missing Dependency",
  duplicate: "Duplicate Container",
  unsupported_combo: "Unsupported Combo",
  duplicate_part_number: "Duplicate Part Number",
  system_incompatible: "System Incompatible",
  aircraft_state_version_conflict: "Aircraft State Conflict",
  version_conflict: "Version Conflict",
  maturity_risk: "Maturity Risk",
};

export default function IssueRecurrenceChart({ data }) {
  if (!data || data.length === 0) {
    return <p className="text-gray-400 text-sm">No recurring issues</p>;
  }

  // Show top 5
  const topIssues = data.slice(0, 5);
  const maxCount = topIssues[0]?.count || 1;

  return (
    <div className="w-full space-y-4">
      {topIssues.map((issue, idx) => {
        const barWidth = (issue.count / maxCount) * 100;
        const displayLabel = issueTypeLabel[issue.type] || issue.type;

        return (
          <div key={idx} className="space-y-1">
            <div className="flex items-center justify-between">
              <label className="text-sm font-semibold text-gray-700 truncate">
                {displayLabel}
              </label>
              <span className="text-xs font-bold text-gray-600 bg-gray-100 px-2 py-1 rounded">
                {issue.count}x
              </span>
            </div>
            <div className="w-full bg-gray-200 rounded-full h-2">
              <div
                className="bg-orange-500 h-2 rounded-full transition-all"
                style={{ width: `${barWidth}%` }}
              />
            </div>
            <p className="text-xs text-gray-400">
              Fix Rate: {(issue.fix_rate * 100).toFixed(0)}%
            </p>
          </div>
        );
      })}
    </div>
  );
}
