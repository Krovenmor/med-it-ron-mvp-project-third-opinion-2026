SELECT phone
FROM gateway.patients
WHERE external_id = @external_id;
