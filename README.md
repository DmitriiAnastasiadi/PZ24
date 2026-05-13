# PZ8: CI/CD для Go backend-проекта

## Тема и цель практической работы

- Тема: настройка CI/CD для backend-проекта на Go.
- Цель: освоение автоматического pipeline для проверки, сборки, упаковки Docker-образа и подготовки приложения к доставке.

## Архитектура проекта

В репозитории реализованы два сервиса: `auth` и `tasks`. Для практики выбран сервис `tasks`, поскольку он является завершённым backend-сервисом с собственным HTTP-обработчиком, внутренним сервисным слоем и отдельным исполняемым модулем.

Выбор `tasks` обусловлен тем, что:
- сервис уже компилируется самостоятельно;
- пакет `services/tasks` содержит все необходимые зависимости;
- он легко проверяется локально и упаковывается в Docker.

## Что такое CI и CD

- CI (Continuous Integration) — непрерывная интеграция. После каждого изменения проект автоматически проверяется: зависимостей, тестов и сборки.
- CD (Continuous Delivery / Continuous Deployment) — готовность к доставке или автоматическое развёртывание. CD отвечает за упаковку результата и доставку на целевую платформу.

## Структура pipeline

Pipeline реализован на GitHub Actions и разбит на два job:

1. `test-and-build`
   - checkout репозитория;
   - установка Go;
   - загрузка зависимостей;
   - запуск тестов;
   - компиляция приложения.
2. `docker-build`
   - выполняется после успешного `test-and-build`;
   - настраивает Docker Buildx;
   - собирает Docker-образ.

## Выбранная платформа

Выбран GitHub Actions.

## Файл pipeline

Файл: `.github/workflows/ci.yml`

```yaml
name: CI Pipeline

on:
  push:
    branches: [ "main", "master" ]
  pull_request:
    branches: [ "main", "master" ]

jobs:
  test-and-build:
    runs-on: ubuntu-latest

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Setup Go
        uses: actions/setup-go@v5
        with:
          go-version: '1.23'

      - name: Show Go version
        run: go version

      - name: Download dependencies
        run: go mod tidy
        working-directory: ./services/tasks

      - name: Run tests
        run: go test ./...
        working-directory: ./services/tasks

      - name: Build application
        run: go build ./...
        working-directory: ./services/tasks

  docker-build:
    runs-on: ubuntu-latest
    needs: test-and-build

    steps:
      - name: Checkout repository
        uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build Docker image
        run: docker build -t techip-tasks:${{ github.sha }} -f Dockerfile .
```

## Пояснение шагов pipeline

- `Checkout repository` — получает актуальный код.
- `Setup Go` — подготавливает окружение Go 1.23.
- `Show Go version` — проверка установленной версии.
- `Download dependencies` — запускает `go mod tidy` в директории `services/tasks`, чтобы проверить зависимости.
- `Run tests` — запускает `go test ./...` в директории сервиса.
- `Build application` — выполняет `go build ./...` в директории сервиса.
- `docker-build` — собирает Docker-образ только после успешных тестов и сборки.

## Формирование тега образа

Тег Docker-образа формируется на основе хеша коммита: `${{ github.sha }}`.

Это позволяет точно идентифицировать версию, из которой собран образ.

## Где хранятся секреты

Секреты и переменные CI должны храниться в защищённом хранилище GitHub Secrets или GitLab CI/CD Variables.

Нельзя хранить их:
- в репозитории;
- в YAML-файле;
- в `.env`, который попадает в контролируемый Git.

## Локальная проверка

Для подтверждения работоспособности проекта выполнены команды:

```powershell
cd c:\Users\dimma\Desktop\ПИШ\2 СЕМ\Технологии создания ПО\PZ24
cd services\tasks
go test ./...
go build ./...
cd ..\..
docker build -t techip-tasks:local -f Dockerfile .
```

## Результаты проверки

- `go test ./services/tasks/...` — успешно.
- `go build ./services/tasks/...` — успешно.
- `docker build -t techip-tasks:local -f Dockerfile .` — успешно.

## Публикация образа и деплой

В этом проекте публикация образа в registry не настроена, но базовая логика выглядит так:

- авторизация в registry через секреты `REGISTRY_USERNAME` и `REGISTRY_PASSWORD`;
- сборка образа с тегом на основе коммита;
- push в `ghcr.io/my-org/techip-tasks:${{ github.sha }}`.

Минимальный деплой на VPS может включать:
- SSH-подключение из pipeline;
- `docker pull ghcr.io/my-org/techip-tasks:<tag>`;
- `docker compose up -d`.

## Контрольные вопросы

1. Чем CI отличается от CD?
   - CI отвечает за проверку кода, запуск тестов и сборку.
   - CD отвечает за подготовку доставки и/или автоматическое развёртывание.

2. Почему pipeline должен запускать тесты?
   - Чтобы гарантировать, что изменения не ломают существующую функциональность.
   - Тесты автоматически обнаруживают баги до интеграции.

3. Зачем нужен автоматический build?
   - Чтобы убедиться, что проект компилируется в стандартизованном окружении.
   - Это исключает разрывы между локальным и CI-окружением.

4. Почему важно собирать Docker-образ в CI, а не только локально?
   - Потому что CI гарантирует повторяемость сборки и проверяет, что Dockerfile работает.
   - Локальная сборка может скрыть проблемы окружения и зависимостей.

5. Что такое CI secrets?
   - Это защищённые значения (токены, пароли, ключи), доступные только в CI.
   - Они применяются в pipeline без сохранения в репозитории.

6. Почему нельзя хранить токены и SSH-ключи в репозитории?
   - Потому что это угрожает безопасности и позволяет получить несанкционированный доступ.
   - Репозиторий может стать доступен другим людям или утечь.

7. Для чего нужен тег Docker-образа?
   - Чтобы однозначно идентифицировать собранную версию.
   - Тег помогает воспроизводить и откатывать конкретные сборки.

8. Что делает job `docker-build`?
   - Собирает Docker-образ после успешного `test-and-build`.
   - Гарантирует, что контейнер строится только из проверенного кода.

9. Почему в multi-service проекте важен `working-directory`?
   - Потому что в репозитории несколько сервисов, и команды должны выполняться в нужной директории.
   - Неправильный `working-directory` приведёт к ошибке сборки или невозможности найти `go.mod`.

10. Какие риски возникают при полностью автоматическом деплое?
   - Возможный разворот в продакшн с багами, если тесты недостаточны.
   - Ошибки конфигурации или секретов могут привести к срыву сервиса.
