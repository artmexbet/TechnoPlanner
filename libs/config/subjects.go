package config

const (
	SubjectUserCreated = "event.user.created"

	SubjectCategoryCreated = "gateway.equipment.category.created"
	SubjectCategoryUpdated = "gateway.equipment.category.updated"
	SubjectCategoryDeleted = "gateway.equipment.category.deleted"

	SubjectEquipmentCreated = "gateway.equipment.item.created"
	SubjectEquipmentUpdated = "gateway.equipment.item.updated"
	SubjectEquipmentDeleted = "gateway.equipment.item.deleted"

	SubjectRequestStatusChanged = "gateway.requests.status.changed"
	SubjectRequestAssigned      = "gateway.requests.responsible.assigned"

	SubjectGatewayRequestCreated = "gateway.requests.created"
	SubjectBotRequestCreated     = "bot.requests.created"
)
