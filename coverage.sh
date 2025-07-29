#!/bin/bash

# Сначала экспортируй переменные окружения
export GO_COVER_IGNORE_SPEC_PATH=".coverage-ignore.yaml"
export GO_COVER_IGNORE_COVER_PROFILE_PATH="coverage.out"

# Выполни тесты с генерацией покрытия
go test ./... -coverprofile=coverage.out -count=1

# Убери игнорируемые файлы из покрытия
go-cover-ignore

# Покажи итоговое покрытие
go tool cover -func=coverage.out

# Удали временный файл
rm coverage.out