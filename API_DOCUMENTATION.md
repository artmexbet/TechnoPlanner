# API Documentation - Request Updates

## Обновленные эндпоинты для работы с заявками

### 1. Обновить заявку
**PATCH** `/api/v1/requests/:id`

**Требования:** JWT токен, роль Admin

**Request Body:** (все поля опциональные)
```json
{
  "request_text": "string",
  "status": "pending|assigned|in_progress|completed|rejected|canceled",
  "schedule_time": "2024-01-01T10:00:00Z",
  "address": "string",
  "responsible_id": "uuid"
}
```

**Response:**
```json
{
  "id": "uuid",
  "request_text": "string",
  "status": "string",
  "schedule_time": "string",
  "end_time": "string",
  "address": "string",
  "responsible_id": "uuid",
  "equipment": [],
  "created_at": "string",
  "updated_at": "string"
}
```

## Существующие эндпоинты для заявок

### 2. Получить список заявок
**GET** `/api/v1/requests`

**Query Parameters:**
- `responsible_id` (опционально) - UUID ответственного

**Требования:** JWT токен

### 3. Получить заявку по ID
**GET** `/api/v1/requests/:id`

**Требования:** JWT токен

### 4. Назначить ответственного за заявку
**POST** `/api/v1/requests/:id/responsible`

**Требования:** JWT токен, роль Admin

**Request Body:**
```json
{
  "responsible_id": "uuid"
}
```

## NATS Subjects

### Gateway -> Requests Service

- `gateway.requests.list` - список заявок
- `gateway.requests.get` - получить заявку
- `gateway.requests.assign.responsible` - назначить ответственного
- `gateway.requests.update` - обновить заявку

## Изменения в сервисах

### Requests Service

**Новые методы в RequestService:**
- `UpdateRequest(ctx, requestID, updates)` - обновление заявки
- `ListResponsibles(ctx)` - список ответственных
- `SaveResponsible(ctx, id, username)` - сохранение ответственного

**Новые NATS обработчики:**
- `handleGatewayUpdateRequest` - обработка обновления заявки
- `handleGatewayListResponsibles` - обработка запроса списка ответственных
- `handleGatewayCreateResponsible` - обработка создания ответственного

### Gateway Service

**Новые клиенты:**
- `ResponsibleClient` - клиент для работы с ответственными

**Новые сервисы:**
- `ResponsibleService` - сервис для работы с ответственными

**Новые роуты:**
- `InitResponsibleRoutes()` - инициализация роутов для ответственных

**Новые обработчики:**
- `listResponsibles()` - список ответственных
- `createResponsible()` - создание ответственного
- `updateRequest()` - обновление заявки

## Структуры данных

### DTO (libs/dto/models.go)

```go
type Responsible struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
}

type ResponsibleCreateRequest struct {
    ID       uuid.UUID `json:"id"`
    Username string    `json:"username"`
}

type RequestUpdateRequest struct {
    RequestID     uuid.UUID      `json:"request_id"`
    RequestText   *string        `json:"request_text,omitempty"`
    Status        *RequestStatus `json:"status,omitempty"`
    ScheduleTime  *string        `json:"schedule_time,omitempty"`
    Address       *string        `json:"address,omitempty"`
    ResponsibleID *uuid.UUID     `json:"responsible_id,omitempty"`
}
```

### Domain (services/requests/internal/domain/domain.go)

```go
type RequestUpdate struct {
    RequestText   *string
    Status        *StatusType
    ScheduleTime  *string
    Address       *string
    ResponsibleID *uuid.UUID
}

type Responsible struct {
    ID       uuid.UUID
    Username string
}
```

