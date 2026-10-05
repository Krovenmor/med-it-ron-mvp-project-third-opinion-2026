package demo

import "time"

const SourceSystem = "mis-demo"

var clinicZone = time.FixedZone("MSK", 3*60*60)

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

var histories = map[string][]visitSpec{
	anna.ID: {
		{"MG-SCREENING", "Маммография (скрининг)", -385, 10},
		{"THER-CONSULT", "Приём терапевта", -3, 11},
	},
	igor.ID: {
		{"THER-CONSULT", "Приём терапевта", -100, 10},
		{"PULM-CONSULT", "Консультация пульмонолога", -85, 14},
		{"THER-CONSULT", "Приём терапевта", -30, 10},
		{"LAB-CBC", "Общий анализ крови", -29, 9},
		{"THER-CONSULT", "Приём терапевта", -4, 11},
	},
	sergey.ID: {
		{"THER-CONSULT", "Приём терапевта", -1, 8},
	},
	olga.ID: {
		{"DX-CHEST", "Рентгенография органов грудной клетки", -366, 9},
		{"THER-CONSULT", "Приём терапевта", -2, 16},
	},
}

var reports = []reportSpec{
	{
		Study: studySpec{ID: "MG-2026-0001", Modality: "MG", BodySite: "breast", Day: -1},
		Conclusion: "Маммография обеих молочных желез. В верхне-наружном квадранте левой молочной железы " +
			"образование неправильной формы 12 мм с нечёткими контурами. Заключение: BI-RADS 4.",
		Patient: anna,
	},
	{
		Study: studySpec{ID: "CT-2026-0002", Modality: "CT", BodySite: "chest", Day: -1},
		Conclusion: "КТ органов грудной клетки. В S6 правого лёгкого солидный очаг 9 мм с ровными контурами. " +
			"Внутригрудные лимфоузлы не увеличены.",
		Patient: igor,
	},
	{
		Study:      studySpec{ID: "DX-2026-0003", Modality: "DX", BodySite: "chest", Day: -1},
		Conclusion: "Рентгенография органов грудной клетки. Справа пневмоторакс с коллапсом лёгкого на 1/3 объёма.",
		Patient:    sergey,
	},
	{
		Study: studySpec{ID: "DX-2026-0004", Modality: "DX", BodySite: "chest", Day: -1},
		Conclusion: "Рентгенография органов грудной клетки. Лёгочные поля без патологических изменений. " +
			"Корни структурны, синусы свободны.",
		Patient: olga,
	},
}
