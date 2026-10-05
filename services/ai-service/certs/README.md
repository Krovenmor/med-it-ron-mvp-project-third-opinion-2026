Сюда кладётся корневой сертификат НУЦ Минцифры для GigaChat (https://www.gosuslugi.ru/crt, формат PEM).
В .env: GIGACHAT_CA_BUNDLE_FILE=certs/russian_trusted_root_ca.pem — путь относительно папки сервиса,
он работает и локально, и в docker-compose (папка монтируется в контейнер).
Сертификат публичный, секретом не является.
