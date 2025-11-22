package types

var vehicleSlotMapping = map[VehicleType]SlotType{
	HelicopterVehicleType: HelipadSlotType,
	ShipVehicleType: DockSlotType,
} 

func GetSlotTypeByVehicleType(vehicle VehicleType) SlotType {
	return vehicleSlotMapping[vehicle]
}
