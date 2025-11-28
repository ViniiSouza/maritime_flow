using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class MetricsMessage
{
    [JsonPropertyName("metadata")]
    public required Metadata Metadata { get; set; }

    [JsonPropertyName("metrics")]
    public required Metrics Metrics { get; set; }
}

