INSERT INTO intakes (case_id, patient_id, source_system, study_id, modality, body_site, performed_at, conclusion,
                     fingerprint, status, received_at, assessed_at)
VALUES (@case_id, @patient_id, @source_system, @study_id, @modality, @body_site, @performed_at, @conclusion,
        @fingerprint, @status, @received_at, @assessed_at);
