SELECT case_id, patient_id, source_system, study_id, modality, body_site, performed_at, conclusion, fingerprint,
       status, received_at, assessed_at
FROM intakes
WHERE source_system = @source_system
  AND study_id = @study_id;
