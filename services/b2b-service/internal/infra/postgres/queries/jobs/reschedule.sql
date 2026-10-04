UPDATE jobs
SET run_at = @run_at, locked_until = NULL, last_error = @last_error
WHERE id = @id AND attempts = @attempts AND state = 'pending';
