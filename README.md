# bank-app-backend

Сервис, сервис по работе с банковскими операциями.

## API сервиса

[Спецификация OpenAPI](docs/swagger.yaml)

### Endpoints

| HTTP method  | Endpoint              | Описание                              |
|--------------|-----------------------|---------------------------------------|
| GET          | `/auth/me`            | Получить профиль пользователя         |
| GET          | `/auth/accounts`      | Получить список счетов пользователя   |
| POST         | `/auth/accounts`      | Создать участника авиамероприятия     |
| GET          | `/auth/accounts/:id`  | Получить детальную инфу о счёте по ID |
| PATCH        | `/auth/accounts/:id`  | Закрыть счёт при нулевом балансе      |
| GET          | `/users`              | Получить список пользователей         |
| PATCH        | `/users/:id`          | Обновить информацию о пользователе    |
| POST         | `/register`           | Регистрация пользователя              |
| POST         | `/login`              | Авторизация пользователя              |
| POST         | `/refresh`            | Обновление токена авторизации         |
| POST         | `/logout`             | Выход пользователя                    |
| GET          | `/swagger/*any`       | Документация swagger                  |
| GET          | `/metrics`            | Сбор метрик prometheus                |

Каждый endpoint с префиксом /auth требует токен авторизации

## Запуск

```bash
cd cmd
go run main.go
```

## Сборка

```bash
cd cmd
go build
```

### Минимальная версия языка go

1.24

### Зависимости от сторонних библиотек

- [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
- [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt)
- [github.com/ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- [github.com/prometheus/client_golang](https://github.com/prometheus/client_golang)
- [github.com/swaggo/files](https://github.com/swaggo/files)
- [go.uber.org/zap](https://go.uber.org/zap)
- [golang.org/x/crypto](https://golang.org/x/crypto)
- [github.com/ilyakaznacheev/cleanenv](https://github.com/ilyakaznacheev/cleanenv)
- [gorm.io/driver/postgres](https://github.com/go-gorm/postgres)
- [gorm.io/gorm](https://github.com/go-gorm/gorm)

## Запуск

Сервис корректно обрабатывает сигналы SIGINT и SIGTERM, настроен graceful shutdown

### База данных

Используется ORM gorm, подключена AutoMigrate

### Файлы конфигурации

***local-yaml:***

[Пример файла конфигурации](config/local.yaml)

### Переменные окружения

CONFIG_PATH="Path" (Path - путь до локального конфига)

## Мониторинг

```bash
curl -X 'GET' \
  'http://0.0.0.0:8080/' \
  -H 'accept: application/json'
```

### Структура проекта

```bash
├── cmd                 // точка входа
├── config              // конфиг запуска приложения
├── docs                // swagger-документация
└──internal             // внутренности приложений
│   ├── app             // код основного приложения
│   ├── config          // всё что относится к конфигурации приложения
│   ├── controllers     // контроллеры приложения
│   │   └── http        // HTTP (REST) контроллеры
│   ├── db              // подключение к БД
│   ├── entities        // основные сущности сервиса (модели)
│   ├── lib             // библиотеки и утилиты
│   ├── repository      // слой взаимподейсвтия с БД
│   ├── server          // запуск сервера, graceful shutdown
│   └── services        // бизнес логика
└── test                // тесты и тестовая инфраструктура
```

### Swagger генерация кода для endpoints

Чтобы сгенерировать swagger-документацию:

```
swag init -g cmd/main.go
```
