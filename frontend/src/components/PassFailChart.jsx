function formatDayLabel(dateString) {
  const date = new Date(dateString);
  return `${date.getMonth() + 1}/${date.getDate()}`;
}

export default function PassFailChart({ data }) {
  if (!data || data.length === 0) {
    return <p className="text-gray-400 text-sm">No trend data available for the selected range.</p>;
  }

  // Keep only recent days so labels stay readable.
  const displayData = data.slice(-10);
  const maxValue = Math.max(...displayData.map((d) => Math.max(d.passes || 0, d.failures || 0)), 1);

  const totalPasses = displayData.reduce((sum, item) => sum + (item.passes || 0), 0);
  const totalFailures = displayData.reduce((sum, item) => sum + (item.failures || 0), 0);

  const latest = displayData[displayData.length - 1] || { passes: 0, failures: 0 };
  const latestStatus = (latest.failures || 0) > (latest.passes || 0) ? "Needs Attention" : "Stable";

  return (
    <div className="w-full space-y-4">
      <div className="grid grid-cols-1 sm:grid-cols-3 gap-3">
        <div className="rounded-lg border border-emerald-100 bg-emerald-50 px-4 py-3">
          <p className="text-xs font-semibold text-emerald-700">Recent Passes</p>
          <p className="text-xl font-black text-emerald-800">{totalPasses}</p>
        </div>
        <div className="rounded-lg border border-red-100 bg-red-50 px-4 py-3">
          <p className="text-xs font-semibold text-red-700">Recent Failures</p>
          <p className="text-xl font-black text-red-800">{totalFailures}</p>
        </div>
        <div className="rounded-lg border border-blue-100 bg-blue-50 px-4 py-3">
          <p className="text-xs font-semibold text-blue-700">Latest Day Status</p>
          <p className="text-xl font-black text-blue-800">{latestStatus}</p>
        </div>
      </div>

      <div className="rounded-lg border border-gray-100 bg-white px-3 py-4">
        <p className="text-xs font-semibold text-gray-500 mb-3">Daily Pass vs Failure (last 10 days)</p>
        <div className="flex items-end gap-2 h-44">
          {displayData.map((item, index) => {
            const passHeight = Math.max(((item.passes || 0) / maxValue) * 130, 2);
            const failHeight = Math.max(((item.failures || 0) / maxValue) * 130, 2);

            return (
              <div key={index} className="flex-1 min-w-0">
                <div className="flex items-end justify-center gap-1 h-36">
                  <div
                    className="w-3 rounded-t bg-emerald-500"
                    style={{ height: `${passHeight}px` }}
                    title={`${item.passes || 0} passes`}
                  />
                  <div
                    className="w-3 rounded-t bg-red-500"
                    style={{ height: `${failHeight}px` }}
                    title={`${item.failures || 0} failures`}
                  />
                </div>
                <p className="text-[10px] text-gray-500 text-center mt-2">{formatDayLabel(item.date)}</p>
              </div>
            );
          })}
        </div>
      </div>

      <div className="rounded-lg border border-gray-100 bg-gray-50 px-3 py-3">
        <p className="text-xs font-semibold text-gray-600 mb-2">Quick Day-wise Summary</p>
        <div className="grid grid-cols-1 sm:grid-cols-2 gap-2 text-xs">
          {displayData.slice().reverse().slice(0, 6).map((item, index) => (
            <div key={index} className="flex items-center justify-between rounded-md bg-white px-2.5 py-2 border border-gray-100">
              <span className="font-semibold text-gray-700">{formatDayLabel(item.date)}</span>
              <span className="text-emerald-700 font-semibold">P: {item.passes || 0}</span>
              <span className="text-red-700 font-semibold">F: {item.failures || 0}</span>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
