UPDATE intakes
SET status      = @status,
    assessed_at = @assessed_at
WHERE case_id = @case_id;
