UPDATE operator_tasks
SET status       = @status,
    assignee     = NULLIF(@assignee, ''),
    next_call_at = @next_call_at,
    attempts     = @attempts,
    updated_at   = @updated_at,
    closed_at    = @closed_at
WHERE id = @id;
