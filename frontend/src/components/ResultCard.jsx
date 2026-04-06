import { useEffect, useRef, useState } from "react";
import RiskBadge from "./RiskBadge";
import ChatBot from "./ChatBot";
import { exportValidationReport } from "../services/api";

const severityIcon = {
  HIGH:   { icon: "↑", bg: "bg-red-50",   border: "border-red-100",  text: "text-red-600",   rec: "border-red-200 bg-red-50" },
  MEDIUM: { icon: "▲", bg: "bg-blue-50",  border: "border-blue-100", text: "text-blue-600",  rec: "border-blue-200 bg-blue-50" },
  LOW:    { icon: "i", bg: "bg-gray-50",  border: "border-gray-100", text: "text-gray-500",  rec: "border-gray-200 bg-gray-50" },
};

const issueTypeLabel = {
  version_mismatch:      "Version Mismatch",
  missing_dependency:    "Missing Dependency",
  duplicate:             "Duplicate Container",
  unsupported_combo:     "Unsupported Combination",
  duplicate_part_number: "Duplicate Part Number",
  system_incompatible:   "System Incompatible",
  aircraft_state_version_conflict: "Aircraft State Version Conflict",
  version_conflict:      "Version Conflict",
  maturity_risk:         "Maturity Risk",
};

