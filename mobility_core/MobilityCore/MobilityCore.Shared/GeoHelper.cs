using MobilityCore.Shared.Models;

namespace MobilityCore.Shared;

public static class GeoHelper
{
    public static GeoPoint MoveTowardsDestination(GeoPoint origin, GeoPoint destination, double metersPerSecond)
    {
        const double R = 6371000; // raio da Terra em metros

        var lat1 = ToRad(origin.Latitude);
        var lon1 = ToRad(origin.Longitude);
        var lat2 = ToRad(destination.Latitude);
        var lon2 = ToRad(destination.Longitude);

        // Distância atual
        var d = HaversineDistance(origin, destination);

        // Se já está perto o suficiente ou passou do destino → retorna destino direto
        if (d == 0 || metersPerSecond >= d)
            return new GeoPoint(destination.Latitude, destination.Longitude);

        // Direção (bearing)
        var bearing = Math.Atan2(
            Math.Sin(lon2 - lon1) * Math.Cos(lat2),
            Math.Cos(lat1) * Math.Sin(lat2) -
            Math.Sin(lat1) * Math.Cos(lat2) * Math.Cos(lon2 - lon1)
        );

        // Nova posição
        var frac = metersPerSecond / R; // distância angular

        var newLat =
            Math.Asin(Math.Sin(lat1) * Math.Cos(frac) +
                      Math.Cos(lat1) * Math.Sin(frac) * Math.Cos(bearing));

        var newLon =
            lon1 + Math.Atan2(
                Math.Sin(bearing) * Math.Sin(frac) * Math.Cos(lat1),
                Math.Cos(frac) - Math.Sin(lat1) * Math.Sin(newLat)
            );

        // Converte de volta para graus
        return new GeoPoint(ToDeg(newLat), ToDeg(newLon));
    }

    private static double ToRad(double deg) => deg * Math.PI / 180.0;
    private static double ToDeg(double rad) => rad * 180.0 / Math.PI;

    public static double HaversineDistance(GeoPoint p1, GeoPoint p2)
    {
        const double R = 6371000; // metros

        var dLat = ToRad(p2.Latitude - p1.Latitude);
        var dLon = ToRad(p2.Longitude - p1.Longitude);

        var lat1 = ToRad(p1.Latitude);
        var lat2 = ToRad(p2.Latitude);

        var a = Math.Sin(dLat / 2) * Math.Sin(dLat / 2) +
                Math.Cos(lat1) * Math.Cos(lat2) *
                Math.Sin(dLon / 2) * Math.Sin(dLon / 2);

        var c = 2 * Math.Atan2(Math.Sqrt(a), Math.Sqrt(1 - a));

        return R * c;
    }
}