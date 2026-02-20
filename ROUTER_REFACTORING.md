# Рефакторинг Router - Разделение на модули

## Описание

Большой файл `router.go` (749 строк) был разделен на несколько логически связанных модулей для улучшения навигации и поддержки кода.

## Структура файлов

### До рефакторинга:
```
router/
  ├── router.go (749 строк - все обработчики)
  ├── user.go
  ├── user-protected.go
  └── middlwares/
```

### После рефакторинга:
```
router/
  ├── router.go (166 строк - основная структура, интерфейсы, middleware)
  ├── porter_handlers.go (78 строк - обработчики для porters)
  ├── equipment_handlers.go (233 строк - обработчики для equipment и categories)
  ├── request_handlers.go (247 строк - обработчики для requests и history)
  ├── responsible_handlers.go (63 строки - обработчики для responsibles)
  ├── helpers.go (27 строк - вспомогательные функции)
  ├── user.go (163 строки - аутентификация пользователей)
  ├── user-protected.go (34 строки - защищенные роуты пользователей)
  └── middlwares/
```

## Детали разделения

### 1. router.go
**Содержание:**
- Интерфейсы сервисов (UserService, AuthSvcConnector, PorterService, EquipmentService, CategoryService, RequestService, HistoryService, ResponsibleService)
- Структура Router и ее поля
- Конструктор NewRouter()
- Инициализация middleware (InitMiddlewares)
- Базовые роуты (InitBaseRoutes)
- Методы Run() и userContext()

**Ответственность:** Основная структура роутера и общая конфигурация

### 2. porter_handlers.go
**Содержание:**
- InitPorterRoutes() - регистрация роутов
- listPorters() - GET /api/v1/porters
- getPorter() - GET /api/v1/porters/:id
- createPorter() - POST /api/v1/porters
- toPorterResponse() - маппер для ответов

**Ответственность:** Управление грузчиками (porters)

### 3. equipment_handlers.go
**Содержание:**
- InitEquipmentRoutes() - регистрация роутов для equipment и categories
- Обработчики equipment:
  - listEquipment() - GET /api/v1/equipment
  - getEquipment() - GET /api/v1/equipment/:id
  - createEquipment() - POST /api/v1/equipment
  - updateEquipment() - PUT /api/v1/equipment/:id
  - deleteEquipment() - DELETE /api/v1/equipment/:id
- Обработчики categories:
  - listCategories() - GET /api/v1/equipment/categories
  - createCategory() - POST /api/v1/equipment/categories
  - updateCategory() - PUT /api/v1/equipment/categories/:id
  - deleteCategory() - DELETE /api/v1/equipment/categories/:id
- Мапперы: toEquipmentResponse(), toCategoryResponse()

**Ответственность:** Управление оборудованием и его категориями

### 4. request_handlers.go
**Содержание:**
- InitRequestRoutes() - регистрация роутов
- Обработчики requests:
  - listRequests() - GET /api/v1/requests
  - getRequest() - GET /api/v1/requests/:id
  - updateRequest() - PATCH /api/v1/requests/:id
  - assignResponsible() - POST /api/v1/requests/:id/responsible
- Обработчики history:
  - listRequestHistory() - GET /api/v1/requests/:id/history
  - addRequestHistory() - POST /api/v1/requests/:id/history
- Мапперы: toRequestResponse(), toHistoryResponse()

**Ответственность:** Управление заявками и их историей

### 5. responsible_handlers.go
**Содержание:**
- InitResponsibleRoutes() - регистрация роутов
- listResponsibles() - GET /api/v1/responsibles
- createResponsible() - POST /api/v1/responsibles

**Ответственность:** Управление ответственными пользователями

### 6. helpers.go
**Содержание:**
- handleServiceError() - обработка ошибок сервисов
- derefString() - безопасное разыменование строк

**Ответственность:** Вспомогательные функции общего назначения

### 7. user.go (уже существовал)
**Содержание:**
- InitUserRoutes() - регистрация роутов
- RegisterUser() - POST /api/v1/users/register
- LoginUser() - POST /api/v1/users/login
- RefreshToken() - POST /api/v1/users/refresh
- Me() - GET /api/v1/users/me

**Ответственность:** Аутентификация и регистрация пользователей

### 8. user-protected.go (уже существовал)
**Содержание:**
- InitProtectedUserRoutes() - регистрация защищенных роутов
- LogoutUser() - POST /api/v1/users/logout

**Ответственность:** Защищенные операции пользователей

## Преимущества рефакторинга

1. **Улучшенная навигация**: Теперь легко найти нужный обработчик по имени файла
2. **Логическое разделение**: Каждый файл отвечает за свою предметную область
3. **Упрощенная поддержка**: Изменения в одной области не влияют на код других областей
4. **Лучшая читаемость**: Файлы стали короче и понятнее
5. **Легче работать в команде**: Меньше конфликтов при слиянии изменений

## Проверка

✅ Компиляция проходит успешно  
✅ Все импорты корректны  
✅ Структура Router и её методы доступны во всех файлах  
✅ Сохранена обратная совместимость с существующим кодом

