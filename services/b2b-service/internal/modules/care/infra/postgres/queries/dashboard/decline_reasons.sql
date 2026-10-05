SELECT declines.decline_reason, count(*) AS count
FROM (SELECT a.decline_reason
      FROM operator_task_attempts a
      WHERE a.outcome = 'declined'
        AND a.created_at >= @from
        AND a.created_at < @to
      UNION ALL
      SELECT i.decline_reason
      FROM plan_items i
      WHERE i.declined_at >= @from
        AND i.declined_at < @to) declines
GROUP BY declines.decline_reason
ORDER BY count DESC, declines.decline_reason;
