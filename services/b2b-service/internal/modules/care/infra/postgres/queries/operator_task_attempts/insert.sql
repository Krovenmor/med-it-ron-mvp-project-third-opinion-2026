INSERT INTO operator_task_attempts (task_id, outcome, decline_reason, comment, callback_at, actor, created_at)
VALUES (@task_id, @outcome, NULLIF(@decline_reason, ''), @comment, @callback_at, @actor, @created_at);
