UPDATE plan_items
SET declined_at     = @declined_at,
    decline_reason  = @decline_reason,
    decline_comment = @decline_comment
WHERE id = @id;
