# Контракт сервиса Requests

Сервис работает через **NATS** (Request-Reply и Pub/Sub).  
Транспортный формат — **JSON**.  
Стандартный ответ на Request-Reply — `dto.GatewayResponse`.

```go
// GatewayResponse — обёртка всех ответов
type GatewayResponse struct {
    Success bool            `json:"success"`
    Message string          `json:"message"`   // "success" | описание ошибки
    Data    json.RawMessage `json:"data,omitempty"`
}
```

---

## Общие типы

### Request (заявка)
```json
{
    "id": "uuid",
    "request_text": "string | null",
    "status": "pending | assigned | in_progress | completed | canceled | rejected",
    "schedule_time": "string",
    "end_time": "RFC3339",
    "address": "string",
    "responsible_user_id": "uuid | null",
    "equipment": [ { "request_id": "uuid", "equipment_id": 0, "quantity": 0, "created_at": "RFC3339", "updated_at": "RFC3339" } ],
    "audit": { "created_at": "RFC3339", "updated_at": "RFC3339", "created_by": "uuid | null", "updated_by": "uuid | null", "deleted_at": "RFC3339 | null" }
}
```

### Responsible (ответственный)
```json
{ "id": "uuid", "username": "string" }
```

---

## Request-Reply каналы

### Bot API

#### `bot.requests.created` — Создать заявку

**Request:**
```json
{
    "text": "string | null",
    "schedule_time": "string",
    "user_id": 123456789,
    "user_name": "string | null",
    "equipments": [ { "id": 1, "quantity": 2 } ],
    "equipment_string": "string | null",
    "address": "string"
}
```

**Response `data`:** `Request`

---

#### `gateway.requests.status.changed` — Обновить статус заявки

**Request:**
```json
{ "request_id": "uuid", "status": "pending | assigned | in_progress | completed | canceled | rejected" }
```

**Response `data`:** `"status updated"`

---

#### `bot.requests.get` — Получить заявку по ID

**Request:**
```json
{ "request_id": "uuid" }
```

**Response `data`:** `Request`

---

#### `bot.requests.list` — Получить список заявок пользователя (по Telegram ID)

**Request:**
```json
{ "telegram_id": 123456789, "limit": 20, "offset": 0 }
```

**Response `data`:** `Request[]`

---

#### `bot.requests.cancel` — Отменить заявку (от бота)

**Request:**
```json
{ "request_id": "uuid" }
```

**Response `data`:** `"request canceled"`

---

#### `bot-svc.equipment.add` — Добавить оборудование в справочник

**Request:**
```json
{
    "equipments": [ { "id": 1, "name": "Проектор", "quantity": 5 } ]
}
```

**Response `data`:** `"equipment added"`

---

### Gateway API

#### `gateway.requests.list` — Получить список заявок (с фильтром по ответственному)

**Request:**
```json
{ "responsible_id": "uuid | null" }
```

> Если `responsible_id` = `null` — вернуть все заявки.

**Response `data`:** `Request[]`

---

#### `gateway.requests.get` — Получить заявку по ID

**Request:**
```json
{ "request_id": "uuid" }
```

**Response `data`:** `Request`

**Ошибки:** `404 not found` если заявка не найдена.

---

#### `gateway.requests.assign.responsible` — Назначить ответственного

**Request:**
```json
{ "request_id": "uuid", "responsible_id": "uuid | null" }
```

> Если `responsible_id` = `null` — снять назначение.

**Response `data`:** `Request`

---

#### `gateway.requests.update` — Обновить заявку (частичное обновление)

**Request:**
```json
{
    "request_id": "uuid",
    "request_text": "string | null",
    "status": "pending | assigned | in_progress | completed | canceled | rejected | null",
    "schedule_time": "string | null",
    "address": "string | null",
    "responsible_id": "uuid | null"
}
```

> Все поля кроме `request_id` опциональны. `null` = не менять.

**Response `data`:** `Request`

---


## Pub/Sub события (исходящие)

Сервис **публикует** события без ожидания ответа:

| Subject | Когда | Payload |
|---|---|---|
| `bot-svc.requests.created` | Заявка успешно создана | `domain.Request` (JSON) |
| `bot-svc.requests.canceled` | Заявка отменена | `domain.Request` (JSON) |
| `bot-svc.users.added` | Пользователь сохранён в БД после первого обращения | `domain.User` (JSON) |

---

## Pub/Sub события (входящие, fire-and-forget)

| Subject | Источник | Описание |
|---|---|---|
| `gateway.requests.canceled` | Gateway | Отмена заявки инициирована из внешнего сервиса; обрабатывается как обычная отмена без reply |
| `event.user.created` | Auth Service | Создан новый пользователь — сервис сохраняет его как потенциального ответственного (`SaveResponsible`) |

---

## Коды ошибок в `GatewayResponse`

| `message` | Значение |
|---|---|
| `"invalid request format"` | Не удалось десериализовать тело запроса |
| `"validation error"` | Ошибка валидации входных данных |
| `"not found"` | Ресурс не найден |
| `"internal server error"` | Внутренняя ошибка сервиса |
| `"success"` | Успех (`success: true`) |

