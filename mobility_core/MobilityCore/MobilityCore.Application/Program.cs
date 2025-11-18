using MobilityCore.Shared;
using MobilityCore.Shared.Models;

// fluxo:
// 1 - helicóptero nasce
// 2 - aguarda algum tempo
// 3 - solicita um slot para ir
// 4 - se permitido, vai para destino. ao acabar, volta para o passo 2
// 5 - se negado, volta para o passo 2

var type = VehicleType.Ship; // args[0]
var lat = Convert.ToDouble(args[1]);
var lon = Convert.ToDouble(args[2]);

Vehicle vehicle = new (type, lat, lon);

// espera algum tempo

// solicita permissão para uma plataforma/central e slot válido
 
# region se permitido: 
// obtém dados do destino (lat e long)
var destino = new GeoPoint(56.62528, 12.81786);
// começa a se movimentar
bool arrived = false;
while (!arrived)
{
    // sleep
    arrived = vehicle.MoveTowardsDestination(destino);
}
#endregion
# region se não permitido
// reagenda solicitação
#endregion