SELECT r.urgency, t.status, t.due_at, t.created_at, t.closed_at
FROM operator_tasks t
JOIN routes r ON r.case_id = t.case_id
WHERE t.created_at >= @from
  AND t.created_at < @to;
