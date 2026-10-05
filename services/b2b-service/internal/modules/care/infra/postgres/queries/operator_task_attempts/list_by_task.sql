SELECT id, task_id, outcome, decline_reason, comment, callback_at, actor, created_at
FROM operator_task_attempts
WHERE task_id = @task_id
ORDER BY created_at;
