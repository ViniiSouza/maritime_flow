using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class MetricsMessage
{
    [JsonPropertyName("latitude")]
    public double Latitude { get; set; }

    [JsonPropertyName("longitude")]
    public double Longitude { get; set; }

    [JsonPropertyName("fuel_level")]
    public double FuelLevel { get; set; }

    [JsonPropertyName("temperature")]
    public double Temperature { get; set; }

    [JsonPropertyName("cpu_usage")]
    public double CpuUsage { get; set; }

    [JsonPropertyName("mem_usage")]
    public double MemUsage { get; set; }

    [JsonPropertyName("mem_usage_bytes")]
    public long MemUsageBytes { get; set; }
}

