import pytest

from app.domain import rules
from conftest import make_request, parse_request

# Поле aiResult.conclusion из примеров Kafka-сообщений (формат БФТ НПКЦ ДиТ, табл. 3.1).
KAFKA_DX_HYDROTHORAX_CARDIOMEGALY = (
    "Обнаружены рентгенологические признаки, которые могут коррелировать с:\n"
    "\t\tГидротораксом справа\n\t\tКардиомегалией\n"
    "Заключение подготовлено медицинским изделием на основе технологий ИИ."
)
KAFKA_CT_NODULE_EMPHYSEMA_CALCIUM = (
    "Очаговое образование верхней доли правого легкого (11 мм), дифференциальный ряд: "
    "периферическое новообразование, гранулема.\n"
    "Эмфизема легких (8% общего объема).\n"
    "Кальциноз коронарных артерий, CAC-DRS 2.\n"
    "Атеросклероз грудной аорты."
)
KAFKA_MG_NORMA = "Правая молочная железа: BI-RADS 1\nЛевая молочная железа: BI-RADS 1"

# Заключения демо-пациентов из mis-demo (services/mis-demo/internal/demo/data.go).
ANNA = (
    "Маммография обеих молочных желез. В верхне-наружном квадранте левой молочной железы "
    "образование неправильной формы 12 мм с нечёткими контурами. Заключение: BI-RADS 4."
)
IGOR = (
    "КТ органов грудной клетки. В S6 правого лёгкого солидный очаг 9 мм с ровными контурами. "
    "Внутригрудные лимфоузлы не увеличены."
)
SERGEY = "Рентгенография органов грудной клетки. Справа пневмоторакс с коллапсом лёгкого на 1/3 объёма."
OLGA = (
    "Рентгенография органов грудной клетки. Лёгочные поля без патологических изменений. "
    "Корни структурны, синусы свободны."
)


def extract(conclusion: str, modality: str = "DX") -> rules.Markers:
    return rules.extract(parse_request(make_request(conclusion, modality)))


def test_kafka_xray_hydrothorax_cardiomegaly():
    m = extract(KAFKA_DX_HYDROTHORAX_CARDIOMEGALY)
    assert {"pleural_effusion", "cardiomegaly"} <= set(m.findings)
    assert "pneumothorax" not in m.findings
    assert m.urgency_floor == "planned"


def test_kafka_ct_nodule_emphysema_calcium():
    m = extract(KAFKA_CT_NODULE_EMPHYSEMA_CALCIUM, "CT")
    assert {"lung_nodule", "emphysema", "aortic_atherosclerosis"} <= set(m.findings)
    assert m.max_nodule_mm == 11
    assert m.cac_drs == 2
    assert m.urgency_floor == "priority"


def test_kafka_mammography_norma():
    m = extract(KAFKA_MG_NORMA, "MG")
    assert m.birads == [1, 1]
    assert m.urgency_floor == "normal"


@pytest.mark.parametrize(
    ("conclusion", "modality", "floor"),
    [(ANNA, "MG", "priority"), (IGOR, "CT", "priority"), (SERGEY, "DX", "emergency"), (OLGA, "DX", "normal")],
)
def test_demo_patients(conclusion, modality, floor):
    assert extract(conclusion, modality).urgency_floor == floor


def test_igor_nodule_size():
    assert extract(IGOR, "CT").max_nodule_mm == 9


@pytest.mark.parametrize(
    "conclusion",
    [
        "Без признаков пневмоторакса.",
        "Пневмоторакс не выявлен.",
        "Свободного газа нет. Пневмоторакс: не определяется.",
    ],
)
def test_negated_pneumothorax_is_ignored(conclusion):
    m = extract(conclusion)
    assert "pneumothorax" not in m.findings
    assert m.urgency_floor == "normal"


def test_birads_in_different_styles():
    assert extract("Заключение: BI-RADS 4.", "MG").birads == [4]
    assert extract("BIRADS: 5", "MG").birads == [5]
    assert extract("категория по BI RADS 3", "MG").urgency_floor == "planned"
