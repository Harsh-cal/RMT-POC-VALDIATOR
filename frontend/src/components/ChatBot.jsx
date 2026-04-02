import { useEffect, useRef, useState } from "react";
import { askValidationBot } from "../services/api";

function Message({ msg }) {
  const isUser = msg.role === "user";

  return (
    <div className={`flex gap-3 ${isUser ? "flex-row-reverse" : "flex-row"}`}>
      <div
        className={`w-7 h-7 rounded-full shrink-0 flex items-center justify-center text-xs font-bold mt-0.5 ${
          isUser ? "bg-blue-600 text-white" : "bg-teal-700 text-white"
        }`}
      >
        {isUser ? "U" : "AI"}
      </div>
      <div
        className={`max-w-[80%] px-4 py-2.5 rounded-2xl text-sm leading-relaxed ${
          isUser
            ? "bg-blue-600 text-white rounded-tr-sm"
            : "bg-gray-100 text-gray-800 rounded-tl-sm"
        }`}
      >
        {msg.content}
      </div>
    </div>
  );
}

export default function ChatBot({ result }) {
  const [messages, setMessages] = useState([
    {
      role: "assistant",
      content: `Hi! I have context of ${result.release_name}. Status is ${result.status}, risk is ${result.risk}, and ${result.issues?.length || 0} issue(s) were detected. Ask what failed and what to change.`,
    },
  ]);
  const [input, setInput] = useState("");
  const [loading, setLoading] = useState(false);
  const [isOpen, setIsOpen] = useState(false);
  const bottomRef = useRef(null);

  useEffect(() => {
    bottomRef.current?.scrollIntoView({ behavior: "smooth" });
  }, [messages]);

  const sendMessage = async (seedText) => {
    const text = (seedText ?? input).trim();
    if (!text || loading) return;

    const userMsg = { role: "user", content: text };
    const updated = [...messages, userMsg];

    setMessages(updated);
    setInput("");
    setLoading(true);

    try {
      const data = await askValidationBot({
        question: text,
        result,
        history: updated.slice(-8),
      });

      setMessages((prev) => [
        ...prev,
        { role: "assistant", content: data?.answer || "No response from assistant." },
      ]);
    } catch {
      setMessages((prev) => [
        ...prev,
        {
          role: "assistant",
          content: "I could not fetch an answer right now. Please check backend availability and try again.",
        },
      ]);
    } finally {
      setLoading(false);
    }
  };

  const handleKey = (e) => {
    if (e.key === "Enter" && !e.shiftKey) {
      e.preventDefault();
      sendMessage();
    }
  };

  const suggestions = [
    "Why is this release failing?",
    "What exact changes should I make first?",
    "Which issue blocks deployment most?",
    "Can I deploy safely anyway?",
  ];

  return (
    <>
      {!isOpen && (
        <button
          onClick={() => setIsOpen(true)}
          className="fixed bottom-6 right-6 z-50 bg-teal-700 hover:bg-teal-800 text-white px-5 py-3 rounded-2xl shadow-xl font-semibold text-sm flex items-center gap-2 transition-all active:scale-95"
        >
          <span className="text-base">*</span>
          Ask AI about this result
          {result.risk === "HIGH" && <span className="w-2 h-2 rounded-full bg-red-400 animate-pulse" />}
        </button>
      )}

      {isOpen && (
        <div
          className="fixed bottom-6 right-6 z-50 w-96 bg-white rounded-2xl shadow-2xl border border-gray-100 flex flex-col overflow-hidden"
          style={{ height: "520px" }}
        >
          <div className="bg-teal-800 px-5 py-4 flex items-center justify-between shrink-0">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 rounded-full bg-teal-600 flex items-center justify-center text-white text-sm font-bold">
                *
              </div>
              <div>
                <p className="text-white font-bold text-sm">RMT AI Assistant</p>
                <p className="text-teal-300 text-xs">Context: {result.release_name}</p>
              </div>
            </div>
            <button
              onClick={() => setIsOpen(false)}
              className="text-teal-300 hover:text-white text-xl leading-none transition-colors"
            >
              x
            </button>
          </div>

          <div
            className={`px-4 py-2 text-xs font-semibold flex items-center gap-2 shrink-0 ${
              result.risk === "HIGH"
                ? "bg-red-50 text-red-600"
                : result.risk === "MEDIUM"
                  ? "bg-amber-50 text-amber-600"
                  : "bg-green-50 text-green-600"
            }`}
          >
            <span
              className={`w-1.5 h-1.5 rounded-full animate-pulse ${
                result.risk === "HIGH"
                  ? "bg-red-500"
                  : result.risk === "MEDIUM"
                    ? "bg-amber-500"
                    : "bg-green-500"
              }`}
            />
            Overall risk: {result.risk} | {result.issues?.length || 0} issue(s)
          </div>

          <div className="flex-1 overflow-y-auto px-4 py-4 space-y-4">
            {messages.map((msg, i) => (
              <Message key={i} msg={msg} />
            ))}

            {loading && (
              <div className="flex gap-3">
                <div className="w-7 h-7 rounded-full bg-teal-700 shrink-0 flex items-center justify-center text-xs font-bold text-white mt-0.5">
                  AI
                </div>
                <div className="bg-gray-100 px-4 py-3 rounded-2xl rounded-tl-sm flex items-center gap-1.5">
                  <span className="w-1.5 h-1.5 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "0ms" }} />
                  <span className="w-1.5 h-1.5 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "150ms" }} />
                  <span className="w-1.5 h-1.5 bg-gray-400 rounded-full animate-bounce" style={{ animationDelay: "300ms" }} />
                </div>
              </div>
            )}
            <div ref={bottomRef} />
          </div>

          {messages.length === 1 && (
            <div className="px-4 pb-2 flex flex-wrap gap-1.5 shrink-0">
              {suggestions.map((s) => (
                <button
                  key={s}
                  onClick={() => sendMessage(s)}
                  className="text-xs bg-gray-50 hover:bg-blue-50 hover:text-blue-600 border border-gray-100 hover:border-blue-100 text-gray-500 px-3 py-1.5 rounded-full transition-colors"
                >
                  {s}
                </button>
              ))}
            </div>
          )}

          <div className="px-4 py-3 border-t border-gray-100 flex gap-2 shrink-0">
            <input
              type="text"
              value={input}
              onChange={(e) => setInput(e.target.value)}
              onKeyDown={handleKey}
              placeholder="Ask about this release..."
              disabled={loading}
              className="flex-1 bg-gray-50 border border-gray-100 rounded-xl px-4 py-2.5 text-sm text-gray-700 placeholder-gray-300 focus:outline-none focus:border-blue-200 focus:bg-white transition-all"
            />
            <button
              onClick={() => sendMessage()}
              disabled={!input.trim() || loading}
              className={`px-4 py-2.5 rounded-xl font-bold text-sm transition-all active:scale-95 ${
                input.trim() && !loading
                  ? "bg-blue-600 hover:bg-blue-700 text-white"
                  : "bg-gray-100 text-gray-300 cursor-not-allowed"
              }`}
            >
              ^
            </button>
          </div>
        </div>
      )}
    </>
  );
}
