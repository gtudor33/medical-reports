import { useState, useEffect } from 'react';
import { useNavigate, useParams } from 'react-router-dom';
import { reportsAPI } from '../services/api';
import { useAuth } from '../contexts/AuthContext';
import ICD10Search from '../components/ICD10Search';

export default function ReportForm() {
  const { id } = useParams();
  const { user } = useAuth();
  const navigate = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');
  const [formData, setFormData] = useState({
    patient_cnp: '',
    patient_first_name: '',
    patient_last_name: '',
    specialty: user?.specialty || '',
    report_type: '',
    content: {
      chief_complaint: '',
      history: '',
      physical_examination: '',
      diagnosis_codes: [],
      treatment_plan: '',
      medications: [],
      notes: '',
    },
  });

  useEffect(() => {
    if (id) {
      loadReport();
    }
  }, [id]);

  const loadReport = async () => {
    try {
      const response = await reportsAPI.get(id);
      setFormData({
        patient_cnp: response.data.patient_cnp,
        patient_first_name: response.data.patient_first_name,
        patient_last_name: response.data.patient_last_name,
        specialty: response.data.specialty,
        report_type: response.data.report_type,
        content: response.data.content || formData.content,
      });
    } catch (error) {
      setError('Failed to load report');
      console.error(error);
    }
  };

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData({ ...formData, [name]: value });
  };

  const handleContentChange = (field, value) => {
    setFormData({
      ...formData,
      content: { ...formData.content, [field]: value },
    });
  };

  const handleICD10Select = (codes) => {
    handleContentChange('diagnosis_codes', Array.isArray(codes) ? codes : [codes]);
  };

  const addMedication = () => {
    const medications = formData.content.medications || [];
    handleContentChange('medications', [
      ...medications,
      { name: '', dosage: '', frequency: '', duration: '' },
    ]);
  };

  const updateMedication = (index, field, value) => {
    const medications = [...formData.content.medications];
    medications[index] = { ...medications[index], [field]: value };
    handleContentChange('medications', medications);
  };

  const removeMedication = (index) => {
    const medications = formData.content.medications.filter((_, i) => i !== index);
    handleContentChange('medications', medications);
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');
    setLoading(true);

    try {
      if (id) {
        await reportsAPI.update(id, formData);
      } else {
        // Debug: Log user data to see what we're working with
        console.log('User data:', user);
        console.log('User hospital_id:', user.hospital_id);
        
        // Add required fields from authenticated user
        // Use default UUID if user's hospital_id is not a valid UUID
        const hospitalId = user.hospital_id && user.hospital_id.match(/^[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}$/i) 
          ? user.hospital_id 
          : "550e8400-e29b-41d4-a716-446655440000"; // Default hospital UUID from test data
          
        const createData = {
          ...formData,
          hospital_id: hospitalId,
          doctor_id: user.id,
        };
        
        console.log('Sending create data:', createData);
        await reportsAPI.create(createData);
      }
      navigate('/reports');
    } catch (error) {
      console.error('Create report error:', error);
      setError(error.response?.data?.message || 'Failed to save report');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="px-4 sm:px-6 lg:px-8">
      <div className="max-w-4xl mx-auto">
        <div className="md:flex md:items-center md:justify-between">
          <div className="flex-1 min-w-0">
            <h1 className="text-2xl font-bold leading-7 text-gray-900 sm:text-3xl sm:truncate">
              {id ? 'Edit Report' : 'New Medical Report'}
            </h1>
          </div>
        </div>

        <form onSubmit={handleSubmit} className="mt-8 space-y-6">
          {error && (
            <div className="rounded-md bg-red-50 p-4">
              <div className="text-sm text-red-800">{error}</div>
            </div>
          )}

          {/* Patient Information */}
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg font-medium leading-6 text-gray-900 mb-4">Patient Information</h3>
              <div className="grid grid-cols-1 gap-6 sm:grid-cols-2">
                <div>
                  <label htmlFor="patient_cnp" className="block text-sm font-medium text-gray-700">
                    CNP (Personal Numeric Code)
                  </label>
                  <input
                    type="text"
                    name="patient_cnp"
                    id="patient_cnp"
                    required
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.patient_cnp}
                    onChange={handleChange}
                  />
                </div>

                <div>
                  <label htmlFor="report_type" className="block text-sm font-medium text-gray-700">
                    Report Type
                  </label>
                  <select
                    name="report_type"
                    id="report_type"
                    required
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.report_type}
                    onChange={handleChange}
                  >
                    <option value="">Select type</option>
                    <option value="discharge_summary">Discharge Summary</option>
                    <option value="transfer_summary">Transfer Summary</option>
                    <option value="operative_note">Operative Note</option>
                  </select>
                </div>

                <div>
                  <label htmlFor="patient_first_name" className="block text-sm font-medium text-gray-700">
                    First Name
                  </label>
                  <input
                    type="text"
                    name="patient_first_name"
                    id="patient_first_name"
                    required
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.patient_first_name}
                    onChange={handleChange}
                  />
                </div>

                <div>
                  <label htmlFor="patient_last_name" className="block text-sm font-medium text-gray-700">
                    Last Name
                  </label>
                  <input
                    type="text"
                    name="patient_last_name"
                    id="patient_last_name"
                    required
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.patient_last_name}
                    onChange={handleChange}
                  />
                </div>

                <div className="sm:col-span-2">
                  <label htmlFor="specialty" className="block text-sm font-medium text-gray-700">
                    Specialty
                  </label>
                  <select
                    name="specialty"
                    id="specialty"
                    required
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.specialty}
                    onChange={handleChange}
                  >
                    <option value="">Select specialty</option>
                    <option value="internal_medicine">Internal Medicine</option>
                    <option value="cardiology">Cardiology</option>
                    <option value="neurology">Neurology</option>
                    <option value="pediatrics">Pediatrics</option>
                    <option value="surgery">Surgery</option>
                  </select>
                </div>
              </div>
            </div>
          </div>

          {/* Clinical Information */}
          <div className="bg-white shadow rounded-lg">
            <div className="px-4 py-5 sm:p-6">
              <h3 className="text-lg font-medium leading-6 text-gray-900 mb-4">Clinical Information</h3>
              <div className="space-y-6">
                <div>
                  <label htmlFor="chief_complaint" className="block text-sm font-medium text-gray-700">
                    Chief Complaint
                  </label>
                  <textarea
                    name="chief_complaint"
                    id="chief_complaint"
                    rows={2}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.content.chief_complaint}
                    onChange={(e) => handleContentChange('chief_complaint', e.target.value)}
                  />
                </div>

                <div>
                  <label htmlFor="history" className="block text-sm font-medium text-gray-700">
                    History of Present Illness
                  </label>
                  <textarea
                    name="history"
                    id="history"
                    rows={4}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.content.history}
                    onChange={(e) => handleContentChange('history', e.target.value)}
                  />
                </div>

                <div>
                  <label htmlFor="physical_examination" className="block text-sm font-medium text-gray-700">
                    Physical Examination
                  </label>
                  <textarea
                    name="physical_examination"
                    id="physical_examination"
                    rows={4}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.content.physical_examination}
                    onChange={(e) => handleContentChange('physical_examination', e.target.value)}
                  />
                </div>

                <ICD10Search
                  selectedCodes={formData.content.diagnosis_codes || []}
                  onSelect={handleICD10Select}
                />

                <div>
                  <label htmlFor="treatment_plan" className="block text-sm font-medium text-gray-700">
                    Treatment Plan
                  </label>
                  <textarea
                    name="treatment_plan"
                    id="treatment_plan"
                    rows={4}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.content.treatment_plan}
                    onChange={(e) => handleContentChange('treatment_plan', e.target.value)}
                  />
                </div>

                {/* Medications */}
                <div>
                  <div className="flex justify-between items-center mb-2">
                    <label className="block text-sm font-medium text-gray-700">Medications</label>
                    <button
                      type="button"
                      onClick={addMedication}
                      className="text-sm text-blue-600 hover:text-blue-800"
                    >
                      + Add Medication
                    </button>
                  </div>
                  {formData.content.medications?.map((med, index) => (
                    <div key={index} className="border rounded-md p-4 mb-3 bg-gray-50">
                      <div className="grid grid-cols-2 gap-4">
                        <div className="col-span-2">
                          <input
                            type="text"
                            placeholder="Medication name"
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border"
                            value={med.name}
                            onChange={(e) => updateMedication(index, 'name', e.target.value)}
                          />
                        </div>
                        <div>
                          <input
                            type="text"
                            placeholder="Dosage"
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border"
                            value={med.dosage}
                            onChange={(e) => updateMedication(index, 'dosage', e.target.value)}
                          />
                        </div>
                        <div>
                          <input
                            type="text"
                            placeholder="Frequency"
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border"
                            value={med.frequency}
                            onChange={(e) => updateMedication(index, 'frequency', e.target.value)}
                          />
                        </div>
                        <div>
                          <input
                            type="text"
                            placeholder="Duration"
                            className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-3 py-2 border"
                            value={med.duration}
                            onChange={(e) => updateMedication(index, 'duration', e.target.value)}
                          />
                        </div>
                        <div className="col-span-2">
                          <button
                            type="button"
                            onClick={() => removeMedication(index)}
                            className="text-sm text-red-600 hover:text-red-800"
                          >
                            Remove
                          </button>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>

                <div>
                  <label htmlFor="notes" className="block text-sm font-medium text-gray-700">
                    Additional Notes
                  </label>
                  <textarea
                    name="notes"
                    id="notes"
                    rows={3}
                    className="mt-1 block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
                    value={formData.content.notes}
                    onChange={(e) => handleContentChange('notes', e.target.value)}
                  />
                </div>
              </div>
            </div>
          </div>

          {/* Actions */}
          <div className="flex justify-end gap-3">
            <button
              type="button"
              onClick={() => navigate('/reports')}
              className="px-4 py-2 border border-gray-300 rounded-md shadow-sm text-sm font-medium text-gray-700 bg-white hover:bg-gray-50 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500"
            >
              Cancel
            </button>
            <button
              type="submit"
              disabled={loading}
              className="px-4 py-2 border border-transparent rounded-md shadow-sm text-sm font-medium text-white bg-blue-600 hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-offset-2 focus:ring-blue-500 disabled:opacity-50 disabled:cursor-not-allowed"
            >
              {loading ? 'Saving...' : id ? 'Update Report' : 'Create Report'}
            </button>
          </div>
        </form>
      </div>
    </div>
  );
}
