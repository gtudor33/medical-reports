import { useState, useEffect, useRef } from 'react';
import { referenceAPI } from '../services/api';

export default function ICD10Search({ onSelect, selectedCodes = [] }) {
  const [query, setQuery] = useState('');
  const [results, setResults] = useState([]);
  const [loading, setLoading] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);
  const dropdownRef = useRef(null);

  useEffect(() => {
    const handleClickOutside = (event) => {
      if (dropdownRef.current && !dropdownRef.current.contains(event.target)) {
        setShowDropdown(false);
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, []);

  useEffect(() => {
    const searchICD10 = async () => {
      if (query.length < 2) {
        setResults([]);
        return;
      }

      setLoading(true);
      try {
        const response = await referenceAPI.searchICD10(query);
        setResults(response.data.codes || []);
        setShowDropdown(true);
      } catch (error) {
        console.error('Failed to search ICD-10 codes:', error);
        setResults([]);
      } finally {
        setLoading(false);
      }
    };

    const debounceTimer = setTimeout(searchICD10, 300);
    return () => clearTimeout(debounceTimer);
  }, [query]);

  const handleSelect = (code) => {
    onSelect(code);
    setQuery('');
    setShowDropdown(false);
    setResults([]);
  };

  const handleRemove = (codeToRemove) => {
    onSelect(selectedCodes.filter((c) => c.code !== codeToRemove.code));
  };

  return (
    <div className="space-y-3">
      <div className="relative" ref={dropdownRef}>
        <label className="block text-sm font-medium text-gray-700 mb-1">
          ICD-10 Diagnosis Codes
        </label>
        <input
          type="text"
          className="block w-full rounded-md border-gray-300 shadow-sm focus:border-blue-500 focus:ring-blue-500 sm:text-sm px-4 py-2 border"
          placeholder="Search ICD-10 codes (e.g., pneumonia, J18)..."
          value={query}
          onChange={(e) => setQuery(e.target.value)}
          onFocus={() => results.length > 0 && setShowDropdown(true)}
        />

        {showDropdown && (
          <div className="absolute z-10 mt-1 w-full bg-white shadow-lg max-h-60 rounded-md py-1 text-base ring-1 ring-black ring-opacity-5 overflow-auto focus:outline-none sm:text-sm">
            {loading ? (
              <div className="px-4 py-2 text-gray-500">Searching...</div>
            ) : results.length === 0 ? (
              <div className="px-4 py-2 text-gray-500">No results found</div>
            ) : (
              results.map((code) => {
                const isSelected = selectedCodes.some((c) => c.code === code.code);
                return (
                  <button
                    key={code.code}
                    type="button"
                    className={`w-full text-left px-4 py-2 hover:bg-gray-100 ${
                      isSelected ? 'bg-blue-50' : ''
                    }`}
                    onClick={() => !isSelected && handleSelect(code)}
                    disabled={isSelected}
                  >
                    <div className="flex justify-between">
                      <span className="font-medium text-gray-900">{code.code}</span>
                      {isSelected && (
                        <span className="text-blue-600 text-xs">✓ Selected</span>
                      )}
                    </div>
                    <div className="text-sm text-gray-600">{code.description}</div>
                  </button>
                );
              })
            )}
          </div>
        )}
      </div>

      {/* Selected codes */}
      {selectedCodes.length > 0 && (
        <div className="space-y-2">
          <div className="text-sm font-medium text-gray-700">Selected Codes:</div>
          <div className="flex flex-wrap gap-2">
            {selectedCodes.map((code) => (
              <div
                key={code.code}
                className="inline-flex items-center gap-2 bg-blue-50 text-blue-700 px-3 py-1 rounded-full text-sm"
              >
                <span className="font-medium">{code.code}</span>
                <span className="text-blue-600">-</span>
                <span>{code.description}</span>
                <button
                  type="button"
                  onClick={() => handleRemove(code)}
                  className="ml-1 text-blue-600 hover:text-blue-800"
                >
                  ×
                </button>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  );
}
