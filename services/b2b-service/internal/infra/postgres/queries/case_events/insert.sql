INSERT INTO case_events (case_id, type, from_status, to_status, actor, occurred_at)
VALUES (@case_id, @type, NULLIF(@from_status, ''), NULLIF(@to_status, ''), @actor, @occurred_at);
