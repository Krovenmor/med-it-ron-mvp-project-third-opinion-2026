package demo

import "strings"

var services = []Service{
	{Code: "THER-CONSULT", Name: "Приём терапевта", Description: "Терапевт обсудит с вами результаты обследований и дальнейший план."},
	{Code: "PULM-CONSULT", Name: "Консультация пульмонолога", Description: "Пульмонолог оценит результаты исследования лёгких и подскажет дальнейшие шаги."},
	{Code: "ONC-MAMMO-CONSULT", Name: "Консультация онколога-маммолога", Description: "Маммолог обсудит с вами результаты обследования и дальнейшие шаги."},
	{Code: "ONC-CONSULT", Name: "Консультация онколога", Description: "Специалист обсудит с вами результаты обследования и подскажет дальнейший план."},
	{Code: "SURG-THOR-CONSULT", Name: "Консультация торакального хирурга", Description: "Хирург оценит результаты исследования и подскажет, нужна ли дополнительная помощь."},
	{Code: "CARDIO-CONSULT", Name: "Консультация кардиолога", Description: "Кардиолог проверит работу сердца и при необходимости скорректирует план."},
	{Code: "US-BREAST", Name: "УЗИ молочных желез", Description: "Ультразвуковое исследование поможет врачу получить больше информации."},
	{Code: "MRI-BREAST", Name: "МРТ молочных желез с контрастированием", Description: "МРТ даст врачу более подробное изображение для точного плана."},
	{Code: "BIOPSY-BREAST", Name: "Биопсия молочной железы под контролем УЗИ", Description: "Небольшая процедура, которая поможет врачу уточнить результаты обследования."},
	{Code: "MG-SCREENING", Name: "Маммография (скрининг)", Description: "Плановое профилактическое исследование молочных желез."},
	{Code: "CT-CHEST", Name: "КТ органов грудной клетки", Description: "Компьютерная томография даст врачу подробное изображение лёгких."},
	{Code: "CT-CHEST-FOLLOWUP", Name: "Контрольная КТ органов грудной клетки", Description: "Повторное исследование покажет, как изменилась картина со временем."},
	{Code: "DX-CHEST", Name: "Рентгенография органов грудной клетки", Description: "Рентгеновский снимок поможет врачу оценить состояние лёгких."},
	{Code: "SPIROMETRY", Name: "Спирометрия", Description: "Простой тест дыхания: врач оценит, как работают лёгкие."},
	{Code: "LAB-CBC", Name: "Общий анализ крови", Description: "Анализ крови поможет врачу оценить общее состояние."},
	{Code: "LAB-BIOCHEM", Name: "Анализ крови (биохимия)", Description: "Биохимический анализ крови даст врачу дополнительную информацию."},
	{Code: "SCREENING-ANNUAL", Name: "Плановый профилактический осмотр", Description: "Ежегодный осмотр, чтобы вовремя заметить изменения."},
}

func SearchServices(query string) []Service {
	query = strings.ToLower(strings.TrimSpace(query))
	found := make([]Service, 0, len(services))
	for _, s := range services {
		if query == "" || strings.Contains(strings.ToLower(s.Name), query) || strings.Contains(strings.ToLower(s.Code), query) {
			found = append(found, s)
		}
	}
	return found
}
