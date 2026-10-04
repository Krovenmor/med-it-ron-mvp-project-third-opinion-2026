SELECT id, patient_id, source_system, study_id, modality, body_site, performed_at, conclusion, fingerprint,
       status, urgency::text AS urgency, catalog_version, guidelines_version, received_at, updated_at
FROM cases
WHERE id = @id;
