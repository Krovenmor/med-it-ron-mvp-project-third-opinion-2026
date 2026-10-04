-- run_at is compared with the business clock (@now), while the lease uses the database clock,
-- so advancing demo time never releases a job that is still being processed.
UPDATE jobs
SET locked_until = now() + make_interval(secs => @lease_seconds),
    attempts     = attempts + 1
WHERE id = (
    SELECT id
    FROM jobs
    WHERE state = 'pending'
      AND run_at <= @now
      AND (locked_until IS NULL OR locked_until < now())
    ORDER BY run_at
    LIMIT 1
    FOR UPDATE SKIP LOCKED
)
RETURNING id, case_id, kind, run_at, attempts;
