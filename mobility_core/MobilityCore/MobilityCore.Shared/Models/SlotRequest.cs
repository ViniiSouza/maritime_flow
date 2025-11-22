using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class SlotRequest
{
    [JsonPropertyName("structure_uuid")]
    public string StructureUuid { get; set; } = string.Empty;
    
    [JsonPropertyName("structure_type")]
    public string StructureType { get; set; } = string.Empty; // "platform" | "central"
    
    [JsonPropertyName("slot_number")]
    public int SlotNumber { get; set; }
    
    [JsonPropertyName("slot_type")]
    public string SlotType { get; set; } = string.Empty; // "dock" | "helipad"
    
    [JsonPropertyName("vehicle_type")]
    public string VehicleType { get; set; } = string.Empty; // "helicopter" | "ship"
    
    [JsonPropertyName("vehicle_uuid")]
    public string VehicleUuid { get; set; } = string.Empty;
}

