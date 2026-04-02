import { useState } from "react";
import { mockReleases } from "../mockdata.js";

export default function ReleaseSelector({ onValidate, loading }) {
  const [selected, setSelected] = useState(null);

  const handleSelect = (release) => setSelected(release);

  return (
    <div className="bg-white rounded-2xl border border-gray-100 shadow-sm overflow-hidden">
      {/* Header */}
      <div className="px-8 pt-8 pb-6 border-b border-gray-100">
        <div className="flex items-center justify-between">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Upload Release Package</h2>
            <p className="text-gray-400 mt-1 text-sm">Initialize a new content audit by selecting a pre-configured mock environment.</p>
          </div>
          <span className="text-xs font-semibold tracking-widest text-teal-600 bg-teal-50 border border-teal-100 px-3 py-1.5 rounded-full uppercase">System Active</span>
        </div>
      </div>

      <div className="p-8 grid grid-cols-1 lg:grid-cols-3 gap-6">
        {/* Left — Mock selector */}
        <div className="lg:col-span-2 space-y-4">
          <p className="text-xs font-bold tracking-widest text-gray-400 uppercase">Select Mock Package</p>
          <div className="space-y-3">
            {mockReleases.map((release) => (
              <button
                key={release.id}
                onClick={() => handleSelect(release)}
                className={`w-full text-left px-5 py-4 rounded-xl border-2 transition-all duration-200 ${
                  selected?.id === release.id
                    ? "border-blue-500 bg-blue-50 shadow-sm"
                    : "border-gray-100 bg-gray-50 hover:border-gray-200 hover:bg-white"
                }`}
              >
                <div className="flex items-center gap-3">
                  <div className={`w-4 h-4 rounded-full border-2 flex items-center justify-center shrink-0 ${
                    selected?.id === release.id ? "border-blue-500" : "border-gray-300"
                  }`}>
                    {selected?.id === release.id && (
                      <div className="w-2 h-2 rounded-full bg-blue-500" />
                    )}
                  </div>
                  <div>
                    <p className={`font-semibold text-sm ${selected?.id === release.id ? "text-blue-700" : "text-gray-700"}`}>
                      {release.payload.release_name}
                    </p>
                    <p className="text-xs text-gray-400 mt-0.5">
                      {release.payload.target_fleet} · {release.payload.containers.length} containers
                    </p>
                  </div>
                </div>
              </button>
            ))}
          </div>

          {/* JSON Preview */}
          {selected && (
            <div className="mt-4">
              <p className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-2">Payload Preview</p>
              <pre className="bg-gray-950 text-green-400 text-xs p-4 rounded-xl overflow-auto max-h-48 font-mono leading-relaxed">
                {JSON.stringify(selected.payload, null, 2)}
              </pre>
            </div>
          )}
        </div>

        {/* Right — Audit params + button */}
        <div className="space-y-4">
          <div className="bg-gray-50 border border-gray-100 rounded-xl p-5">
            <p className="text-xs font-bold tracking-widest text-gray-400 uppercase mb-4">Audit Parameters</p>
            <div className="space-y-3">
              {[
                { icon: "✓", label: "Schema Validation", desc: "Automatic JSON schema check" },
                { icon: "⬡", label: "Version Checks", desc: "Dependency version constraints" },
                { icon: "⚑", label: "Duplicate Detection", desc: "Duplicate container scan" },
                { icon: "◈", label: "AI Insights", desc: "AI-powered insights" },
              ].map((item) => (
                <div key={item.label} className="flex items-start gap-3">
                  <div className="w-6 h-6 rounded-full bg-blue-100 text-blue-600 flex items-center justify-center text-xs font-bold shrink-0 mt-0.5">
                    {item.icon}
                  </div>
                  <div>
                    <p className="text-sm font-semibold text-gray-700">{item.label}</p>
                    <p className="text-xs text-gray-400">{item.desc}</p>
                  </div>
                </div>
              ))}
            </div>
          </div>

          {/* Validate button */}
          <button
            onClick={() => selected && onValidate(selected.payload)}
            disabled={!selected || loading}
            className={`w-full py-4 rounded-xl font-bold text-sm tracking-wide transition-all duration-200 flex items-center justify-center gap-2 ${
              selected && !loading
                ? "bg-blue-600 hover:bg-blue-700 text-white shadow-lg shadow-blue-200 active:scale-95"
                : "bg-gray-100 text-gray-400 cursor-not-allowed"
            }`}
          >
            {loading ? (
              <>
                <svg className="animate-spin w-4 h-4" fill="none" viewBox="0 0 24 24">
                  <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                  <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
                </svg>
                Validating...
              </>
            ) : (
              <>
                Validate Package
                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 5l7 7m0 0l-7 7m7-7H3" />
                </svg>
              </>
            )}
          </button>
        </div>
      </div>
    </div>
  );
}