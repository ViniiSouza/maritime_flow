using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class Tower
{
    [JsonPropertyName("tower_uuid")]
    public string TowerUuid { get; set; } = string.Empty;
    
    [JsonPropertyName("latitude")]
    public double Latitude { get; set; }
    
    [JsonPropertyName("longitude")]
    public double Longitude { get; set; }
}

