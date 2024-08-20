# KV хранилище на базе Tarantool
## Описание задания

Реализовать API для KV-хранилища на базе Tarantool.

## Описание API

Зарегистрированы следующие эндпоинты:

* POST /api/login - для получения токена авторизации
* POST /api/write - для записи данных пачками
* POST /api/read - для чтения данных пачками

Схема API представлена в файле [openapi.yaml](api/openapi.yaml)

## Запуск проекта

1. Необходимо создать .env файл с помощью команды

```shell
make env
```

В файле может задать следующие параметры

TARANTOOL_USER_NAME - логин администратора Tarantool
TARANTOOL_USER_PASSWORD - пароль администратора
TARANTOOL_PORT - порт Tarantool
TARANTOOL_REQUEST_TIMEOUT - таймаут на исполнение запроса к Tarantool

APP_PORT - порт сервиса
APP_HOST - хост сервиса 
TOKEN_TTL - время жизни токена аутентификации
APP_SECRET_KEY - секретный ключ приложения для JWT
TARANTOOL_ADDRESS=tarantool:4000 - адрес TARANTOOL для приложения

2. Запустить контейнеры docker 

```shell
make deploy
```

## Тестирование приложения

Приложение использует протокол HTTP для обмена данными между сервером и клиентом.
В качестве токена авторизации используются JWT токены.

Примеры запросов показаны на рисунках ниже:

Все запросы имеют валидацию `json`

![invalid_json.png](./docs/test_requests/invalid_json.png)

У запросов к `/write` и `/read` должен указываться токен для доступа к данным, иначе система не пропустит запрос.

Стандартные хедеры:

![headers.png](./docs/test_requests/headers.png)

Ответы при невалидном токене:

![invalid_token.png](./docs/test_requests/write/invalid_token.png)
![token.png](./docs/test_requests/read/token.png)


1) Запросы к `/login`

![empty_data_2.png](./docs/test_requests/login/empty_data_2.png)
![empty_fields.png](./docs/test_requests/login/empty_fields.png)
![invalid_creds.png](./docs/test_requests/login/invalid_creds.png)
![valid_data.png](./docs/test_requests/login/valid_data.png)

------

Запросы к `/read`

![empty_req.png](./docs/test_requests/read/empty_req.png)
![full_output.png](./docs/test_requests/read/full_output.png)
![part_output.png](./docs/test_requests/read/part_output.png)
![token.png](./docs/test_requests/read/token.png)