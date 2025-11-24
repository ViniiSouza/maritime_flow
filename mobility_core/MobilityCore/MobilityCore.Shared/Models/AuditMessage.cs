using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class AuditMessage
{
    [JsonPropertyName("vehicle_type")]
    public string VehicleType { get; set; } = string.Empty;

    [JsonPropertyName("vehicle_uuid")]
    public string VehicleUuid { get; set; } = string.Empty;

    [JsonPropertyName("structure_type")]
    public string StructureType { get; set; } = string.Empty;

    [JsonPropertyName("structure_uuid")]
    public string StructureUuid { get; set; } = string.Empty;

    [JsonPropertyName("timestamp")]
    public long Timestamp { get; set; }

    [JsonPropertyName("event")]
    public string Event { get; set; } = string.Empty; // "arrived" | "departed"

    [JsonPropertyName("slot_number")]
    public int SlotNumber { get; set; }

    [JsonPropertyName("tower_id")]
    public string TowerId { get; set; } = string.Empty;
}

