SELECT kind, state
FROM (SELECT kind, state, run_at
      FROM doctor.jobs
      WHERE key = @case_id
        AND kind <> 'deliver_event'
      UNION ALL
      SELECT kind, state, run_at
      FROM care.jobs
      WHERE key = @case_id
        AND kind <> 'deliver_event') jobs
ORDER BY run_at, kind;
