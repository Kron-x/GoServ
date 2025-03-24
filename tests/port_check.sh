#!/bin/bash

# Пути к файлам
FIRST_FILE="./go/config.json"
SECOND_FILE="docker-compose.yml"

# Извлечение значения порта из первого файла
PORT_FIRST_FILE=$(jq -r '.port' "$FIRST_FILE")

# Извлечение значения порта из второго файла (после ":")
PORT_SECOND_FILE=$(grep "ports:" -A 1 "$SECOND_FILE" | grep -oP ':(\d+)' | tr -d ':')

# Проверка соответствия
if [[ "$PORT_FIRST_FILE" == "$PORT_SECOND_FILE" ]]; then
  echo "Порты совпадают: $PORT_FIRST_FILE"
else
  echo "Порты не совпадают! Первый файл: $PORT_FIRST_FILE, Второй файл: $PORT_SECOND_FILE"
fi