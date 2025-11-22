package types

var vehicleSlotMapping = map[VehicleType]SlotType{
	HelicopterVehicleType: HelipadSlotType,
	ShipVehicleType:       DockSlotType,
}

var slotResultMapping = map[SlotState]ResultType{
	FreeSlotState:  AllowedResultType,
	InUseSlotState: DeniedResultType,
}

func GetSlotTypeByVehicleType(vehicle VehicleType) SlotType {
	return vehicleSlotMapping[vehicle]
}

func GetResultTypeBySlotState(state SlotState) ResultType {
	return slotResultMapping[state]
}
