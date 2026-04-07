import { useState, useEffect } from "react";
import {
  getReleasesHistory,
  getTrendsData,
  getRecurringIssues,
} from "../services/historyApi";
import SummaryMetrics from "./SummaryMetrics";
import PassFailChart from "./PassFailChart";
import TimelineView from "./TimelineView";
import IssueRecurrenceChart from "./IssueRecurrenceChart";

export default function ReleaseHistoryDashboard() {
  const [releases, setReleases] = useState([]);
  const [trends, setTrends] = useState([]);
  const [issues, setIssues] = useState([]);
  const [loading, setLoading] = useState(true);
  const [metrics, setMetrics] = useState(null);
  const [filters, setFilters] = useState({ days: 90 });

  useEffect(() => {
    loadData();
  }, [filters]);

  const loadData = async () => {
    setLoading(true);
    try {
      const [historyRes, trendsRes, issuesRes] = await Promise.all([
        getReleasesHistory({ limit: 50 }),
        getTrendsData({ days: filters.days }),
        getRecurringIssues({ days: filters.days }),
      ]);

      setReleases(historyRes.releases || []);
      setMetrics(historyRes.metrics || {});
      setTrends(trendsRes.data || []);
      setIssues(issuesRes.issues || []);
    } catch (err) {
      console.error("Failed to load history:", err);
    } finally {
      setLoading(false);
    }
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-gray-50 flex items-center justify-center">
        <div className="text-center">
          <div className="w-12 h-12 border-4 border-blue-200 border-t-blue-600 rounded-full animate-spin mx-auto mb-4" />
          <p className="text-gray-600 font-semibold">Loading Release History...</p>
        </div>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-gray-50 p-6">
      <div className="max-w-7xl mx-auto space-y-6">
        {/* Header */}
        <div className="flex items-center justify-between">
          <div>
            <h1 className="text-3xl font-black text-gray-900">📊 Release History</h1>
            <p className="text-gray-500 text-sm mt-1">
              Day-wise validation health with readable trend lines and issue patterns
            </p>
          </div>

          {/* Days Filter */}
          <div className="flex items-center gap-2">
            <label className="text-sm font-semibold text-gray-700">Last:</label>
            <select
              value={filters.days}
              onChange={(e) => setFilters({ days: parseInt(e.target.value) })}
              className="px-3 py-2 border border-gray-300 rounded-lg text-sm font-medium focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              <option value={30}>30 Days</option>
              <option value={90}>90 Days</option>
              <option value={180}>6 Months</option>
              <option value={365}>1 Year</option>
            </select>
          </div>
        </div>

        {/* Summary Metrics */}
        {metrics && <SummaryMetrics metrics={metrics} />}

        {/* Charts Row */}
        <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
          {/* Pass/Fail Trend */}
          <div className="bg-white border border-gray-100 rounded-xl shadow-sm p-6">
            <h2 className="text-lg font-bold text-gray-900 mb-1">Daily Pass vs Fail</h2>
            <p className="text-xs text-gray-500 mb-4">
              Simple day-wise view of validations, with clear pass and failure counts.
            </p>
            <PassFailChart data={trends} />
          </div>

          {/* Recurring Issues */}
          <div className="bg-white border border-gray-100 rounded-xl shadow-sm p-6">
            <h2 className="text-lg font-bold text-gray-900 mb-4">Recurring Issues</h2>
            <IssueRecurrenceChart data={issues} />
          </div>
        </div>

        {/* Recent Releases */}
        <div className="bg-white border border-gray-100 rounded-xl shadow-sm p-6">
          <h2 className="text-lg font-bold text-gray-900 mb-4">Recent Releases</h2>
          <TimelineView releases={releases} />
        </div>
      </div>
    </div>
  );
}
