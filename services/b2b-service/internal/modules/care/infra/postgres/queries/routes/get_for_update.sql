SELECT case_id, patient_id, status, urgency, modality, performed_at, received_at, assessed_at, review_due_at, confirmed_at,
       updated_at
FROM routes
WHERE case_id = @case_id
FOR UPDATE;
