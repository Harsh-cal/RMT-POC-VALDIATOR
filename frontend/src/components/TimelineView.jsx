import { useState } from "react";
import { exportValidationReport } from "../services/api";

import {ChevronDown, ChevronRight} from "lucide-react";

const riskBadgeColor = {
  HIGH: "bg-red-50 border-red-100 text-red-600",
  MEDIUM: "bg-yellow-50 border-yellow-100 text-yellow-600",
  LOW: "bg-blue-50 border-blue-100 text-blue-600",
  SAFE: "bg-green-50 border-green-100 text-green-600",
};

const statusIcon = {
  PASS: "✅",
  FAILED: "🔴",
};

const issueTypeLabel = {
  version_mismatch: "Version Mismatch",
  missing_dependency: "Missing Dependency",
  duplicate: "Duplicate Container",
  unsupported_combo: "Unsupported Combination",
  duplicate_part_number: "Duplicate Part Number",
  system_incompatible: "System Incompatible",
  version_conflict: "Version Conflict",
  maturity_risk: "Maturity Risk",
};

function DetailsModal({ release, onClose }) {
  const [format, setFormat] = useState("csv");
  const [exporting, setExporting] = useState(false);

  const handleExport = async () => {
    setExporting(true);
    try {
      await exportValidationReport({
        releaseId: release.release_id,
        releaseName: release.release_name,
        format,
      });
    } catch (err) {
      console.error("Export failed:", err);
      alert("Export failed. Please try again.");
    } finally {
      setExporting(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/50">
      <div className="bg-white rounded-2xl shadow-xl w-full max-w-lg mx-4 max-h-[90vh] overflow-y-auto">
        {/* Header */}
        <div className="flex items-center justify-between px-6 pt-6 pb-4 border-b border-gray-100">
          <div>
            <h2 className="text-lg font-black text-gray-900">
              {release.release_name} v{release.version}
            </h2>
            <p className="text-xs text-gray-400 mt-0.5">
              {new Date(release.validated_at).toLocaleString()}
            </p>
          </div>
          <button
            onClick={onClose}
            className="text-gray-400 hover:text-gray-700 text-xl font-bold leading-none"
          >
            ✕
          </button>
        </div>

        {/* Body */}
        <div className="px-6 py-5 space-y-4">
          {/* Status + Risk row */}
          <div className="flex gap-3">
            <span
              className={`px-3 py-1 rounded-full text-xs font-bold border ${
                release.status === "PASS"
                  ? "bg-green-50 border-green-200 text-green-700"
                  : "bg-red-50 border-red-200 text-red-700"
              }`}
            >
              {statusIcon[release.status]} {release.status}
            </span>
            <span
              className={`px-3 py-1 rounded-full text-xs font-bold border ${
                riskBadgeColor[release.risk] || "bg-gray-100"
              }`}
            >
              {release.risk} RISK
            </span>
            <span className="px-3 py-1 rounded-full text-xs font-bold border bg-gray-100 border-gray-200 text-gray-600">
              {release.issue_count} Issues
            </span>
          </div>

          {/* Release ID */}
          <div>
            <p className="text-xs font-bold text-gray-500 uppercase mb-1">Release ID</p>
            <p className="font-mono text-xs text-gray-700 bg-gray-50 px-3 py-2 rounded-lg break-all">
              {release.release_id}
            </p>
          </div>

          {/* Top Issues */}
          <div>
            <p className="text-xs font-bold text-gray-500 uppercase mb-2">Top Issues</p>
            {release.top_issues && release.top_issues.length > 0 ? (
              <ul className="space-y-1">
                {release.top_issues.map((issue, i) => (
                  <li
                    key={i}
                    className="flex items-center gap-2 text-sm text-gray-700 bg-gray-50 px-3 py-2 rounded-lg"
                  >
                    <span className="w-2 h-2 rounded-full bg-orange-400 shrink-0" />
                    {issueTypeLabel[issue] || issue}
                  </li>
                ))}
              </ul>
            ) : (
              <p className="text-sm text-green-600 font-semibold bg-green-50 px-3 py-2 rounded-lg">
                ✓ No issues detected
              </p>
            )}
          </div>

          {/* Export Row */}
          <div className="pt-2 border-t border-gray-100">
            <p className="text-xs font-bold text-gray-500 uppercase mb-2">Export Report</p>
            <div className="flex gap-2 items-center">
              <select
                value={format}
                onChange={(e) => setFormat(e.target.value)}
                className="text-sm border border-gray-200 rounded-lg px-3 py-1.5 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                <option value="csv">CSV</option>
                <option value="pdf">PDF</option>
                <option value="excel">Excel</option>
              </select>
              <button
                onClick={handleExport}
                disabled={exporting}
                className="text-sm px-4 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors disabled:opacity-60 font-semibold"
              >
                {exporting ? "Exporting..." : "⬇ Download"}
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}

export default function TimelineView({ releases }) {
  const [expandedIdx, setExpandedIdx] = useState(null);
  const [detailRelease, setDetailRelease] = useState(null);

  if (!releases || releases.length === 0) {
    return <p className="text-gray-400 text-sm">No releases in history</p>;
  }

  return (
    <>
      {detailRelease && (
        <DetailsModal release={detailRelease} onClose={() => setDetailRelease(null)} />
      )}

      <div className="space-y-3 max-h-96 overflow-y-auto">
        {releases.map((release, idx) => (
          <div key={idx} className="bg-gray-50 border border-gray-100 rounded-lg p-4">
            <div
              className="cursor-pointer flex items-center justify-between"
              onClick={() => setExpandedIdx(expandedIdx === idx ? null : idx)}
            >
              <div className="flex items-center gap-3 flex-1 min-w-0">
                <span className="text-xl">{statusIcon[release.status] || "—"}</span>

                <div className="flex-1 min-w-0">
                  <p className="font-bold text-gray-900 text-sm truncate">
                    {release.release_name} v{release.version}
                  </p>
                  <p className="text-xs text-gray-500">
                    {new Date(release.validated_at).toLocaleString()}
                  </p>
                </div>

                <div className="flex items-center gap-2">
                  <span
                    className={`inline-block px-2 py-1 rounded text-xs font-semibold border ${
                      riskBadgeColor[release.risk] || "bg-gray-100"
                    }`}
                  >
                    {release.risk}
                  </span>
                  <span className="bg-gray-200 text-gray-700 px-2 py-1 rounded text-xs font-semibold">
                    {release.issue_count} Issues
                  </span>
                </div>
              </div>

              <span className="text-gray-400 ml-2">{expandedIdx === idx ? <ChevronDown strokeWidth={1} /> : <ChevronRight strokeWidth={1} />}</span>
            </div>

            {/* Expanded Details */}
            {expandedIdx === idx && (
              <div className="mt-4 pt-4 border-t border-gray-200 space-y-2">
                {release.top_issues && release.top_issues.length > 0 ? (
                  <>
                    <p className="text-xs font-bold text-gray-700 uppercase">Top Issues:</p>
                    <ul className="space-y-1">
                      {release.top_issues.map((issue, i) => (
                        <li key={i} className="text-xs text-gray-600 flex items-center gap-2">
                          <span className="w-1.5 h-1.5 rounded-full bg-gray-400" />
                          {issueTypeLabel[issue] || issue}
                        </li>
                      ))}
                    </ul>
                  </>
                ) : (
                  <p className="text-xs text-green-600 font-semibold">✓ No issues detected</p>
                )}

                <p className="text-xs text-gray-500 mt-3 font-mono truncate">
                  ID: {release.release_id}
                </p>

                <div className="flex gap-2 mt-3">
                  <button
                    onClick={(e) => { e.stopPropagation(); setDetailRelease(release); }}
                    className="text-xs px-3 py-1.5 bg-blue-600 text-white rounded-lg hover:bg-blue-700 transition-colors font-semibold"
                  >
                    View Details
                  </button>
                  <button
                    onClick={(e) => { e.stopPropagation(); setDetailRelease(release); }}
                    className="text-xs px-3 py-1.5 bg-gray-200 text-gray-700 rounded-lg hover:bg-gray-300 transition-colors font-semibold"
                  >
                    Export
                  </button>
                </div>
              </div>
            )}
          </div>
        ))}
      </div>
    </>
  );
}
