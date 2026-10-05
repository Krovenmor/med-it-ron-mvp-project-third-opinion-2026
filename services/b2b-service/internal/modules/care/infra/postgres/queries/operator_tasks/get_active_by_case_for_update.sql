SELECT id, case_id, reason, status, assignee, due_at, next_call_at, attempts, created_at, updated_at, closed_at
FROM operator_tasks
WHERE case_id = @case_id
  AND status IN ('new', 'in_progress', 'no_answer', 'callback')
FOR UPDATE;
