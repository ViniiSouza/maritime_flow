namespace MobilityCore.Shared.Models;

public class Vehicle
{
    public string Uuid { get; set; }
    public VehicleType Type { get; set; }
    public double Velocity { get; set; }
    public GeoPoint Position { get; set; }
    public StatusMovimento Status { get; set; }

    public Vehicle(VehicleType type, double lat, double lon, double vel, string? uuid = null)
    {
        Uuid = uuid ?? Guid.NewGuid().ToString("N");
        Type = type;
        Velocity = vel;
        Position = new GeoPoint(lat, lon);
        Status = StatusMovimento.Stationary;
    }

    /// <summary>
    /// Move its own velocity in meters towards the destination
    /// </summary>
    /// <param name="destination"></param>
    /// <returns>If arrived at the destination, returns true. Otherwise, returns false.</returns>
    public bool MoveTowardsDestination(GeoPoint destination)
    {
        Position = GeoHelper.MoveTowardsDestination(Position, destination, Velocity);
        return GeoHelper.HaversineDistance(Position, destination) == 0d;
    }
}