INSERT INTO operator_tasks (id, case_id, reason, status, assignee, due_at, next_call_at, attempts, created_at,
                            updated_at, closed_at)
VALUES (@id, @case_id, @reason, @status, NULLIF(@assignee, ''), @due_at, @next_call_at, @attempts, @created_at,
        @updated_at, @closed_at);
