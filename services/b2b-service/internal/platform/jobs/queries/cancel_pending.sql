UPDATE jobs
SET state = 'cancelled', locked_until = NULL
WHERE key = @key
  AND state = 'pending'
  AND kind = ANY(@kinds);