export default function ResultCard({ result, releaseName, validationTarget, onRevalidate, activeTab }) {
  const high   = result.issues?.filter(i => i.severity === "HIGH").length || 0;
  const medium = result.issues?.filter(i => i.severity === "MEDIUM").length || 0;
  const low    = result.issues?.filter(i => i.severity === "LOW").length || 0;
  const total  = result.issues?.length || 0;
  const isDeploymentApproved = result.status === "PASS";
  const [selectedFormat, setSelectedFormat] = useState("csv");
  const [exporting, setExporting] = useState(false);
  const exportSectionRef = useRef(null);

  useEffect(() => {
    if (activeTab === "export" && exportSectionRef.current) {
      exportSectionRef.current.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  }, [activeTab]);

  const handleExport = async () => {
    setExporting(true);
    try {
      await exportValidationReport({
        releaseId: result.release_id,
        releaseName,
        format: selectedFormat,
      });
    } catch (err) {
      alert(err?.response?.data?.error || err?.message || "Failed to export report.");
    } finally {
      setExporting(false);
    }
  };

  return (
    <div className="space-y-5">
      {/* Top bar */}
      <div className="bg-white rounded-2xl border border-gray-100 shadow-sm px-8 py-6">
        <div className="flex items-center justify-between flex-wrap gap-4">
          <div>
            <p className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-1">Release</p>
            <h2 className="text-3xl font-black text-gray-900">{releaseName}</h2>
          </div>
          <div className="flex items-center gap-3">
            <RiskBadge risk={result.risk} large />
            {result.status && (
              <div className={`inline-flex items-center gap-2 px-4 py-2 rounded-lg border font-bold text-lg tracking-wide ${
                result.status === "PASS" 
                  ? "bg-green-50 border-green-200 text-green-600" 
                  : "bg-red-50 border-red-200 text-red-600"
              }`}>
                <span className={`w-2.5 h-2.5 rounded-full animate-pulse ${
                  result.status === "PASS" ? "bg-green-500" : "bg-red-500"
                }`} />
                {result.status === "PASS" ? "PASS" : "FAIL"}
              </div>
            )}
          </div>
        </div>

        {validationTarget && (
          <div className="mt-4 grid grid-cols-1 sm:grid-cols-2 gap-3">
            <div className="rounded-xl border border-gray-100 bg-gray-50 px-4 py-3">
              <p className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-1">Affected Tail Number</p>
              <p className="text-sm font-bold text-gray-800">{validationTarget.tailNumber}</p>
              <p className="text-xs text-gray-500 mt-0.5">Aircraft: {validationTarget.aircraftType}</p>
            </div>
            <div className="rounded-xl border border-blue-100 bg-blue-50 px-4 py-3">
              <p className="text-xs font-bold tracking-widest text-blue-400 uppercase mb-1">Affected Fleet / System</p>
              <p className="text-sm font-bold text-blue-700">{validationTarget.targetFleet}</p>
              <p className="text-xs text-blue-600 mt-0.5">System: {validationTarget.aircraftSystem}</p>
            </div>
          </div>
        )}

        {/* Stat pills */}
        <div className="grid grid-cols-2 sm:grid-cols-4 gap-3 mt-6">
          {[
            { count: high,   label: "High Risk",   sub: "Must Fix",  color: "text-red-600",   bg: "bg-red-50",   border: "border-red-100"  },
            { count: medium, label: "Medium Risk",  sub: "Should Fix", color: "text-blue-600",  bg: "bg-blue-50",  border: "border-blue-100" },
            { count: low,    label: "Low Risk",     sub: "Nice to Fix",color: "text-gray-500",  bg: "bg-gray-50",  border: "border-gray-100" },
            { count: result.risk === "SAFE" ? result.issues?.length || 0 : 0, label: "Passed", sub: "No Issues", color: "text-green-600", bg: "bg-green-50", border: "border-green-100" },
          ].map((s) => (
            <div key={s.label} className={`rounded-xl border px-4 py-3 flex items-center justify-between ${s.bg} ${s.border}`}>
              <div>
                <p className={`text-2xl font-black ${s.color}`}>{s.count}</p>
                <p className={`text-xs font-semibold ${s.color}`}>{s.label}</p>
                <p className="text-xs text-gray-400">{s.sub}</p>
              </div>
              <div className={`w-8 h-8 rounded-full flex items-center justify-center text-lg ${s.color} opacity-30 font-black`}>
                {s.label === "Passed" ? "✓" : s.label === "High Risk" ? "✕" : s.label === "Medium Risk" ? "▲" : "i"}
              </div>
            </div>
          ))}
        </div>
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-3 gap-5">
        {/* Issues list */}
        <div className="lg:col-span-2 space-y-4">
          <div className="flex items-center gap-2">
            <h3 className="font-bold text-gray-800 text-lg">Issues Detected</h3>
            <span className="text-xs font-bold text-gray-400 bg-gray-100 px-2 py-0.5 rounded-full">{total} Total</span>
          </div>

          {total === 0 ? (
            <div className="bg-green-50 border border-green-100 rounded-xl p-6 text-center">
              <p className="text-green-600 font-bold text-lg">All checks passed!</p>
              <p className="text-green-400 text-sm mt-1">This release is safe to deploy.</p>
            </div>
          ) : (
            <div className="max-h-144 overflow-y-auto pr-1 space-y-4">
              {result.issues.map((issue, idx) => {
                const s = severityIcon[issue.severity] || severityIcon.LOW;
                const rec = result.recommendations?.find(r => r.issue_type === issue.type);
                return (
                  <div key={idx} className={`bg-white rounded-xl border ${s.border} shadow-sm overflow-hidden`}>
                    <div className={`px-5 py-4 ${s.bg} border-b ${s.border}`}>
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className={`w-8 h-8 rounded-full flex items-center justify-center font-black text-sm ${s.bg} border ${s.border} ${s.text}`}>
                            {s.icon}
                          </div>
                          <div>
                            <p className="font-bold text-gray-900 text-sm">{issueTypeLabel[issue.type] || issue.type}</p>
                            <p className="text-xs text-gray-500 mt-0.5">{issue.message}</p>
                          </div>
                        </div>
                        <RiskBadge risk={issue.severity} />
                      </div>
                    </div>

                    {rec && (
                      <div className="px-5 py-3 bg-white">
                        <p className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-1">Recommendation</p>
                        <p className="text-sm text-gray-700">{rec.action}</p>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>
          )}
        </div>

        {/* Right column — AI Insights + Next Actions */}
        <div className="space-y-4">
          {/* AI Insight */}
          <div className="bg-teal-800 rounded-xl p-5 text-white">
            <div className="flex items-center gap-2 mb-3">
              <span className="text-teal-300 text-lg">✦</span>
              <p className="text-xs font-bold tracking-widest text-teal-300 uppercase">AI Insights</p>
            </div>
            {typeof result.insight === "object" ? (
              <div className="space-y-3">
                <p className="text-sm text-teal-100 leading-relaxed italic">
                  "{result.insight.summary}"
                </p>
                <div className="border-t border-teal-600 pt-3">
                  <p className="text-xs font-bold text-teal-200 mb-1">Deployment Decision:</p>
                  <p className="text-sm text-teal-100 font-semibold">{result.insight.impact}</p>
                </div>
              </div>
            ) : (
              <p className="text-sm text-teal-100 leading-relaxed italic">
                "{result.insight}"
              </p>
            )}
          </div>

          {/* Next Actions */}
          <div className="bg-white border border-gray-100 rounded-xl p-5 shadow-sm">
            <p className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-4">Next Actions</p>
            <div className="space-y-3">
              <button
                onClick={onRevalidate}
                className="w-full py-3 rounded-xl bg-blue-600 hover:bg-blue-700 text-white font-bold text-sm flex items-center justify-center gap-2 transition-all active:scale-95"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Go-back
              </button>
              {isDeploymentApproved ? (
                <button onClick={()=>alert("The release is deployed successfully")} className="w-full py-3 rounded-xl bg-green-600 hover:bg-green-700 text-white font-bold text-sm flex items-center justify-center gap-2 transition-all active:scale-95">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                  </svg>
                  Approve & Deploy
                </button>
              ) : (
                <button onClick={()=>alert("Proceeding with High Risk issues may lead to production failure")} className="w-full py-3 rounded-xl border-2 border-gray-200 text-gray-600 hover:bg-gray-50 font-bold text-sm flex items-center justify-center gap-2 transition-all active:scale-95">
                  <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 9l3 3m0 0l-3 3m3-3H8m13 0a9 9 0 11-18 0 9 9 0 0118 0z" />
                  </svg>
                  Proceed Anyway
                </button>
              )}
              {!isDeploymentApproved && (
                <p className="text-xs text-gray-400 text-center leading-relaxed">
                  Note: Proceeding with High Risk issues may lead to production failure.
                </p>
              )}
            </div>
          </div>

          <div ref={exportSectionRef} className="bg-white border border-gray-100 rounded-xl p-5 shadow-sm">
            <p className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-4">Export Report</p>
            <div className="space-y-3">
              <select
                value={selectedFormat}
                onChange={(e) => setSelectedFormat(e.target.value)}
                className="w-full rounded-xl border border-gray-200 px-3 py-2.5 text-sm text-gray-700 focus:outline-none focus:ring-2 focus:ring-blue-200"
              >
                <option value="csv">CSV (.csv)</option>
                <option value="pdf">PDF (.pdf)</option>
                <option value="xlsx">Excel (.xlsx)</option>
              </select>
              <button
                onClick={handleExport}
                disabled={exporting}
                className="w-full py-3 rounded-xl bg-gray-900 hover:bg-black text-white font-bold text-sm flex items-center justify-center gap-2 transition-all active:scale-95"
              >
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 3v12m0 0l4-4m-4 4l-4-4m-6 8h20" />
                </svg>
                {exporting ? "Preparing download..." : `Download ${selectedFormat.toUpperCase()}`}
              </button>
            </div>
          </div>

          {/* Validated at */}
          <p className="text-xs text-gray-300 text-center">
            Validated at {new Date(result.validated_at).toLocaleString()}
          </p>
        </div>
      </div>

      <ChatBot result={result} />
    </div>
  );
}