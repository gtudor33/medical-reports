#!/bin/bash

# Medical Reports API Test Script
# This script demonstrates the complete API workflow

set -e

API_URL="${API_URL:-http://localhost:8080}"
DOCTOR_ID="660e8400-e29b-41d4-a716-446655440001"
HOSPITAL_ID="550e8400-e29b-41d4-a716-446655440000"

echo "================================"
echo "Medical Reports API Test Script"
echo "================================"
echo ""
echo "API URL: $API_URL"
echo ""

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print section headers
print_section() {
    echo ""
    echo -e "${GREEN}===> $1${NC}"
    echo ""
}

# Function to print test results
print_result() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}✓ $2${NC}"
    else
        echo -e "${RED}✗ $2${NC}"
    fi
}

# 1. Health Check
print_section "1. Health Check"
curl -s "$API_URL/health" | jq '.'
print_result $? "Health check"

# 2. Search ICD-10 Codes
print_section "2. Search ICD-10 Codes (pneumonie)"
curl -s "$API_URL/api/v1/reference/icd10?q=pneumonie" | jq '.'
print_result $? "ICD-10 search"

# 3. Search Medications
print_section "3. Search Medications (amox)"
curl -s "$API_URL/api/v1/reference/medications?q=amox" | jq '.'
print_result $? "Medication search"

# 4. Create a New Report
print_section "4. Create a New Report"
RESPONSE=$(curl -s -X POST "$API_URL/api/v1/reports" \
  -H "Content-Type: application/json" \
  -d "{
    \"hospital_id\": \"$HOSPITAL_ID\",
    \"patient_cnp\": \"1850312400123\",
    \"patient_first_name\": \"Ion\",
    \"patient_last_name\": \"Popescu\",
    \"specialty\": \"internal_medicine\",
    \"report_type\": \"discharge_summary\",
    \"doctor_id\": \"$DOCTOR_ID\"
  }")

echo "$RESPONSE" | jq '.'
REPORT_ID=$(echo "$RESPONSE" | jq -r '.id')
print_result $? "Create report (ID: $REPORT_ID)"

# 5. Get the Created Report
print_section "5. Get Report by ID"
curl -s "$API_URL/api/v1/reports/$REPORT_ID" | jq '.'
print_result $? "Get report"

# 6. Update Report Content
print_section "6. Update Report Content"
curl -s -X PUT "$API_URL/api/v1/reports/$REPORT_ID/content" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$DOCTOR_ID\",
    \"content\": {
      \"patient_data\": {
        \"first_name\": \"Ion\",
        \"last_name\": \"Popescu\",
        \"cnp\": \"1850312400123\",
        \"birth_date\": \"1985-03-12T00:00:00Z\",
        \"department\": \"Medicină Internă\",
        \"ward\": \"12\",
        \"bed\": \"3\",
        \"admission_date\": \"2025-10-20T08:00:00Z\",
        \"discharge_date\": \"2025-10-29T14:00:00Z\"
      },
      \"anamnesis\": {
        \"chief_complaint\": \"Dispnee și tuse productivă\",
        \"history_of_present_illness\": \"Pacient prezentat cu simptomatologie respiratorie\",
        \"past_medical_history\": \"HTA diagnosticată în 2020, diabet zaharat tip 2\",
        \"allergies\": \"Negat\",
        \"social_history\": \"Fumător 20 țigări/zi, consumă ocazional alcool\"
      },
      \"examination\": {
        \"general_condition\": \"Stare generală bună\",
        \"consciousness\": \"Lucidă\",
        \"vital_signs\": {
          \"blood_pressure\": \"130/80\",
          \"heart_rate\": 82,
          \"temperature\": 36.8,
          \"respiratory_rate\": 18,
          \"oxygen_saturation\": 96
        },
        \"systems_review\": \"Aparat respirator: MV prezente bilateral, raluri crepitante bazale stângi\"
      },
      \"diagnosis\": {
        \"primary_diagnosis\": {
          \"code\": \"J18.1\",
          \"description\": \"Pneumonie lobară nespecificată\"
        },
        \"secondary_diagnoses\": [
          {
            \"code\": \"I10\",
            \"description\": \"Hipertensiune arterială esențială\"
          },
          {
            \"code\": \"E11\",
            \"description\": \"Diabet zaharat tip 2\"
          }
        ],
        \"clinical_observations\": \"Evoluție favorabilă sub tratament antibiotic\"
      },
      \"treatment\": {
        \"medications\": [
          {
            \"name\": \"Amoxicilină + Acid Clavulanic\",
            \"dosage\": \"1g/200mg\",
            \"frequency\": \"3x/zi\",
            \"route\": \"iv\",
            \"start_date\": \"2025-10-20T08:00:00Z\",
            \"end_date\": \"2025-10-27T08:00:00Z\"
          },
          {
            \"name\": \"Paracetamol\",
            \"dosage\": \"500mg\",
            \"frequency\": \"3x/zi\",
            \"route\": \"po\",
            \"start_date\": \"2025-10-20T08:00:00Z\"
          }
        ],
        \"procedures\": []
      },
      \"recommendations\": {
        \"discharge_plan\": \"Pacient externat cu stare generală bună, apiretic\",
        \"medications\": \"Continuare Enalapril 10mg 1cp/zi, Metformin 850mg 2cp/zi\",
        \"follow_up\": \"Control medicină internă peste 2 săptămâni, radiografie toracică de control\",
        \"diet_restrictions\": \"Regim hiposodat, hipoglucidic\",
        \"activity_restrictions\": \"Repaus relativ 2 săptămâni, evitarea eforturilor intense\"
      }
    }
  }" | jq '.'
print_result $? "Update report content"

# 7. Get Report Versions
print_section "7. Get Report Versions"
curl -s "$API_URL/api/v1/reports/$REPORT_ID/versions" | jq '.'
print_result $? "Get versions"

# 8. Update Report Status to In Review
print_section "8. Update Report Status (draft → in_review)"
curl -s -X PUT "$API_URL/api/v1/reports/$REPORT_ID/status" \
  -H "Content-Type: application/json" \
  -d '{"status": "in_review"}' | jq '.'
print_result $? "Update status to in_review"

# 9. List All Reports
print_section "9. List All Reports for Doctor"
curl -s "$API_URL/api/v1/reports?doctor_id=$DOCTOR_ID&limit=10" | jq '.'
print_result $? "List reports"

# 10. Try to edit a non-draft report (should fail)
print_section "10. Try to Edit Non-Draft Report (Expected to fail)"
curl -s -X PUT "$API_URL/api/v1/reports/$REPORT_ID/content" \
  -H "Content-Type: application/json" \
  -d "{
    \"user_id\": \"$DOCTOR_ID\",
    \"content\": {}
  }" | jq '.'
echo -e "${YELLOW}This should return an error (cannot edit non-draft)${NC}"

echo ""
echo "================================"
echo -e "${GREEN}Test Script Complete!${NC}"
echo "================================"
echo ""
echo "Summary:"
echo "  - Created report: $REPORT_ID"
echo "  - Status changed: draft → in_review"
echo "  - Saved 2 versions"
echo ""
