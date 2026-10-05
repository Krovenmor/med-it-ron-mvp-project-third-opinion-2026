UPDATE jobs
SET state = 'failed', locked_until = NULL, last_error = @last_error
WHERE id = @id AND attempts = @attempts AND state = 'pending';
