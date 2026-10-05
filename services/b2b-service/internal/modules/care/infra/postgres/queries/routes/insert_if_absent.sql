INSERT INTO routes (case_id, patient_id, status, urgency, modality, performed_at, received_at, assessed_at,
                    review_due_at, confirmed_at, updated_at)
VALUES (@case_id, @patient_id, @status, NULLIF(@urgency, ''), @modality, @performed_at, @received_at, @assessed_at,
        @review_due_at, @confirmed_at, @updated_at)
ON CONFLICT (case_id) DO NOTHING;
