INSERT INTO cases (patient_id, source_system, study_id, modality, body_site, performed_at, conclusion,
                   fingerprint, status, received_at, updated_at)
VALUES (@patient_id, @source_system, @study_id, @modality, @body_site, @performed_at, @conclusion,
        @fingerprint, @status, @received_at, @updated_at)
ON CONFLICT (source_system, study_id) DO NOTHING
RETURNING id;
