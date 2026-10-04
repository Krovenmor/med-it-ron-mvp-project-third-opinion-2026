FROM python:3.13-slim

WORKDIR /srv
RUN pip install --no-cache-dir \
    "fastapi>=0.142" "gigachat>=0.2.3" "httpx>=0.28" "pydantic>=2.13" \
    "python-dotenv>=1.2" "uvicorn[standard]>=0.54"

COPY app ./app
COPY knowledge ./knowledge
# default.env — запасной источник несекретных настроек; переменные из env_file compose важнее.
COPY config/default.env ./config/default.env

EXPOSE 8083
CMD ["python", "-m", "app"]
