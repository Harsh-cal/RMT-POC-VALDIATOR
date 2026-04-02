import { useState } from "react";
import ReleaseSelector from "./components/ReleaseSelector";
import ResultCard from "./components/ResultCard";
import FileUploadModal from "./components/FileUploadModal";
import { validateRelease } from "./services/api";

export default function App() {
  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const [releaseName, setReleaseName] = useState("");
  const [activeTab, setActiveTab] = useState("dashboard");
  const [showModal, setShowModal] = useState(false);

  const handleValidate = async (payload) => {
    setLoading(true);
    setError(null);
    setResult(null);
    setReleaseName(payload.release_name);
    setShowModal(false);
    try {
      const data = await validateRelease(payload);
      setResult(data);
    } catch (err) {
      setError(err?.response?.data?.message || "Something went wrong. Is the backend running?");
    } finally {
      setLoading(false);
    }
  };

  const handleRevalidate = () => {
    setResult(null);
    setError(null);
  };

  return (
    <div className="min-h-screen bg-gray-50 font-sans">

      {showModal && (
        <FileUploadModal
          onClose={() => setShowModal(false)}
          onValidate={handleValidate}
          loading={loading}
        />
      )}

      <header className="bg-white border-b border-gray-100 sticky top-0 z-10">
        <div className="max-w-7xl mx-auto px-6 h-14 flex items-center justify-between">
          <div className="flex items-center gap-8">
            <span className="font-black text-gray-900 text-base tracking-tight">
              Smart Content Validator <span className="text-gray-300 font-normal">(POC)</span>
            </span>
            <nav className="hidden md:flex items-center gap-1">
              {["Dashboard", "History", "Settings"].map((tab) => (
                <button
                  key={tab}
                  onClick={() => setActiveTab(tab.toLowerCase())}
                  className={`px-3 py-1.5 rounded-lg text-sm font-medium transition-colors ${
                    activeTab === tab.toLowerCase()
                      ? "text-blue-600 bg-blue-50"
                      : "text-gray-400 hover:text-gray-600"
                  }`}
                >
                  {tab}
                </button>
              ))}
            </nav>
          </div>
          <button
            onClick={() => setShowModal(true)}
            className="bg-blue-600 hover:bg-blue-700 text-white text-sm font-semibold px-4 py-2 rounded-lg transition-all active:scale-95 flex items-center gap-2"
          >
            <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M4 16v1a3 3 0 003 3h10a3 3 0 003-3v-1m-4-8l-4-4m0 0L8 8m4-4v12"/>
            </svg>
            Validate Another Package
          </button>
        </div>
      </header>

      <div className="max-w-7xl mx-auto px-6 py-6 flex gap-6">
        <aside className="hidden lg:flex flex-col w-52 shrink-0">
          <div className="bg-white border border-gray-100 rounded-2xl p-4 shadow-sm">
            <p className="text-xs font-black text-blue-600 mb-0.5">Smart Validator </p>
            <p className="text-xs text-gray-400 mb-4">Internal Review</p>
            <nav className="space-y-1">
              {[
                { icon: "◈", label: "Validation Results", key: "dashboard" },
                { icon: "↓", label: "Export",        key: "export"    },
              ].map((item) => (
                <button
                  key={item.key}
                  onClick={() => setActiveTab(item.key)}
                  className={`w-full flex items-center gap-2.5 px-3 py-2 rounded-lg text-sm font-medium transition-colors text-left ${
                    activeTab === item.key
                      ? "bg-blue-50 text-blue-600"
                      : "text-gray-500 hover:bg-gray-50 hover:text-gray-700"
                  }`}
                >
                  <span className="text-base">{item.icon}</span>
                  {item.label}
                </button>
              ))}
            </nav>
          </div>
        </aside>

        <main className="flex-1 min-w-0 space-y-6">
          {error && (
            <div className="bg-red-50 border border-red-200 rounded-xl px-5 py-4 flex items-center gap-3">
              <span className="text-red-500 text-xl">⚠</span>
              <div>
                <p className="font-semibold text-red-700 text-sm">Validation Failed</p>
                <p className="text-red-500 text-xs mt-0.5">{error}</p>
              </div>
              <button onClick={() => setError(null)} className="ml-auto text-red-300 hover:text-red-500 text-lg">✕</button>
            </div>
          )}

          {loading && (
            <div className="bg-white border border-gray-100 rounded-2xl p-12 flex flex-col items-center justify-center shadow-sm">
              <svg className="animate-spin w-10 h-10 text-blue-500 mb-4" fill="none" viewBox="0 0 24 24">
                <circle className="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" strokeWidth="4"/>
                <path className="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8v8z"/>
              </svg>
              <p className="font-bold text-gray-700">Validating {releaseName}...</p>
              <p className="text-sm text-gray-400 mt-1">Running rule engine + AI insights</p>
            </div>
          )}

          {!loading && !result && (
            <ReleaseSelector onValidate={handleValidate} loading={loading} />
          )}

          {!loading && result && (
            <ResultCard
              result={result}
              releaseName={releaseName}
              onRevalidate={handleRevalidate}
              activeTab={activeTab}
            />
          )}
        </main>
      </div>
    </div>
  );
}