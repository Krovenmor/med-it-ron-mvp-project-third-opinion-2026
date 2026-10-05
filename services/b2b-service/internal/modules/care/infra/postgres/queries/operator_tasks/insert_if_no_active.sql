INSERT INTO operator_tasks (case_id, reason, status, due_at, created_at, updated_at)
VALUES (@case_id, @reason, @status, @due_at, @created_at, @created_at)
ON CONFLICT (case_id) WHERE status IN ('new', 'in_progress', 'no_answer', 'callback') DO NOTHING
RETURNING id, case_id, reason, status, assignee, due_at, next_call_at, attempts, created_at, updated_at, closed_at;
