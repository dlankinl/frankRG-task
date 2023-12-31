# Тестовое задание для FrankRG на позицию Golang Backend

## Техническое задание
```
Файловый менеджер в браузере.
Возможности: 
  - загружать/удалять файлы;
  - создавать папки;
  - скачивать файлы;
  - переименовать файлы/папки;
  - просматривать содержимое папок;
  - переименовывать файлы/папки.
```

Файловая система работает на базе реляционной СУБД **PostgreSQL** + **REST API**.

## Инструкция по запуску
```
docker-compose up
```

## Стек

Backend: Golang, PostgreSQL (migrations), REST API (go chi), docker, docker-compose.
Frontend: JS.

## Документация

### Создать файл/директорию - POST /api/createfile
Принимает на вход json с параметрами файла (среди которых флаг создания директории) и создает файл в текущей директории (из веб-интерфейса).

#### Параметры
- name (string) - наименование файла/директории;
- size (int) - размер файла (в случае директории 0);
- content ([]byte) - содержимое файла (в случае директории - пустой массив);
- is_dir (bool) - флаг создания директории;
- parent_dir (string) - название родительской директории.

#### Пример возможного запроса:
```json
{
    "name": "newDir",
    "content": null,
    "size": 0,
    "is_dir": true,
    "parent_dir": "root"
}
```

#### Пример успешного ответа:
```json
{
    "ID": 3,
    "Name": "newDir",
    "Size": 0,
    "ModTime": "2023-09-19T20:32:15.984192662Z",
    "IsDirectory": true,
    "Content": "",
    "ParentID": 1
}
```

####  Пример неудачного ответа:
```json
{
    "error": true,
    "message": "no parentDir found"
}
```

### Загрузить файл - POST /api/uploadfile/{parent_dir_name}
Принимает на вход массив байтов через форму в JS и создает файл в данной директории.

### Переименовать файл/директорию - POST /api/file
Принимает на вход json с данными о переименовании. Идентификатор файла/директории определяется на стороне клиента.

#### Параметры
- id (int) - идентификатор файла/директории;
- new_name (string) - новое наименование.

#### Пример возможного запроса:
```json
{
    "id": 5,
    "new_name": "new_DIR"
}
```

#### Пример успешного ответа:
```json
{
    "status": "OK"
}
```

#### Пример неудачного ответа:
```json
{
    "error": true,
    "message": "nothing found to update"
}
```

### Получить содержимое директории - GET /dir/{name}
Принимает на вход наименование директории (наличие в системе директории с таким же именем обходится путем получения id директории-родителя в момент запроса со стороны клиента) и возвращает ее содержимое. 

Возвращает html-разметку страницы с учетом содержимого директории (go template).

#### Параметры
- name (string) - наименование директории.

#### Пример возможного запроса:
```
curl --location 'http://localhost:8080/dir/root'
```

#### Пример успешного ответа:
```html
<!--огромная пелена разметки страницы, удовлетворяющей шаблону...-->
```

#### Пример неудачного ответа:
```json
{
    "error": true,
    "message": "dir wasn't found"
}
```

### Получить содержимое файла - GET /file/{id}/{name}
Принимает на вход наименование файла, его идентификатор и возвращает его содержимое.

#### Параметры
- name (string) - наименование файла.

#### Пример возможного запроса:
```
curl --location 'http://localhost:8080/file/6/hello_world.txt'
```

#### Пример успешного ответа:
```json
{
    "filename": "hello_world.txt",
    "id": 6,
    "content": "Hello, Frank RG team!"
}
```

#### Пример неудачного ответа:
```json
{
    "error": true,
    "message": "no files were found"
}
```

### Скачать файл - GET /api/downloadfile/{id}
Принимает на вход id интересующего файла и запускает загрузку со стороны клиента.

#### Параметры
- id (int) - идентификатор файла.

#### Пример возможного запроса:
```
curl --location 'http://localhost:8080/api/downloadfile/6'
```

#### Пример успешного ответа - **скачанный файл**.

#### Пример неудачного ответа:
```json
{
    "error": true,
    "message": "no files were found"
}
```

### Удалить файл/директорию - DELETE /api/file/{id}
Удаляет файлы или директории. Идентификатор файла/директории определяется на стороне клиента.

#### Параметры
- id (int) - идентификатор файла/директории.

#### Пример возможного запроса:
```
curl --location --request DELETE 'http://localhost:8080/api/file/4'
```

#### Пример успешного ответа:
```json
{
    "status": "OK",
    "deleted_rows": 1
}
```

#### Пример неудачного ответа:
```json
{
    "error": true,
    "message": "no files were found"
}
```