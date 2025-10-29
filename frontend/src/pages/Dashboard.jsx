import { useEffect, useState } from 'react';
import { Link } from 'react-router-dom';
import { reportsAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';

export default function Dashboard() {
  const { user } = useAuth();
  const [stats, setStats] = useState({
    draft: 0,
    inReview: 0,
    finalized: 0,
    total: 0,
  });
  const [recentReports, setRecentReports] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    loadDashboardData();
  }, []);

  const loadDashboardData = async () => {
    try {
      const [draftRes, reviewRes, finalizedRes] = await Promise.all([
        reportsAPI.list({ status: 'draft', limit: 100 }),
        reportsAPI.list({ status: 'in_review', limit: 100 }),
        reportsAPI.list({ status: 'finalized', limit: 100 }),
      ]);

      setStats({
        draft: draftRes.data.length,
        inReview: reviewRes.data.length,
        finalized: finalizedRes.data.length,
        total: draftRes.data.length + reviewRes.data.length + finalizedRes.data.length,
      });

      // Get recent reports
      const recentRes = await reportsAPI.list({ limit: 5 });
      setRecentReports(recentRes.data);
    } catch (error) {
      console.error('Failed to load dashboard data:', error);
    } finally {
      setLoading(false);
    }
  };

  const statCards = [
    { name: 'Total Reports', value: stats.total, color: 'bg-blue-500', icon: '📊' },
    { name: 'Draft', value: stats.draft, color: 'bg-yellow-500', icon: '✏️' },
    { name: 'In Review', value: stats.inReview, color: 'bg-orange-500', icon: '👁️' },
    { name: 'Finalized', value: stats.finalized, color: 'bg-green-500', icon: '✅' },
  ];

  const getStatusBadge = (status) => {
    const styles = {
      draft: 'bg-yellow-100 text-yellow-800',
      in_review: 'bg-orange-100 text-orange-800',
      finalized: 'bg-green-100 text-green-800',
    };
    return styles[status] || 'bg-gray-100 text-gray-800';
  };

  const formatDate = (dateString) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-gray-500">Loading...</div>
      </div>
    );
  }

  return (
    <div className="px-4 sm:px-6 lg:px-8">
      <div className="sm:flex sm:items-center">
        <div className="sm:flex-auto">
          <h1 className="text-2xl font-semibold text-gray-900">Welcome back, Dr. {user?.lastName}</h1>
          <p className="mt-2 text-sm text-gray-700">
            Here's an overview of your medical reports
          </p>
        </div>
        <div className="mt-4 sm:mt-0 sm:ml-16 sm:flex-none">
          <Link
            to="/reports/new"
            className="inline-flex items-center justify-center rounded-md border border-transparent bg-blue-600 px-4 py-2 text-sm font-medium text-white shadow-sm hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2"
          >
            ➕ New Report
          </Link>
        </div>
      </div>

      {/* Stats */}
      <div className="mt-8 grid grid-cols-1 gap-5 sm:grid-cols-2 lg:grid-cols-4">
        {statCards.map((stat) => (
          <div key={stat.name} className="bg-white overflow-hidden shadow rounded-lg">
            <div className="p-5">
              <div className="flex items-center">
                <div className={`flex-shrink-0 rounded-md p-3 ${stat.color}`}>
                  <span className="text-2xl">{stat.icon}</span>
                </div>
                <div className="ml-5 w-0 flex-1">
                  <dl>
                    <dt className="text-sm font-medium text-gray-500 truncate">{stat.name}</dt>
                    <dd className="text-3xl font-semibold text-gray-900">{stat.value}</dd>
                  </dl>
                </div>
              </div>
            </div>
          </div>
        ))}
      </div>

      {/* Recent Reports */}
      <div className="mt-8">
        <div className="bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg leading-6 font-medium text-gray-900">Recent Reports</h3>
            <div className="mt-5">
              {recentReports.length === 0 ? (
                <p className="text-gray-500 text-center py-8">No reports yet. Create your first report!</p>
              ) : (
                <div className="overflow-hidden shadow ring-1 ring-black ring-opacity-5 rounded-lg">
                  <table className="min-w-full divide-y divide-gray-300">
                    <thead className="bg-gray-50">
                      <tr>
                        <th className="py-3.5 pl-4 pr-3 text-left text-sm font-semibold text-gray-900">Patient</th>
                        <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Type</th>
                        <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Status</th>
                        <th className="px-3 py-3.5 text-left text-sm font-semibold text-gray-900">Last Modified</th>
                        <th className="relative py-3.5 pl-3 pr-4"></th>
                      </tr>
                    </thead>
                    <tbody className="divide-y divide-gray-200 bg-white">
                      {recentReports.map((report) => (
                        <tr key={report.id}>
                          <td className="whitespace-nowrap py-4 pl-4 pr-3 text-sm font-medium text-gray-900">
                            {report.patient_first_name} {report.patient_last_name}
                          </td>
                          <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">{report.report_type}</td>
                          <td className="whitespace-nowrap px-3 py-4 text-sm">
                            <span className={`inline-flex rounded-full px-2 text-xs font-semibold leading-5 ${getStatusBadge(report.status)}`}>
                              {report.status}
                            </span>
                          </td>
                          <td className="whitespace-nowrap px-3 py-4 text-sm text-gray-500">
                            {formatDate(report.last_modified)}
                          </td>
                          <td className="relative whitespace-nowrap py-4 pl-3 pr-4 text-right text-sm font-medium">
                            <Link to={`/reports/${report.id}`} className="text-blue-600 hover:text-blue-900">
                              View
                            </Link>
                          </td>
                        </tr>
                      ))}
                    </tbody>
                  </table>
                </div>
              )}
            </div>
            <div className="mt-4 text-center">
              <Link to="/reports" className="text-sm font-medium text-blue-600 hover:text-blue-500">
                View all reports →
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
