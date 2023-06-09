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

Отличительной особенностью сервиса является интеграция с **Keycloak** вашего предприятия. Сервис может принимать и верифицировать
JWT-токены.

### Основная конфигурация

Обратите внимание на следующие секции в конфигурации сервиса.

```yaml
storage:
  path: /s3/storage
  volumes: 1
  max-volume-size: 10
general:
  use-jwt-auth: true
  jwks-url: "http://localhost:8534/realms/infra/protocol/openid-connect/certs"
  except-urls:
    - /api/v1/file/
```

Параметр `storage.path` задаёт полный путь в файловой системе к месту, где будут сохраняться базы данных SQLite, играющие
роль томов-хранилищ для файлов.

Параметр `storage.volumes` задает, сколько будет создано отдельных баз данных, каждая отдельная база данных называется
томом в составе общего хранилища.

Параметр `storage.max-volume-size` предполагается использовать для ограничения размера каждого тома. Пока что этот параметр не используется и 
зарезервирован на будущее.

Секция параметров `general` содержит ряд параметров для работы с JWT-токенами и с сервисом **Keycloak** вашего предприятия.
Параметр `use-jwt-auth` определяет, будет ли сервис требовать наличия JWT-токена в запросах. В параметре `except-urls` можно
передать список URL, которые будут работать без токена. В параметре `jwks-url` передаётся URL, по которому можно получить
открытые JWKS-ключи из вашего Keycloak.


### API

POST /api/v1/store

```bash
curl --location 'http://localhost:8903/api/v1/store' \
--form 'file_to_store=@"/home/user/test.bin"'
--form 'file_type="application/octet-stream"'
```

Сохраняет файл в один из томов. В качестве типа сохраняемого файла используется значение параметра `file_type`. Задавать это
значение важно для последующего адекватного получения файла. В случае успеха возвращается идентификатор файла следующего вида.

```json
{
    "file": "0-5"
}
```

Здесь первое число представляет собой идентификатор тома, а второе число представляет собой идентификатор файла в пределах тома.

---

GET /api/v1/file/{fileId}

Получает файл по его идентификатору. Идентификатор файла выглядит как `0-1`, где первое число представляет собой идентификатор
тома в составе общего хранилища, а второе число &mdash; идентификатор файла в томе. В случае успеха возвращается содержимое файла
с типом контента, соответствующим сохранённому в базе.



## Полезные ссылки
* [OpenBlocks "Пользователи"](https://github.com/IgorIvkin/openblocks-users)
* [OpenBlocks "Роли"](https://github.com/IgorIvkin/openblocks-roles)
* [OpenBlocks "Команды"](https://github.com/IgorIvkin/openblocks-teams)
* [Сервис "Ограничитель запросов" на Java](https://github.com/IgorIvkin/openblocks-ratelimiter)
* [Документация по Keycloak](https://www.keycloak.org/documentation)