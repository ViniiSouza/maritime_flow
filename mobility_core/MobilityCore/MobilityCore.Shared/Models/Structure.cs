using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class Structure
{
    [JsonPropertyName("central_uuid")]
    public string CentralUuid { get; set; } = string.Empty;
    
    [JsonPropertyName("platform_uuid")]
    public string PlatformUuid { get; set; } = string.Empty;
    
    [JsonPropertyName("latitude")]
    public double Latitude { get; set; }
    
    [JsonPropertyName("longitude")]
    public double Longitude { get; set; }
    
    [JsonPropertyName("slots")]
    public Slots Slots { get; set; } = new();
}

public class Slots
{
    [JsonPropertyName("docks_qtt")]
    public int DocksQtt { get; set; }
    
    [JsonPropertyName("helipads_qtt")]
    public int HelipadsQtt { get; set; }
}

