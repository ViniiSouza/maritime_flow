using MobilityCore.Shared;
using MobilityCore.Shared.Models;

namespace MobilityCore.Application.Services;

public class StructureSelector
{
    public static (Structure? structure, string structureType, int slotNumber, string slotType) SelectStructureAndSlot(
        StructuresResponse structures,
        VehicleType vehicleType,
        string currentStructureUuid)
    {
        var slotType = vehicleType == VehicleType.Helicopter ? "helipad" : "dock";
        var allStructures = new List<(Structure structure, string type)>();

        foreach (var central in structures.Centrals)
        {
            allStructures.Add((central, "central"));
        }

        foreach (var platform in structures.Platforms)
        {
            allStructures.Add((platform, "platform"));
        }

        var availableStructures = allStructures.Where(s =>
        {
            var slots = s.structure.Slots;
            return slotType == "helipad" ? slots.HelipadsQtt > 0 : slots.DocksQtt > 0;
        }).ToList();

        if (availableStructures.Count == 0)
        {
            return (null, string.Empty, -1, string.Empty);
        }

        var random = new Random();
        (Structure structure, string type) selected;
        do
        {
            selected = availableStructures[random.Next(availableStructures.Count)];
        } while ((selected.type == "platform" ? selected.structure.PlatformUuid : selected.structure.CentralUuid) == currentStructureUuid);

        var maxSlots = slotType == "helipad"
            ? selected.structure.Slots.HelipadsQtt
            : selected.structure.Slots.DocksQtt;
        var slotNumber = random.Next(1, maxSlots);

        return (selected.structure, selected.type, slotNumber, slotType);
    }

    public static string GetStructureUuid(Structure structure, string structureType)
    {
        return structureType == "central" ? structure.CentralUuid : structure.PlatformUuid;
    }
}

