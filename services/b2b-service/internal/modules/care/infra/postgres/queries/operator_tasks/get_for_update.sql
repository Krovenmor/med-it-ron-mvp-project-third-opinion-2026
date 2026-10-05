SELECT id, case_id, reason, status, assignee, due_at, next_call_at, attempts, created_at, updated_at, closed_at
FROM operator_tasks
WHERE id = @id
FOR UPDATE;
