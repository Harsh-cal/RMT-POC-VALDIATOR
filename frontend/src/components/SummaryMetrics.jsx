export default function SummaryMetrics({ metrics }) {
  const formatPercent = (val) => (val * 100).toFixed(1);

  return (
    <div className="grid grid-cols-2 sm:grid-cols-4 gap-4">
      <div className="bg-green-50 border border-green-100 rounded-xl p-4">
        <p className="text-xs font-bold text-green-600 uppercase">Pass Rate</p>
        <p className="text-2xl font-black text-green-700 mt-2">
          {formatPercent(metrics.pass_rate)}%
        </p>
        <p className="text-xs text-green-600 mt-1">Validation Success</p>
      </div>

      <div className="bg-blue-50 border border-blue-100 rounded-xl p-4">
        <p className="text-xs font-bold text-blue-600 uppercase">Avg Issues</p>
        <p className="text-2xl font-black text-blue-700 mt-2">
          {metrics.avg_issues?.toFixed(1) || 0}
        </p>
        <p className="text-xs text-blue-600 mt-1">Per Release</p>
      </div>

      <div className="bg-red-50 border border-red-100 rounded-xl p-4">
        <p className="text-xs font-bold text-red-600 uppercase">High Risk</p>
        <p className="text-2xl font-black text-red-700 mt-2">{metrics.high_risk_count || 0}</p>
        <p className="text-xs text-red-600 mt-1">Releases</p>
      </div>

      <div className="bg-yellow-50 border border-yellow-100 rounded-xl p-4">
        <p className="text-xs font-bold text-yellow-600 uppercase">Recurring</p>
        <p className="text-2xl font-black text-yellow-700 mt-2">{metrics.recurring_issues || 0}</p>
        <p className="text-xs text-yellow-600 mt-1">Issue Types</p>
      </div>
    </div>
  );
}
