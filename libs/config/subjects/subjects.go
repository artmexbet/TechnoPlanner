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

	// Equipment subjects
	GatewayEquipmentList   = "gateway.equipment.list"
	GatewayEquipmentGet    = "gateway.equipment.get"
	GatewayEquipmentCreate = "gateway.equipment.create"
	GatewayEquipmentUpdate = "gateway.equipment.update"
	GatewayEquipmentDelete = "gateway.equipment.delete"

	// Category subjects
	GatewayCategoryList   = "gateway.category.list"
	GatewayCategoryGet    = "gateway.category.get"
	GatewayCategoryCreate = "gateway.category.create"
	GatewayCategoryUpdate = "gateway.category.update"
	GatewayCategoryDelete = "gateway.category.delete"

	// User subjects
	GatewayUserList   = "gateway.user.list"
	GatewayUserGet    = "gateway.user.get"
	GatewayUserCreate = "gateway.user.create"

	// History subjects
	GatewayHistoryList = "gateway.history.list"
	GatewayHistoryAdd  = "gateway.history.add"

	ServiceRequestCreated  = "bot-svc.requests.created"
	ServiceRequestCanceled = "bot-svc.requests.canceled"
	ServiceUserAdded       = "bot-svc.users.added"
	ServiceEquipmentAdd    = "bot-svc.equipment.add"

	SubjectBotRequestCreated = "bot.requests.created"
	SubjectBotRequestGet     = "bot.requests.get"
	SubjectBotRequestList    = "bot.requests.list"
	SubjectBotRequestCancel  = "bot.requests.cancel"
)
