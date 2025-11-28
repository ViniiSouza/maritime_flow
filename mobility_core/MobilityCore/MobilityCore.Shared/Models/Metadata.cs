using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class Metadata
{
    [JsonPropertyName("vehicle_uuid")]
    public required string VehicleUuid { get; set; }

    [JsonPropertyName("vehicle_type")]
    public required string VehicleType { get; set; }
}
