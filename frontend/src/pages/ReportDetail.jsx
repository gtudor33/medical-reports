import { useEffect, useState } from 'react';
import { useParams, useNavigate, Link } from 'react-router-dom';
import { reportsAPI } from '../services/api';

export default function ReportDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [report, setReport] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [finalizing, setFinalizing] = useState(false);

  useEffect(() => {
    loadReport();
  }, [id]);

  const loadReport = async () => {
    try {
      const response = await reportsAPI.get(id);
      setReport(response.data);
    } catch (error) {
      setError('Failed to load report');
      console.error(error);
    } finally {
      setLoading(false);
    }
  };

  const handleFinalize = async () => {
    if (!window.confirm('Are you sure you want to finalize this report? This action cannot be undone.')) {
      return;
    }

    setFinalizing(true);
    try {
      await reportsAPI.finalize(id);
      await loadReport();
    } catch (error) {
      setError('Failed to finalize report');
      console.error(error);
    } finally {
      setFinalizing(false);
    }
  };

  const handleDelete = async () => {
    if (!window.confirm('Are you sure you want to delete this report? This action cannot be undone.')) {
      return;
    }

    try {
      await reportsAPI.delete(id);
      navigate('/reports');
    } catch (error) {
      setError('Failed to delete report');
      console.error(error);
    }
  };

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
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  };

  if (loading) {
    return (
      <div className="flex justify-center items-center h-64">
        <div className="text-gray-500">Loading report...</div>
      </div>
    );
  }

  if (error && !report) {
    return (
      <div className="px-4 sm:px-6 lg:px-8">
        <div className="rounded-md bg-red-50 p-4">
          <div className="text-sm text-red-800">{error}</div>
        </div>
      </div>
    );
  }

  return (
    <div className="px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        {/* Header */}
        <div className="md:flex md:items-center md:justify-between">
          <div className="flex-1 min-w-0">
            <h1 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
              Medical Report
            </h1>
            <div className="mt-2 flex items-center text-sm text-gray-500">
              <span className={`inline-flex rounded-full px-3 py-1 text-xs font-semibold ${getStatusBadge(report.status)}`}>
                {report.status.replace('_', ' ').toUpperCase()}
              </span>
              <span className="ml-4">Created: {formatDate(report.created_at)}</span>
              <span className="ml-4">Last Modified: {formatDate(report.last_modified)}</span>
            </div>
          </div>
          <div className="mt-4 flex md:mt-0 md:ml-4 gap-2">
            {report.status === 'draft' && (
              <>
                <Link
                  to={`/reports/${id}/edit`}
                  className="inline-flex items-center px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
                >
                  Edit
                </Link>
                <button
                  onClick={handleFinalize}
                  disabled={finalizing}
                  className="inline-flex items-center px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-green-600 hover:bg-green-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-green-500 disabled:opacity-50"
                >
                  {finalizing ? 'Finalizing...' : 'Finalize Report'}
                </button>
              </>
            )}
            {report.status === 'draft' && (
              <button
                onClick={handleDelete}
                className="inline-flex items-center px-4 py-2 border border-red-300 rounded-md shadow-sm text-sm font-medium text-red-700 bg-white hover:bg-red-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-red-500"
              >
                Delete
              </button>
            )}
          </div>
        </div>

        {error && (
          <div className="mt-4 rounded-md bg-red-50 p-4">
            <div className="text-sm text-red-800">{error}</div>
          </div>
        )}

        {/* Patient Information */}
        <div className="mt-6 bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium leading-6 text-gray-900 mb-4">Patient Information</h3>
            <dl className="grid grid-cols-1 gap-x-4 gap-y-4 sm:grid-cols-2">
              <div>
                <dt className="text-sm font-medium text-gray-500">Full Name</dt>
                <dd className="mt-1 text-sm text-gray-900">
                  {report.patient_first_name} {report.patient_last_name}
                </dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">CNP</dt>
                <dd className="mt-1 text-sm text-gray-900">{report.patient_cnp}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Report Type</dt>
                <dd className="mt-1 text-sm text-gray-900 capitalize">{report.report_type}</dd>
              </div>
              <div>
                <dt className="text-sm font-medium text-gray-500">Specialty</dt>
                <dd className="mt-1 text-sm text-gray-900 capitalize">{report.specialty}</dd>
              </div>
            </dl>
          </div>
        </div>

        {/* Clinical Information */}
        <div className="mt-6 bg-white shadow rounded-lg">
          <div className="px-4 py-5 sm:p-6">
            <h3 className="text-lg font-medium leading-6 text-gray-900 mb-4">Clinical Information</h3>
            <div className="space-y-6">
              {report.content?.chief_complaint && (
                <div>
                  <dt className="text-sm font-medium text-gray-500 mb-1">Chief Complaint</dt>
                  <dd className="text-sm text-gray-900 whitespace-pre-wrap">{report.content.chief_complaint}</dd>
                </div>
              )}

              {report.content?.history && (
                <div>
                  <dt className="text-sm font-medium text-gray-500 mb-1">History of Present Illness</dt>
                  <dd className="text-sm text-gray-900 whitespace-pre-wrap">{report.content.history}</dd>
                </div>
              )}

              {report.content?.physical_examination && (
                <div>
                  <dt className="text-sm font-medium text-gray-500 mb-1">Physical Examination</dt>
                  <dd className="text-sm text-gray-900 whitespace-pre-wrap">{report.content.physical_examination}</dd>
                </div>
              )}

              {report.content?.diagnosis_codes && report.content.diagnosis_codes.length > 0 && (
                <div>
                  <dt className="text-sm font-medium text-gray-500 mb-2">ICD-10 Diagnosis Codes</dt>
                  <dd className="space-y-2">
                    {report.content.diagnosis_codes.map((code) => (
                      <div key={code.code} className="flex items-start bg-blue-50 rounded-lg p-3">
                        <span className="inline-flex items-center px-2.5 py-0.5 rounded-md text-sm font-medium bg-blue-100 text-blue-800">
                          {code.code}
                        </span>
                        <span className="ml-3 text-sm text-gray-900">{code.description}</span>
                      </div>
                    ))}
                  </dd>
                </div>
              )}

              {report.content?.treatment_plan && (
                <div>
                  <dt className="text-sm font-medium text-gray-500 mb-1">Treatment Plan</dt>
                  <dd className="text-sm text-gray-900 whitespace-pre-wrap">{report.content.treatment_plan}</dd>
                </div>
              )}

              {report.content?.medications && report.content.medications.length > 0 && (
                <div>
                  <dt className="text-sm font-medium text-gray-500 mb-2">Medications</dt>
                  <dd className="space-y-2">
                    {report.content.medications.map((med, index) => (
                      <div key={index} className="bg-gray-50 rounded-lg p-3">
                        <div className="font-medium text-gray-900">{med.name}</div>
                        <div className="mt-1 text-sm text-gray-600">
                          {med.dosage && <span>Dosage: {med.dosage}</span>}
                          {med.frequency && <span className="ml-4">Frequency: {med.frequency}</span>}
                          {med.duration && <span className="ml-4">Duration: {med.duration}</span>}
                        </div>
                      </div>
                    ))}
                  </dd>
                </div>
              )}

              {report.content?.notes && (
                <div>
                  <dt className="text-sm font-medium text-gray-500 mb-1">Additional Notes</dt>
                  <dd className="text-sm text-gray-900 whitespace-pre-wrap">{report.content.notes}</dd>
                </div>
              )}
            </div>
          </div>
        </div>

        {/* Back Button */}
        <div className="mt-6 mb-8">
          <Link
            to="/reports"
            className="inline-flex items-center text-sm font-medium text-blue-600 hover:text-blue-500"
          >
            ← Back to Reports
          </Link>
        </div>
      </div>
    </div>
  );
}
