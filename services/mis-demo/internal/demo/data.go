package demo

import "time"

const SourceSystem = "mis-demo"

var (
	anna = Patient{
		ID: "P-ANNA", FullName: "Анна Смирнова", BirthDate: "1974-03-12", Sex: "female",
		Phone: "+79990000001", Email: "anna@example.com",
	}
	igor = Patient{
		ID: "P-IGOR", FullName: "Игорь Петров", BirthDate: "1965-05-01", Sex: "male",
		Phone: "+79990000002", Email: "igor@example.com",
	}
	sergey = Patient{
		ID: "P-SERGEY", FullName: "Сергей Иванов", BirthDate: "1988-07-20", Sex: "male",
		Phone: "+79990000003", Email: "sergey@example.com",
	}
	olga = Patient{
		ID: "P-OLGA", FullName: "Ольга Кузнецова", BirthDate: "1981-11-02", Sex: "female",
		Phone: "+79990000004", Email: "olga@example.com",
	}
)

var histories = map[string]History{
	anna.ID: {
		Visits: []Visit{
			{ServiceCode: "MG-SCREENING", ServiceName: "Маммография (скрининг)", VisitedAt: date(2025, 9, 15)},
		},
	},
	igor.ID: {
		Visits: []Visit{
			{ServiceCode: "THER-CONSULT", ServiceName: "Приём терапевта", VisitedAt: date(2026, 8, 20)},
		},
	},
	olga.ID: {
		Visits: []Visit{
			{ServiceCode: "DX-CHEST", ServiceName: "Рентгенография органов грудной клетки", VisitedAt: date(2025, 10, 1)},
		},
	},
}

var reports = []Report{
	{
		SourceSystem: SourceSystem,
		Study:        Study{ID: "MG-2026-0001", Modality: "MG", BodySite: "breast", PerformedAt: date(2026, 10, 3)},
		Conclusion: "Маммография обеих молочных желез. В верхне-наружном квадранте левой молочной железы " +
			"образование неправильной формы 12 мм с нечёткими контурами. Заключение: BI-RADS 4.",
		Patient: anna,
	},
	{
		SourceSystem: SourceSystem,
		Study:        Study{ID: "CT-2026-0002", Modality: "CT", BodySite: "chest", PerformedAt: date(2026, 10, 3)},
		Conclusion: "КТ органов грудной клетки. В S6 правого лёгкого солидный очаг 9 мм с ровными контурами. " +
			"Внутригрудные лимфоузлы не увеличены.",
		Patient: igor,
	},
	{
		SourceSystem: SourceSystem,
		Study:        Study{ID: "DX-2026-0003", Modality: "DX", BodySite: "chest", PerformedAt: date(2026, 10, 3)},
		Conclusion:   "Рентгенография органов грудной клетки. Справа пневмоторакс с коллапсом лёгкого на 1/3 объёма.",
		Patient:      sergey,
	},
	{
		SourceSystem: SourceSystem,
		Study:        Study{ID: "DX-2026-0004", Modality: "DX", BodySite: "chest", PerformedAt: date(2026, 10, 3)},
		Conclusion: "Рентгенография органов грудной клетки. Лёгочные поля без патологических изменений. " +
			"Корни структурны, синусы свободны.",
		Patient: olga,
	},
}

func HistoryOf(patientID string) History {
	h := histories[patientID]
	return History{
		Visits:       append([]Visit{}, h.Visits...),
		Appointments: append([]Appointment{}, h.Appointments...),
	}
}

func Reports() []Report {
	return reports
}

func date(year int, month time.Month, day int) time.Time {
	return time.Date(year, month, day, 9, 0, 0, 0, time.UTC)
}
