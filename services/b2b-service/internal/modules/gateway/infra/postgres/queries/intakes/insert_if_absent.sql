INSERT INTO intakes (patient_id, source_system, study_id, modality, body_site, performed_at, conclusion, fingerprint,
                     status, received_at)
VALUES (@patient_id, @source_system, @study_id, @modality, @body_site, @performed_at, @conclusion, @fingerprint,
        @status, @received_at)
ON CONFLICT (source_system, study_id) DO NOTHING
RETURNING case_id;
