UPDATE cases
SET status       = @status,
    urgency      = @urgency::urgency,
    confirmed_at = @confirmed_at,
    updated_at   = @updated_at
WHERE id = @id;
