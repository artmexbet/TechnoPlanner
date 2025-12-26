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

	ServiceRequestCreated  = "bot-svc.requests.created"
	ServiceRequestCanceled = "bot-svc.requests.canceled"
	ServiceUserAdded       = "bot-svc.users.added"
	ServiceEquipmentAdd    = "bot-svc.equipment.add"

	SubjectBotRequestCreated = "bot.requests.created"
	SubjectBotRequestGet     = "bot.requests.get"
	SubjectBotRequestList    = "bot.requests.list"
	SubjectBotRequestCancel  = "bot.requests.cancel"
)
