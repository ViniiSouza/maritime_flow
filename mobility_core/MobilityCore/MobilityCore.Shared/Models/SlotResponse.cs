using System.Text.Json.Serialization;

namespace MobilityCore.Shared.Models;

public class SlotResponse
{
    [JsonPropertyName("state")]
    public string State { get; set; } = string.Empty; // "free" | "in_use"
}

