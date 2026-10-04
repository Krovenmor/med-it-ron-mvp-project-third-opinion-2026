UPDATE jobs
SET state = 'done', locked_until = NULL
WHERE id = @id AND attempts = @attempts AND state = 'pending';
