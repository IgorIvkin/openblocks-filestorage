# OpenBlocks "Хранилище файлов"

### Инициатива OpenBlocks

Инициатива OpenBlocks &mdash; это проект с открытым исходным кодом, целью которого
является предоставить открытые и масштабируемые решения уровня предприятия.

### Описание

Сервис "Хранилище файлов" представляет собой API для хранения файлов в больших бинарных объектах - блобах. 
Термин "блоб" происходит от английского **Binary Large Object**, что означает "большой бинарный объект".
В роли блобов выступают базы данных SQLite.

В сервисе может быть создано несколько блобов, каждый из которых представлен отдельной базой данных SQLite, 
каждая такая БД называется "том". При сохранении файла выбирается случайный том, который не занят сохранением файлов,
если же все тома заняты, сервис ожидает освобождения первого попавшегося тома.

Наличие множества томов даёт возможность сохранять много файлов одновременно, можно сохранять одновременно столько файлов, сколько в сервисе
создано томов.

Структура каждого тома единообразна и представляет собой отдельную таблицу для хранения файлов (файлы хранятся в бинарном виде),
а также таблицу мета-информации, в которой хранится служебная информация о томе.

### Основная конфигурация

Обратите внимание на следующие секции в конфигурации сервиса.

```yaml
storage:
  path: /s3/storage
  volumes: 1
  max-volume-size: 10
```

Параметр `path` задаёт полный путь в файловой системе к месту, где будут сохраняться базы данных SQLite, играющие
роль томов-хранилищ для файлов.

Параметр `volumes` задает, сколько будет создано отдельных баз данных, каждая отдельная база данных называется
томом в составе общего хранилища.

Параметр `max-volume-size` предполагается использовать для ограничения размера каждого тома. Пока что этот параметр не используется и 
зарезервирован на будущее.


### API

POST /api/v1/store

```bash
curl --location 'http://localhost:8903/api/v1/store' \
--form 'file_to_store=@"/home/user/test.bin"'
```

Сохраняет файл в один из томов. В случае успеха возвращается идентификатор файла следующего вида.

```json
{
    "file": "0-5"
}
```

Здесь первое число представляет собой идентификатор тома, а второе число представляет собой идентификатор файла в пределах тома.



## Полезные ссылки
* [OpenBlocks "Пользователи"](https://github.com/IgorIvkin/openblocks-users)
* [OpenBlocks "Роли"](https://github.com/IgorIvkin/openblocks-roles)
* [OpenBlocks "Команды"](https://github.com/IgorIvkin/openblocks-teams)
* [Сервис "Ограничитель запросов" на Java](https://github.com/IgorIvkin/openblocks-ratelimiter)
* [Документация по Keycloak](https://www.keycloak.org/documentation)