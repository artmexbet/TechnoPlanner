**План Реализации CRUD в gateway**

## Plan: CRUD porter/equipment/requests

Создаем миграции для ролей, porter-метаданных, категорий и истории статусов, обновляем sqlc-модели и domain-структуры,
расширяем storage/service/router слоями для porter-ограниченного CRUD и RBAC, обеспечиваем soft-delete/audit-поля и
связь equipment↔categories, добавляем отдельные эндпоинты и OpenAPI-описания, плюс seed-данные и настройки
фильтрации/истории в requests.

### Steps

1. Спроектировать новые таблицы в `services/gateway/migrations` (porter флаги в users, audit/soft-delete,
   equipment_categories, equipment_category_links, request_status_history, role seed) и описать seed-скрипты.
2. Обновить `services/gateway/internal/postgres/queries/*.sql` и `sqlc.yaml` + прогнать sqlc, чтобы сгенерировать модели
   и CRUD-запросы для equipment, categories, requests, history и porter-фильтров.
3. Расширить `services/gateway/internal/domain/models.go` и `internal/models` DTO для новых сущностей, включая массив
   категорий и audit-поля.
4. Доработать `internal/storage` и `internal/service` слоя: CRUD-логика equipment/categories, soft-delete, porter
   фильтры, назначение ответственных, история статусов, RBAC проверки (admin-only).
5. Обновить `internal/router` (маршруты + middleware RBAC-checks), `internal/router/middlwares` (ролевая проверка) и
   `docs/openapi.yaml` для новых эндпоинтов, включая отдельный эндпоинт истории request-статусов.
6. Добавить seed-механизм (через миграцию или init-скрипт) для ролей admin/porter и убедиться, что Make/README описывают
   запуск обновлённых миграций и sqlc.

### Further Considerations

1. Уточнить, где хранится связь user↔role: использовать существующую таблицу roles или добавить user_roles?
2. Требуется ли версионирование audit-полей для чужих сервисов (например, capture service ID) или достаточно user_id?
3. Нужно ли покрыть новые сервисные методы e2e/интеграционными тестами в `services/gateway/tests` или достаточно
   unit-тестов?

## Progress
- [x] migrations updated
- [x] sqlc models regenerated and domain/DTO synced
- [x] postgres adapters for equipment/categories/requests/status history/porters
- [x] storage layer with publisher hooks
- [x] service layer (porter/equipment/category/request/history + RBAC)
- [ ] router layer & OpenAPI