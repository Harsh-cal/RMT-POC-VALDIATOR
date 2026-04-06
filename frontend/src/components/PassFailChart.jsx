export default function PassFailChart({ data }) {
  if (!data || data.length === 0) {
    return <p className="text-gray-400 text-sm">No data available</p>;
  }

  // Get last 14 data points
  const displayData = data.slice(-14);
  const maxCount = Math.max(...displayData.map((d) => d.passes + d.failures), 1);

  return (
    <div className="w-full">
      <div className="flex items-end justify-between h-48 gap-1">
        {displayData.map((item, idx) => {
          const passHeight = (item.passes / maxCount) * 100;
          const failHeight = (item.failures / maxCount) * 100;
          const date = new Date(item.date);
          const label = `${date.getMonth() + 1}/${date.getDate()}`;

          return (
            <div key={idx} className="flex-1 flex flex-col items-center gap-1">
              <div className="w-full flex flex-col items-stretch gap-px h-40">
                {passHeight > 0 && (
                  <div
                    className="bg-green-500"
                    style={{ height: `${passHeight}%` }}
                    title={`${item.passes} passes`}
                  />
                )}
                {failHeight > 0 && (
                  <div
                    className="bg-red-500"
                    style={{ height: `${failHeight}%` }}
                    title={`${item.failures} failures`}
                  />
                )}
              </div>
              <p className="text-xs text-gray-500">{label}</p>
            </div>
          );
        })}
      </div>

      <div className="flex items-center gap-4 mt-4 justify-center">
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded bg-green-500" />
          <span className="text-xs font-semibold text-gray-700">Passes</span>
        </div>
        <div className="flex items-center gap-2">
          <div className="w-3 h-3 rounded bg-red-500" />
          <span className="text-xs font-semibold text-gray-700">Failures</span>
        </div>
      </div>
    </div>
  );
}
