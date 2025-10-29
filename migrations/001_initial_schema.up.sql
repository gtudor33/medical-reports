-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- ============================================================================
-- Main reports table
-- ============================================================================
CREATE TABLE reports (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    hospital_id UUID NOT NULL,
    patient_cnp CHAR(13) NOT NULL,
    patient_first_name VARCHAR(100) NOT NULL,
    patient_last_name VARCHAR(100) NOT NULL,
    specialty VARCHAR(50) NOT NULL,
    report_type VARCHAR(50) NOT NULL,
    status VARCHAR(20) NOT NULL,
    
    -- Audit fields
    created_by UUID NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_modified TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    finalized_at TIMESTAMPTZ,
    
    CONSTRAINT chk_status CHECK (status IN ('draft', 'in_review', 'approved', 'signed', 'cancelled'))
);

CREATE INDEX idx_reports_hospital_patient ON reports(hospital_id, patient_cnp);
CREATE INDEX idx_reports_doctor_status ON reports(created_by, status);
CREATE INDEX idx_reports_created_at ON reports(created_at DESC);

-- ============================================================================
-- Report versions (immutable snapshots)
-- ============================================================================
CREATE TABLE report_versions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES reports(id) ON DELETE CASCADE,
    version_number INT NOT NULL,
    
    -- Content stored as JSONB
    content JSONB NOT NULL,
    
    -- Metadata
    saved_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    saved_by UUID NOT NULL,
    comment TEXT,
    
    UNIQUE(report_id, version_number)
);

CREATE INDEX idx_versions_report ON report_versions(report_id, version_number DESC);

-- ============================================================================
-- ICD-10 codes reference table
-- ============================================================================
CREATE TABLE icd10_codes (
    code VARCHAR(10) PRIMARY KEY,
    description_ro TEXT NOT NULL,
    category VARCHAR(50),
    
    search_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('romanian', description_ro)
    ) STORED
);

CREATE INDEX idx_icd10_search ON icd10_codes USING GIN(search_vector);

-- Seed ICD-10 codes
INSERT INTO icd10_codes (code, description_ro, category) VALUES
('J18.1', 'Pneumonie lobară nespecificată', 'Respiratory'),
('J18.0', 'Bronhopneumonie nespecificată', 'Respiratory'),
('J15.9', 'Pneumonie bacteriană nespecificată', 'Respiratory'),
('J18.9', 'Pneumonie nespecificată', 'Respiratory'),
('I10', 'Hipertensiune arterială esențială (primară)', 'Cardiovascular'),
('I25.1', 'Boală cardiacă ischemică aterosclerotică', 'Cardiovascular'),
('I50.0', 'Insuficiență cardiacă congestivă', 'Cardiovascular'),
('E11', 'Diabet zaharat tip 2', 'Endocrine'),
('E11.9', 'Diabet zaharat tip 2 fără complicații', 'Endocrine'),
('E78.5', 'Hiperlipidemie nespecificată', 'Endocrine'),
('K29.7', 'Gastrită cronică nespecificată', 'Digestive'),
('K21.9', 'Boală de reflux gastroesofagian fără esofagită', 'Digestive'),
('M79.3', 'Paniculită nespecificată', 'Musculoskeletal'),
('M54.5', 'Durere lombară', 'Musculoskeletal'),
('R50.9', 'Febră nespecificată', 'Symptoms'),
('R05', 'Tuse', 'Symptoms'),
('R06.0', 'Dispnee', 'Symptoms');

-- ============================================================================
-- Medications reference table
-- ============================================================================
CREATE TABLE medications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    active_substance TEXT NOT NULL,
    form VARCHAR(50),
    dosage TEXT,
    manufacturer TEXT,
    
    search_vector tsvector GENERATED ALWAYS AS (
        to_tsvector('romanian', name || ' ' || active_substance)
    ) STORED
);

CREATE INDEX idx_medications_search ON medications USING GIN(search_vector);

-- Seed medications
INSERT INTO medications (name, active_substance, form, dosage, manufacturer) VALUES
('Augmentin', 'Amoxicilină + Acid Clavulanic', 'tablet', '1g/200mg', 'GSK'),
('Amoxicilină', 'Amoxicilină', 'capsule', '500mg', 'Antibiotice Iași'),
('Paracetamol', 'Paracetamol', 'tablet', '500mg', 'Terapia'),
('Ibuprofen', 'Ibuprofen', 'tablet', '400mg', 'Terapia'),
('Metformin', 'Metformin', 'tablet', '850mg', 'Zentiva'),
('Enalapril', 'Enalapril', 'tablet', '10mg', 'Terapia'),
('Atorvastatină', 'Atorvastatină', 'tablet', '20mg', 'Pfizer'),
('Omeprazol', 'Omeprazol', 'capsule', '20mg', 'Zentiva'),
('Furosemid', 'Furosemid', 'tablet', '40mg', 'Terapia'),
('Aspirină', 'Acid Acetilsalicilic', 'tablet', '75mg', 'Bayer');

-- ============================================================================
-- Audit log
-- ============================================================================
CREATE TABLE audit_log (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    report_id UUID NOT NULL REFERENCES reports(id),
    event_type VARCHAR(50) NOT NULL,
    event_data JSONB NOT NULL,
    user_id UUID NOT NULL,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ip_address INET
);

CREATE INDEX idx_audit_report ON audit_log(report_id, timestamp DESC);
CREATE INDEX idx_audit_user ON audit_log(user_id, timestamp DESC);
