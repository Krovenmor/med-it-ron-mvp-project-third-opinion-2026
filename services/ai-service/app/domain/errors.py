"""Ошибки, после которых оценку нельзя выдать. Все они дают не-200, и main-B2B повторит запрос."""


class CatalogUnavailableError(Exception):
    """Справочник услуг не получен из МИС ни разу."""


class LLMError(Exception):
    """LLM не ответила или ответила непригодно."""


class InvalidDraftError(Exception):
    """Черновик LLM нельзя превратить в валидный ответ."""
