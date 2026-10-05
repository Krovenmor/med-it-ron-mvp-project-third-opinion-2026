"""Детерминированный разбор заключения: маркеры и нижняя граница срочности.

Маркеры передаются LLM как подсказка, а нижняя граница срочности страхует от самой
опасной ошибки — недооценки срочности. LLM может поднять срочность, но не опустить
её ниже границы, которую выставили правила.

ВНИМАНИЕ: пороги ниже — черновик, их должен проверить медицинский эксперт.
"""

import re
from dataclasses import dataclass, field

from app.domain.assessment import URGENCY_ORDER, AssessmentRequest, Urgency

_BIRADS = re.compile(r"BI[-\s]?RADS\D{0,3}([0-6])", re.IGNORECASE)
_CAC_DRS = re.compile(r"CAC[-\s]?DRS\D{0,3}([0-3])", re.IGNORECASE)
_SIZE_MM = re.compile(
    r"(\d+(?:[.,]\d+)?)\s*(?:[×xх*]\s*\d+(?:[.,]\d+)?\s*)*мм",
    re.IGNORECASE,
)
_CLAUSE_SPLIT = re.compile(r"[.;\n]")
_NEGATION_BEFORE = re.compile(r"(без|нет|отсутств\w*|исключ\w*)\s+(\w+\s+){0,2}$")
_NEGATION_AFTER = re.compile(r"^\W*(\w+\W+){0,3}(не\s+выявл|не\s+определя|отсутств)")

# Находки: ключ маркера → основы слов, по которым ищем упоминание.
_FINDINGS: dict[str, tuple[str, ...]] = {
    "pneumothorax": ("пневмоторакс",),
    "pleural_effusion": ("гидроторакс", "плевральн\\w* выпот", "выпот\\w* в плевр"),
    "cardiomegaly": ("кардиомегал", "тень сердца\\W+расширен"),
    "lung_nodule": ("очаг",),
    "emphysema": ("эмфизем",),
    "coronary_calcium": ("кальциноз коронарн", "коронарн\\w* кальци"),
    "aortic_atherosclerosis": ("атеросклероз\\w* (грудн\\w* )?аорт",),
    "breast_mass": ("образовани\\w*", "узлов\\w* образ"),
}


@dataclass
class Markers:
    findings: list[str] = field(default_factory=list)
    birads: list[int] = field(default_factory=list)
    cac_drs: int | None = None
    max_nodule_mm: float | None = None
    urgency_floor: Urgency = "normal"
    reasons: list[str] = field(default_factory=list)

    def as_hint(self) -> dict:
        """Подсказка для LLM: только непустые поля."""
        hint = {
            "findings": self.findings,
            "birads": self.birads,
            "cac_drs": self.cac_drs,
            "max_nodule_mm": self.max_nodule_mm,
        }
        return {k: v for k, v in hint.items() if v not in (None, [])}


def extract(req: AssessmentRequest) -> Markers:
    text = req.conclusion
    m = Markers()

    for key, stems in _FINDINGS.items():
        if key == "breast_mass" and req.study.modality != "MG":
            continue
        if any(_mentioned(text, stem) for stem in stems):
            m.findings.append(key)

    m.birads = [int(x) for x in _BIRADS.findall(text)]
    cac = _CAC_DRS.findall(text)
    m.cac_drs = max(int(x) for x in cac) if cac else None
    if "lung_nodule" in m.findings:
        m.max_nodule_mm = _max_size_near(text, "очаг")

    _apply_floor(m)
    return m


def _mentioned(text: str, stem: str) -> bool:
    """Есть ли утвердительное упоминание находки (грубая проверка отрицаний)."""
    for match in re.finditer(stem, text, re.IGNORECASE):
        before = text[: match.start()].lower()
        after = text[match.end() :].lower()
        clause_before = _CLAUSE_SPLIT.split(before)[-1]
        clause_after = _CLAUSE_SPLIT.split(after)[0]
        if _NEGATION_BEFORE.search(clause_before):
            continue
        if _NEGATION_AFTER.search(clause_after):
            continue
        return True
    return False


def _max_size_near(text: str, stem: str) -> float | None:
    sizes = []
    for clause in _CLAUSE_SPLIT.split(text):
        if re.search(stem, clause, re.IGNORECASE):
            sizes += [float(v.replace(",", ".")) for v in _SIZE_MM.findall(clause)]
    return max(sizes) if sizes else None


def _raise(m: Markers, level: Urgency, reason: str) -> None:
    m.reasons.append(reason)
    if URGENCY_ORDER.index(level) > URGENCY_ORDER.index(m.urgency_floor):
        m.urgency_floor = level


def _apply_floor(m: Markers) -> None:
    if "pneumothorax" in m.findings:
        _raise(m, "emergency", "пневмоторакс")

    worst_birads = max(m.birads, default=None)
    if worst_birads in (4, 5, 6):
        _raise(m, "priority", f"BI-RADS {worst_birads}")
    elif worst_birads in (0, 3):
        _raise(m, "planned", f"BI-RADS {worst_birads}")

    if m.max_nodule_mm is not None:
        if m.max_nodule_mm > 8:
            _raise(m, "priority", f"очаг {m.max_nodule_mm:g} мм")
        elif m.max_nodule_mm >= 6:
            _raise(m, "planned", f"очаг {m.max_nodule_mm:g} мм")
    elif "lung_nodule" in m.findings:
        _raise(m, "planned", "очаг без указания размера")

    for finding in ("pleural_effusion", "cardiomegaly", "emphysema"):
        if finding in m.findings:
            _raise(m, "planned", finding)
    if m.cac_drs or "coronary_calcium" in m.findings:
        _raise(m, "planned", "кальциноз коронарных артерий")
