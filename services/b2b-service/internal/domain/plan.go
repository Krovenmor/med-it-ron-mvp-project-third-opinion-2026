package domain

type PatientPlan struct {
	Patient Patient
	Cases   []PlanCase
}

type PlanCase struct {
	Case            Case
	UrgentContact   bool
	Recommendations []Recommendation
}
