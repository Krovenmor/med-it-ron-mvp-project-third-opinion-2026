SELECT EXISTS (
    SELECT 1
    FROM case_events
    WHERE case_id = @case_id
      AND type = @type
);
