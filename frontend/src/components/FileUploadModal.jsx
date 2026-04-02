import { useState, useRef, useCallback } from "react";

const REQUIRED_FIELDS = ["release_name", "version", "target_fleet", "containers", "aircraft"];

function validateSchema(parsed) {
  const missing = REQUIRED_FIELDS.filter(f => !(f in parsed));
  if (missing.length > 0) return `Missing required fields: ${missing.join(", ")}`;
  if (!Array.isArray(parsed.containers)) return '"containers" must be an array';
  if (parsed.containers.length === 0) return '"containers" array cannot be empty';
  if (typeof parsed.aircraft !== 'object') return '"aircraft" must be an object';
  if (!parsed.aircraft.tailNumber) return '"aircraft.tailNumber" is required';
  if (!parsed.aircraft.system) return '"aircraft.system" is required';
  for (let i = 0; i < parsed.containers.length; i++) {
    const c = parsed.containers[i];
    if (!c.name) return `Container at index ${i} is missing "name"`;
    if (!c.version) return `Container at index ${i} is missing "version"`;
    if (!c.partNumber) return `Container at index ${i} is missing "partNumber"`;
    if (!c.systemType) return `Container at index ${i} is missing "systemType"`;
  }
  return null;
}

export default function FileUploadModal({ onClose, onValidate, loading }) {
  const [dragOver, setDragOver] = useState(false);
  const [file, setFile] = useState(null);
  const [parsed, setParsed] = useState(null);
  const [parseError, setParseError] = useState(null);
  const [schemaError, setSchemaError] = useState(null);
  const inputRef = useRef();

  const processFile = useCallback((f) => {
    setFile(f);
    setParsed(null);
    setParseError(null);
    setSchemaError(null);

    const reader = new FileReader();
    reader.onload = (e) => {
      try {
        const content = e.target.result;
        const json = JSON.parse(content);
        const err = validateSchema(json);
        if (err) {
          setSchemaError(err);
        } else {
          setParsed(json);
        }
      } catch {
        setParseError("Invalid JSON — could not parse this file. Please check the format.");
      }
    };
    reader.readAsText(f);
  }, []);

  const onDrop = useCallback((e) => {
    e.preventDefault();
    setDragOver(false);
    const f = e.dataTransfer.files[0];
    if (f) processFile(f);
  }, [processFile]);

  const onFileChange = (e) => {
    const f = e.target.files[0];
    if (f) processFile(f);
  };

  const handleValidate = () => {
    if (parsed) onValidate(parsed);
  };

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center p-4"
      style={{ background: "rgba(0,0,0,0.45)" }}
      onClick={(e) => e.target === e.currentTarget && onClose()}
    >
      <div className="bg-white rounded-2xl shadow-2xl w-full max-w-2xl overflow-hidden animate-in">

        {/* Header */}
        <div className="px-7 pt-7 pb-5 border-b border-gray-100 flex items-start justify-between">
          <div>
            <h2 className="text-xl font-black text-gray-900">Upload Release Package</h2>
            <p className="text-sm text-gray-400 mt-1">Drop a <code className="bg-gray-100 px-1 rounded text-xs">.json</code> or <code className="bg-gray-100 px-1 rounded text-xs">.txt</code> file — we'll parse and validate it in real time.</p>
          </div>
          <button
            onClick={onClose}
            className="text-gray-300 hover:text-gray-500 text-2xl leading-none mt-0.5 transition-colors"
          >
            ✕
          </button>
        </div>

        <div className="p-7 space-y-5">

          {/* Drop zone */}
          <div
            onDragOver={(e) => { e.preventDefault(); setDragOver(true); }}
            onDragLeave={() => setDragOver(false)}
            onDrop={onDrop}
            onClick={() => inputRef.current.click()}
            className={`relative cursor-pointer border-2 border-dashed rounded-xl p-10 flex flex-col items-center justify-center text-center transition-all duration-200 ${
              dragOver
                ? "border-blue-400 bg-blue-50"
                : file && !parseError && !schemaError
                ? "border-green-300 bg-green-50"
                : parseError || schemaError
                ? "border-red-300 bg-red-50"
                : "border-gray-200 bg-gray-50 hover:border-blue-300 hover:bg-blue-50"
            }`}
          >
            <input
              ref={inputRef}
              type="file"
              accept=".json,.txt"
              className="hidden"
              onChange={onFileChange}
            />

            {/* Icon */}
            <div className={`w-14 h-14 rounded-2xl flex items-center justify-center mb-4 text-2xl ${
              dragOver ? "bg-blue-100" :
              file && !parseError && !schemaError ? "bg-green-100" :
              parseError || schemaError ? "bg-red-100" :
              "bg-gray-100"
            }`}>
              {dragOver ? "↓" :
               file && !parseError && !schemaError ? "✓" :
               parseError || schemaError ? "✕" : "↑"}
            </div>

            {!file ? (
              <>
                <p className="font-bold text-gray-700 text-sm">Drop your JSON package here</p>
                <p className="text-gray-400 text-xs mt-1">or click to browse files</p>
                <p className="text-gray-300 text-xs mt-3">Supports .json and .txt files</p>
              </>
            ) : parseError ? (
              <>
                <p className="font-bold text-red-600 text-sm">{file.name}</p>
                <p className="text-red-400 text-xs mt-1">Parse error — click to try another file</p>
              </>
            ) : schemaError ? (
              <>
                <p className="font-bold text-red-600 text-sm">{file.name}</p>
                <p className="text-red-400 text-xs mt-1">Schema error — click to try another file</p>
              </>
            ) : (
              <>
                <p className="font-bold text-green-700 text-sm">{file.name}</p>
                <p className="text-green-500 text-xs mt-1">
                  {parsed?.containers?.length} containers · {parsed?.target_fleet}
                </p>
                <p className="text-gray-300 text-xs mt-2">Click to replace file</p>
              </>
            )}
          </div>

          {/* Error messages */}
          {parseError && (
            <div className="flex items-start gap-3 bg-red-50 border border-red-100 rounded-xl px-4 py-3">
              <span className="text-red-400 text-lg mt-0.5">⚠</span>
              <div>
                <p className="text-sm font-semibold text-red-700">Parse Error</p>
                <p className="text-xs text-red-500 mt-0.5">{parseError}</p>
              </div>
            </div>
          )}

          {schemaError && (
            <div className="flex items-start gap-3 bg-red-50 border border-red-100 rounded-xl px-4 py-3">
              <span className="text-red-400 text-lg mt-0.5">⚠</span>
              <div>
                <p className="text-sm font-semibold text-red-700">Schema Validation Error</p>
                <p className="text-xs text-red-500 mt-0.5">{schemaError}</p>
              </div>
            </div>
          )}

          {/* Parsed preview */}
          {parsed && (
            <div>
              <div className="flex items-center justify-between mb-2">
                <p className="text-xs font-bold tracking-widest text-gray-400 uppercase">Parsed Payload</p>
                <span className="text-xs text-green-600 font-semibold bg-green-50 border border-green-100 px-2 py-0.5 rounded-full">Valid JSON</span>
              </div>
              <pre className="bg-gray-950 text-green-400 text-xs p-4 rounded-xl overflow-auto max-h-48 font-mono leading-relaxed">
                {JSON.stringify(parsed, null, 2)}
              </pre>
            </div>
          )}

          {/* Actions */}
          <div className="flex items-center gap-3 pt-1">
            <button
              onClick={onClose}
              className="flex-1 py-3 rounded-xl border-2 border-gray-100 text-gray-500 hover:bg-gray-50 font-semibold text-sm transition-all"
            >
              Cancel
            </button>
            <button
              onClick={handleValidate}
              disabled={!parsed || loading}
              className={`flex-1 py-3 rounded-xl font-bold text-sm flex items-center justify-center gap-2 transition-all active:scale-95 ${
                parsed && !loading
                  ? "bg-blue-600 hover:bg-blue-700 text-white shadow-lg shadow-blue-100"
                  : "bg-gray-100 text-gray-300 cursor-not-allowed"
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
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M14 5l7 7m0 0l-7 7m7-7H3"/>
                  </svg>
                </>
              )}
            </button>
          </div>
        </div>
      </div>
    </div>
  );
}