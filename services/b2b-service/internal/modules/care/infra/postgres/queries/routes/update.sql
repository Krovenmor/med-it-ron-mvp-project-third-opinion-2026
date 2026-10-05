UPDATE routes
SET patient_id    = @patient_id,
    status        = @status,
    urgency       = NULLIF(@urgency, ''),
    modality      = @modality,
    performed_at  = @performed_at,
    received_at   = @received_at,
    assessed_at   = @assessed_at,
    review_due_at = @review_due_at,
    confirmed_at  = @confirmed_at,
    updated_at    = @updated_at
WHERE case_id = @case_id;
