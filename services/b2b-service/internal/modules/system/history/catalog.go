package history

var (
	doctors   = []string{"dr-demo", "dr-sokolova", "dr-orlov"}
	operators = []string{"admin-demo", "admin-2"}

	femaleNames = []string{"Анна", "Елена", "Ольга", "Мария", "Наталья", "Ирина", "Татьяна", "Светлана", "Юлия", "Екатерина"}
	maleNames   = []string{"Сергей", "Игорь", "Андрей", "Дмитрий", "Алексей", "Михаил", "Николай", "Владимир", "Павел", "Олег"}
	surnames    = []string{"Смирнов", "Петров", "Иванов", "Кузнецов", "Соколов", "Попов", "Лебедев", "Козлов", "Новиков", "Морозов", "Волков", "Алексеев"}

	modalities = []weighted[string]{
		{"CT", 40},
		{"DX", 35},
		{"MG", 25},
	}

	reviewOnTime = map[string]float64{
		"emergency": 0.86,
		"priority":  0.8,
		"planned":   0.88,
		"normal":    0.9,
	}

	bookingRate = map[string]float64{
		"priority": 0.62,
		"planned":  0.48,
	}

	modalityBookingShift = map[string]float64{
		"MG": 0.06,
		"CT": 0,
		"DX": -0.05,
	}

	emergencyOutcomes = []weighted[taskOutcome]{
		{outcomeBooked, 62},
		{outcomeContacted, 22},
		{outcomeUnreachable, 10},
		{outcomeDeclined, 6},
	}

	missedBookingOutcomes = []weighted[taskOutcome]{
		{outcomeDeclined, 50},
		{outcomeUnreachable, 30},
		{outcomeContacted, 20},
	}

	declineReasons = []weighted[string]{
		{"other_clinic", 28},
		{"expensive", 22},
		{"not_needed", 22},
		{"far", 18},
		{"other", 10},
	}

	declineComments = []string{"уезжает в отпуск", "лечится по ДМС в другой сети", "решил подождать до осени"}
	rejectComments  = []string{"уже наблюдается у специалиста", "исследование выполнено месяц назад"}
)

var (
	pulmonologist = historyService{
		Code: "PULM-CONSULT", Name: "Консультация пульмонолога", Importance: "high",
		Rationale: "Находка в лёгких требует оценки специалистом.", GuidelineRef: "КР МЗ РФ",
		PatientText: "Рекомендуем консультацию пульмонолога по результатам исследования.",
	}
	thoracicSurgeon = historyService{
		Code: "SURG-THOR-CONSULT", Name: "Консультация торакального хирурга", Importance: "high",
		Rationale: "Требуется неотложная оценка хирургом.", GuidelineRef: "КР МЗ РФ «Пневмоторакс»",
		PatientText: "Врач просит вас срочно связаться с клиникой.",
	}
	chestCTFollowUp = historyService{
		Code: "CT-CHEST-FOLLOWUP", Name: "Контрольная КТ органов грудной клетки", Importance: "low",
		Rationale: "Динамическое наблюдение.", GuidelineRef: "КР МЗ РФ",
		PatientText: "Повторное исследование покажет, как изменилась картина со временем.",
	}
	chestCT = historyService{
		Code: "CT-CHEST", Name: "КТ органов грудной клетки", Importance: "high",
		Rationale: "Уточнение изменений, выявленных на рентгенограмме.", GuidelineRef: "КР МЗ РФ",
		PatientText: "Компьютерная томография даст врачу подробное изображение лёгких.",
	}
	therapist = historyService{
		Code: "THER-CONSULT", Name: "Приём терапевта", Importance: "low",
		Rationale:   "Обсуждение результатов и плана наблюдения.",
		PatientText: "Терапевт обсудит с вами результаты обследований и дальнейший план.",
	}
	mammologist = historyService{
		Code: "ONC-MAMMO-CONSULT", Name: "Консультация онколога-маммолога", Importance: "high",
		Rationale: "Образование требует консультации специалиста.", GuidelineRef: "КР МЗ РФ «Рак молочной железы»",
		PatientText: "Рекомендуем консультацию маммолога, чтобы обсудить результаты обследования.",
	}
	breastUltrasound = historyService{
		Code: "US-BREAST", Name: "УЗИ молочных желез", Importance: "low",
		Rationale: "Дополнительная визуализация.", GuidelineRef: "КР МЗ РФ «Рак молочной железы»",
		PatientText: "Ультразвуковое исследование поможет врачу получить больше информации.",
	}
	annualScreening = historyService{
		Code: "SCREENING-ANNUAL", Name: "Плановый профилактический осмотр", Importance: "low",
		Rationale:   "Патологии не выявлено, плановый скрининг в установленный срок.",
		PatientText: "Всё в порядке. Напомним вам о следующем плановом обследовании.",
	}
	mammographyScreening = historyService{
		Code: "MG-SCREENING", Name: "Маммография (скрининг)", Importance: "low",
		Rationale:   "Патологии не выявлено, плановый скрининг в установленный срок.",
		PatientText: "Всё в порядке. Напомним вам о следующей плановой маммографии.",
	}
)

var profiles = map[string]modalityProfile{
	"CT": {
		bodySite: "chest",
		urgencies: []weighted[string]{
			{"normal", 40}, {"planned", 25}, {"priority", 30}, {"emergency", 5},
		},
		conclusions: map[string]string{
			"normal":    "Органы грудной клетки без патологических изменений.",
			"planned":   "Единичные очаги до 5 мм, рекомендовано динамическое наблюдение.",
			"priority":  "Солидный очаг в нижней доле правого лёгкого 9 мм.",
			"emergency": "Признаки пневмоторакса справа.",
		},
		services: map[string][]historyService{
			"normal":    {annualScreening},
			"planned":   {chestCTFollowUp, pulmonologist},
			"priority":  {pulmonologist, chestCTFollowUp},
			"emergency": {thoracicSurgeon},
		},
	},
	"DX": {
		bodySite: "chest",
		urgencies: []weighted[string]{
			{"normal", 50}, {"planned", 25}, {"priority", 18}, {"emergency", 7},
		},
		conclusions: map[string]string{
			"normal":    "Лёгкие без очаговых и инфильтративных изменений.",
			"planned":   "Усиление лёгочного рисунка в нижних отделах.",
			"priority":  "Округлая тень в верхней доле левого лёгкого.",
			"emergency": "Пневмоторакс слева.",
		},
		services: map[string][]historyService{
			"normal":    {annualScreening},
			"planned":   {therapist, chestCT},
			"priority":  {chestCT, pulmonologist},
			"emergency": {thoracicSurgeon},
		},
	},
	"MG": {
		bodySite: "breast",
		urgencies: []weighted[string]{
			{"normal", 50}, {"planned", 25}, {"priority", 25},
		},
		conclusions: map[string]string{
			"normal":   "Молочные железы без патологических изменений, BI-RADS 1.",
			"planned":  "Вероятно доброкачественное образование, BI-RADS 3.",
			"priority": "Образование в верхне-наружном квадранте, BI-RADS 4.",
		},
		services: map[string][]historyService{
			"normal":   {mammographyScreening},
			"planned":  {breastUltrasound, mammologist},
			"priority": {mammologist, breastUltrasound},
		},
	},
}
