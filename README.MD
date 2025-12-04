# Yayo CLI

yayo-cli — CLI утилита для управления ENV переменными проектов с безопасным локальным хранилищем (SQLCipher) и возможностью работы через SSH.

## Установка

Установи утилиту:

`go install github.com/1Storm3/yayo-cli`

## Использование
Команды

**init**

Инициализирует локальное зашифрованное хранилище ENV переменных.

`yayo-cli init`

**destroy**

Удаляет локальное хранилище проекта и все связанные файлы.

Флаги:

-p — название проекта

`yayo-cli destroy -p molybdenum`

**add-bulk**

Добавляет или обновляет несколько ENV переменных через stdin.

Флаги:

-p — название проекта

-pass — пароль SQLCipher БД


`echo '[{"key":"FOO","value":"BAR","service":"backend"}]' | yayo-cli add-bulk -p molybdenum -pass "test"`

**list-envs**

Выводит список ENV переменных проекта локально.

Флаги:

-p — название проекта

-s — фильтр по сервису (необязательно)

`yayo-cli list-envs -p molybdenum`


**list**

Выводит список ENV переменных проекта на удалённой машине через SSH.

Флаги:

-h — SSH host, например root@1.2.3.4

-p — название проекта

-s — фильтр по сервису (необязательно)


`yayo-cli list -h root@1.2.3.4 -p molybdenum -s backend`

**push**

Загружает локальный .env на удалённую машину через SSH.

Флаги:

-h — SSH host, например root@1.2.3.4

-p — название проекта

-s — название сервиса

-f — путь к .env файлу


`yayo-cli push -h root@1.2.3.4 -p molybdenum -s backend -f .env`