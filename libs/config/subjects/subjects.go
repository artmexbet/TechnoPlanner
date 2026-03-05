package subjects

const (
	UserCreated = "event.user.created"

	CategoryCreated = "gateway.equipment.category.created"
	CategoryUpdated = "gateway.equipment.category.updated"
	CategoryDeleted = "gateway.equipment.category.deleted"

	EquipmentCreated = "gateway.equipment.item.created"
	EquipmentUpdated = "gateway.equipment.item.updated"
	EquipmentDeleted = "gateway.equipment.item.deleted"

	RequestStatusChanged   = "gateway.requests.status.changed"
	RequestAssigned        = "gateway.requests.responsible.assigned"
	GatewayRequestCanceled = "gateway.requests.canceled"

	GatewayRequestCreated = "gateway.requests.created"

	// GatewayRequestList - запрос списка заявок от gateway
	GatewayRequestList = "gateway.requests.list"
	// GatewayRequestGet - запрос заявки по ID от gateway
	GatewayRequestGet = "gateway.requests.get"
	// GatewayRequestAssignResponsible - назначение ответственного от gateway
	GatewayRequestAssignResponsible = "gateway.requests.assign.responsible"
	// GatewayRequestUpdate - обновление заявки от gateway
	GatewayRequestUpdate = "gateway.requests.update"

	// Porter subjects (хранение портеров в Requests сервисе)

	// GatewayPorterList - список всех портеров
	GatewayPorterList = "gateway.porters.list"
	// GatewayPorterGet - получение портера по ID
	GatewayPorterGet = "gateway.porters.get"
	// GatewayPorterDelete - удаление портера
	GatewayPorterDelete = "gateway.porters.delete"
	// GatewayPorterSave - сохранение (upsert) портера
	GatewayPorterSave = "gateway.porters.save"

	// User subjects

	GatewayUserList   = "gateway.user.list"
	GatewayUserGet    = "gateway.user.get"
	GatewayUserCreate = "gateway.user.create"
	GatewayUserUpdate = "gateway.user.update"
	GatewayUserDelete = "gateway.user.delete"

	// History subjects

	GatewayHistoryList = "gateway.history.list"
	GatewayHistoryAdd  = "gateway.history.add"

	// Equipment subjects (новое API Equipment Service)

	GatewayEquipmentList          = "gateway.equipment.list"
	GatewayEquipmentGet           = "gateway.equipment.get"
	GatewayEquipmentCreate        = "gateway.equipment.create"
	GatewayEquipmentUpdate        = "gateway.equipment.update"
	GatewayEquipmentDelete        = "gateway.equipment.delete"
	GatewayEquipmentGetByCategory = "gateway.equipment.get.by_category"

	// Equipment reservation subjects (Request-Reply, equipment service)
	GatewayEquipmentReserve           = "gateway.equipment.reserve"
	GatewayEquipmentRelease           = "gateway.equipment.release"
	GatewayEquipmentCheckAvailability = "gateway.equipment.check"

	// Equipment Category subjects (новое API Equipment Service)

	GatewayEquipmentCategoryList   = "gateway.equipment_category.list"
	GatewayEquipmentCategoryGet    = "gateway.equipment_category.get"
	GatewayEquipmentCategoryCreate = "gateway.equipment_category.create"
	GatewayEquipmentCategoryUpdate = "gateway.equipment_category.update"
	GatewayEquipmentCategoryDelete = "gateway.equipment_category.delete"

	ServiceRequestCreated  = "bot-svc.requests.created"
	ServiceRequestCanceled = "bot-svc.requests.canceled"
	ServiceUserAdded       = "bot-svc.users.added"
	ServiceEquipmentAdd    = "bot-svc.equipment.add"

	SubjectBotRequestCreated = "bot.requests.created"
	SubjectBotRequestGet     = "bot.requests.get"
	SubjectBotRequestList    = "bot.requests.list"
	SubjectBotRequestCancel  = "bot.requests.cancel"

	// Raw request subjects (сырые запросы от бота)

	// SubjectBotRawRequestCreated — бот отправляет сырой запрос
	SubjectBotRawRequestCreated = "bot.raw_requests.created"
	// GatewayRawRequestList — gateway запрашивает список сырых запросов
	GatewayRawRequestList = "gateway.raw_requests.list"
	// GatewayRawRequestGet — gateway запрашивает сырой запрос по ID
	GatewayRawRequestGet = "gateway.raw_requests.get"
	// GatewayRawRequestProcess — gateway запрашивает создание нормальной заявки из сырого запроса
	GatewayRawRequestProcess = "gateway.raw_requests.process"
)
