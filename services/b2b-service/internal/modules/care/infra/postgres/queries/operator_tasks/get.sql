SELECT t.id, t.case_id, t.reason, t.status, t.assignee, t.due_at, t.next_call_at, t.attempts,
       t.created_at, t.updated_at, t.closed_at
FROM operator_tasks t
WHERE t.id = @id;
