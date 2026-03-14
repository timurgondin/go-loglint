# go-loglint

[![CI](https://img.shields.io/github/actions/workflow/status/timurgondin/go-loglint/ci.yml?style=flat-square&label=CI)](https://github.com/timurgondin/go-loglint/actions/workflows/ci.yml)
![Go](https://img.shields.io/badge/go-1.25-00ADD8?style=flat-square&logo=go&logoColor=white)
![golangci-lint](https://img.shields.io/badge/golangci--lint-plugin-4A90D9?style=flat-square&logo=go&logoColor=white)

Линтер для Go, который проверяет стиль и содержимое сообщений в вызовах функций логирования. Поддерживает `log/slog` и `go.uber.org/zap`, интегрируется в [golangci-lint](https://golangci-lint.run/) как плагин или запускается как самостоятельный инструмент.

## Содержание

- [Поддерживаемые логгеры](#поддерживаемые-логгеры)
- [Правила](#правила)
- [Установка и использование](#установка-и-использование)
- [Конфигурация](#конфигурация)
- [Автоисправление](#автоисправление)

## Поддерживаемые логгеры

| Пакет | Методы |
|-------|--------|
| `log/slog` | `Info`, `Error`, `Warn`, `Debug` и их `*Context`-варианты |
| `go.uber.org/zap` | `Info`, `Error`, `Warn`, `Debug`; а также через `zap.L()` и `zap.S()` |

## Правила

### 1. `check-lowercase` — сообщение должно начинаться со строчной буквы

```go
// Плохо
slog.Info("Starting server on port 8080")
slog.Error("Failed to connect")

// Хорошо
slog.Info("starting server on port 8080")
slog.Error("failed to connect")
```

Поддерживает автоматическое исправление (`--fix`).

### 2. `check-english` — сообщение должно быть на английском

Проверяет, что все символы находятся в диапазоне ASCII (≤ 127).

```go
// Плохо
slog.Info("запуск сервера")
slog.Error("ошибка подключения к базе данных")

// Хорошо
slog.Info("starting server")
slog.Error("failed to connect to database")
```

Автоматическое исправление не поддерживается.

### 3. `check-special-chars` — сообщение не должно содержать спецсимволы

Запрещены эмодзи и прочие символы вне ASCII, а также следующие символы:

```
! ? @ # $ % ^ & * ~ ` | \
```

```go
// Плохо
slog.Info("server started!!!")
slog.Warn("server started!!!")

// Хорошо
slog.Info("server started")
slog.Warn("something went wrong")
```

Поддерживает автоматическое исправление (`--fix`) — удаляет запрещённые символы.

### 4. `check-sensitive` — сообщение не должно содержать чувствительные данные

Проверяет наличие ключевых слов, связанных с чувствительной информацией (без учёта регистра, по целым словам).

Встроенные паттерны: `password`, `passwd`, `secret`, `token`, `api_key`, `apikey`, `auth`, `credential`, `private_key`, `credit_card`.

```go
// Плохо
slog.Info("user password: 123456")
slog.Debug("api_key=abc123")

// Хорошо
slog.Info("user authenticated successfully")
slog.Debug("api request completed")
```

Автоматическое исправление не поддерживается.

## Установка и использование

### Как плагин golangci-lint (рекомендуется)

Для интеграции в golangci-lint используется механизм [module plugins](https://golangci-lint.run/plugins/module-plugins/). Кастомный бинарник собирается командой `golangci-lint custom`.

**1. Создайте файл `.custom-gcl.yml` в корне проекта:**

```yaml
version: v2.8.0
name: custom-gcl
destination: .
plugins:
  - module: 'github.com/timurgondin/go-loglint'
    import: 'github.com/timurgondin/go-loglint/plugin'
    version: v1.0.0
```

**2. Соберите кастомный бинарник:**

```bash
golangci-lint custom
# или с подробными логами:
golangci-lint custom -v
```

Появится файл `./custom-gcl` (или с именем из поля `name`).

**3. Запустите линтер:**

```bash
./custom-gcl run ./...
./custom-gcl run --fix ./...
```

> `.custom-gcl.yml` используется только при сборке бинарника. Настройки самого линтера задаются в `.golangci.yml` и читаются при каждом запуске.

### Standalone

```bash
go install github.com/timurgondin/go-loglint/cmd/loglint@latest
```

```bash
# Запуск проверки
loglint ./...

# С отключёнными правилами
loglint -check-sensitive=false ./...

# С дополнительными паттернами чувствительных данных
loglint -extra-patterns="ssn,phone" ./...

# Автоисправление (только правила, поддерживающие fix)
loglint -fix ./...
```

## Конфигурация

### Флаги CLI (standalone)

| Флаг | Тип | По умолчанию | Описание |
|------|-----|--------------|----------|
| `-check-lowercase` | bool | `true` | Проверка заглавной буквы |
| `-check-english` | bool | `true` | Проверка английского языка |
| `-check-special-chars` | bool | `true` | Проверка спецсимволов |
| `-check-sensitive` | bool | `true` | Проверка чувствительных данных |
| `-extra-patterns` | string | `""` | Дополнительные паттерны через запятую |

```bash
loglint -check-lowercase=true -check-english=true -extra-patterns="ssn,phone_number" ./...
```

### Конфигурация `.golangci.yml`

```yaml
version: "2"

linters:
  default: none
  enable:
    - loglint
  settings:
    custom:
      loglint:
        type: module
        description: "checks log messages style"
        settings:
          check-lowercase: true
          check-english: true
          check-special-chars: true
          check-sensitive: true
          extra-patterns:
            - "ssn"
            - "phone_number"
            - "passport"
```

## Автоисправление

Флаг `--fix` (в golangci-lint) или `-fix` (в standalone) применяет автоматические исправления там, где это возможно.

| Правило | Поддержка `--fix` |
|---------|:-----------------:|
| `check-lowercase` | ✓ |
| `check-english` | — |
| `check-special-chars` | ✓ |
| `check-sensitive` | — |

```bash
# Standalone
loglint -fix ./...

# golangci-lint
./custom-gcl run --fix ./...
```
