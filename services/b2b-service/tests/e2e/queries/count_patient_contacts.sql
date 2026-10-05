SELECT count(*)
FROM care.patients
WHERE external_id = @external_id
  AND phone = @phone;
