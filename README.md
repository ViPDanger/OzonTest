# OzonTest 
тестовое задания для Ozon SDET mock www.cbr.ru

## ============Описание сервиса===============

### Возможные аргументы
- --host            - host сервиса в GIN  (default 127.0.0.1:8080)
- --mongoURI        - URI mongo db        (default 127.0.0.1:27017)
- --mongodbName     - mongo db name       (default xml_daily)
- --mongoUser       - mongo user          (default admin)
- --mongoPassword   - mongo password      (default admin)
- --grpcHost        - host сервиса в gRPC (default 127.0.0.1:2020)

## Gin routes
- GET /   - получение XML документа ValCurs на дату и имя, установленные по gRPC (default date: current date, name `Foreign Currency Market`)

### Ручки gRPC

- rpc AddValCurs(AddValCursRequest{ValCurs}) returns (AddValCursResponse{message}); - добавить в mongo данные по курсу валют

- rpc DeleteValCurs(DeleteValCursRequest{date string,name string}) returns (DeleteValCursResponse{message}); - удалить в монго данные по курсу валют

- rpc SetState(SetStateRequest{date string, name string}) returns (SetStateResponse{message}) - Установить состояние GIN handler 

- rpc GetState(GetStateRequest{}) returns (GetStateResponse{date string, name string}); -  Просмотреть состояние GIN handler 

- rpc Reset(ResetRequest{}) returns (ResetResponse{message}); - Очистить всю бд

### Архитектура
- `cmd/` — точка входа
- `docs/` — Swagger, описание архитектуры
- `internal/` - внутренняя часть приложения
- `internal/` - внутренняя часть приложения
- `internal/app` - собранное приложение с gRPC и gin
- `internal/domin` - entity и интерфейсы репозиториев
- `internal/interfaces` - контроллеры приложения (grpc/gin handlers)
- `interal/interfaces\dto` - Data Transfer Objects GIN
- `interal/interfaces\gRPC` - gRPC handler-ы
- `interal/interfaces\handlers` - gin handler-ы
- `interal/interfaces\mapper` - мапперы для proto и DTO к entity обьектам сервиса
- `internal/mongodb` - реализация репозиториев на mongoDB
- `internal/usecase` — бизнес-логика (слой приложения)
- `proto/` — .proto файлы и их генерации



## =========Описание задачи================

### Вводные данные по окружению:
- Микросервисная архитектура
- Асинхронное взаимодействие
- Единый тестовый контур
- Сотни параллельно запускаемых е2е автотестов

### Суть задания
Написать мок-сервис имитирующий ответ от
https://www.cbr.ru/scripts/XML_daily.asp?date_req=02/03/2002

### Юзкейсы
- успешное получение данных, http код 200
- не успешное получение данных, http код 500

### Доп. требования:
- под каждый тест предполагается формировать уникальные данные/ответы
- общий хард-код ответ недопустим
- нельзя использовать хэддеры ввиду технического ограничения

### Будет плюсом:
- makefile/taskfile с командой для запуска
- наличие proto
- наличие swagger
- golang
- использование БД

## ----------Коментарии разработчкика-----------
- Добавил идентификацию пользователя по IP(без порта). В случае большого количества клиентов запросы могут мешать тесту друг друга, поэтому было принято решения добавить идентификацию.
  Стоит учесть, что порты у запроса GET / и gRPC клиента могут быть разные, соответственно идентефицировать по ним не получится. Это приводит нас к тому, что два процесса с одного адреса могут всё ещё мешать тестированию
  друг с друга, но как идентефицировать клиентов по другому без Header-ов я исключительно без понятия.
- Добавил удаление данных из бд запроса на их просмотр. Не уверен, насколько это необходимо, но это поможет держать базу чистой, к тому же в локальных моках так и сделано.
- Всё ещё не уверен насчёт Swagger. Им я не разу до этого не занимался и не уверен, правильно ли я подошёл к его генерации
