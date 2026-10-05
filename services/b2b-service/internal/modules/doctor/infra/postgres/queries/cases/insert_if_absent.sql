INSERT INTO cases (id, patient_id, patient_source_system, patient_external_id, patient_birth_date, patient_sex,
                   study_id, modality, body_site, performed_at, conclusion, status, urgency, catalog_version,
                   guidelines_version, received_at, assessed_at, confirmed_at, updated_at)
VALUES (@id, @patient_id, @patient_source_system, @patient_external_id, @patient_birth_date, @patient_sex,
        @study_id, @modality, @body_site, @performed_at, @conclusion, @status, @urgency::urgency, @catalog_version,
        @guidelines_version, @received_at, @assessed_at, @confirmed_at, @updated_at)
ON CONFLICT (id) DO NOTHING;
