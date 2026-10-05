UPDATE recommendations
SET mark           = @mark,
    reject_reason  = NULLIF(@reject_reason, ''),
    reject_comment = NULLIF(@reject_comment, ''),
    marked_by      = @marked_by,
    marked_at      = @marked_at,
    patient_text   = @patient_text
WHERE id = @id;
