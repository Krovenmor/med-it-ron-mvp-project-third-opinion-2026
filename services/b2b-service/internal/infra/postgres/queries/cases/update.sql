UPDATE cases
SET status             = @status,
    urgency            = NULLIF(@urgency, '')::urgency,
    catalog_version    = @catalog_version,
    guidelines_version = @guidelines_version,
    updated_at         = @updated_at
WHERE id = @id;
