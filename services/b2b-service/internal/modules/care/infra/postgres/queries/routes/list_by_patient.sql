SELECT case_id, patient_id, status, urgency, modality, performed_at, received_at, assessed_at, review_due_at, confirmed_at,
       updated_at
FROM routes
WHERE patient_id = @patient_id
ORDER BY performed_at DESC, received_at DESC;
