using MobilityCore.Application.Services;
using MobilityCore.Shared;
using MobilityCore.Shared.Models;

// fluxo:
// 1 - veículo nasce
// 2 - aguarda algum tempo
// 3 - solicita um slot para ir
// 4 - se permitido, vai para destino. ao acabar, volta para o passo 2
// 5 - se negado, volta para o passo 2

if (args.Length < 3)
{
    Console.WriteLine("Uso: MobilityCore.Application <vehicle_type> <latitude> <longitude> [tower_address]");
    Console.WriteLine("vehicle_type: Ship ou Helicopter");
    Environment.Exit(1);
}

var typeStr = args[0];
var vehicleType = typeStr.Equals("Helicopter", StringComparison.OrdinalIgnoreCase)
    ? VehicleType.Helicopter
    : VehicleType.Ship;

var lat = Convert.ToDouble(args[1]);
var lon = Convert.ToDouble(args[2]);
var towerAddress = args.Length > 3 ? args[3] : "localhost:5000";

Vehicle vehicle = new(vehicleType, lat, lon);
Console.WriteLine($"Veículo criado: UUID={vehicle.Uuid}, Tipo={vehicle.Type}, Posição=({lat}, {lon})");

var httpClient = new HttpClient();
var towerService = new TowerService(httpClient);

const int initialWaitSeconds = 5;
const int retryWaitSeconds = 10;
const int movementIntervalMs = 1000;

while (true)
{
    try
    {
        Console.WriteLine($"Aguardando {initialWaitSeconds} segundos antes de solicitar slot...");
        await Task.Delay(TimeSpan.FromSeconds(initialWaitSeconds));

        Console.WriteLine($"Buscando torres em {towerAddress}...");
        var towers = await towerService.GetTowersAsync(towerAddress);

        if (towers.Count == 0)
        {
            Console.WriteLine("Nenhuma torre encontrada. Tentando novamente...");
            await Task.Delay(TimeSpan.FromSeconds(retryWaitSeconds));
            continue;
        }

        var selectedTower = towers[0];
        Console.WriteLine($"Torre selecionada: {selectedTower.TowerUuid} em {selectedTower.TowerAddress}");

        Console.WriteLine("Buscando estruturas disponíveis...");
        var structuresResponse = await towerService.GetStructuresAsync(selectedTower.TowerAddress);

        if (structuresResponse == null)
        {
            Console.WriteLine("Erro ao buscar estruturas. Tentando novamente...");
            await Task.Delay(TimeSpan.FromSeconds(retryWaitSeconds));
            continue;
        }

        var (structure, structureType, slotNumber, slotType) =
            StructureSelector.SelectStructureAndSlot(structuresResponse, vehicle.Type);

        if (structure == null)
        {
            Console.WriteLine("Nenhuma estrutura disponível com slots adequados. Tentando novamente...");
            await Task.Delay(TimeSpan.FromSeconds(retryWaitSeconds));
            continue;
        }

        var structureUuid = StructureSelector.GetStructureUuid(structure, structureType);
        Console.WriteLine($"Estrutura selecionada: {structureType} {structureUuid}, Slot: {slotType} #{slotNumber}");

        var slotRequest = new SlotRequest
        {
            StructureUuid = structureUuid,
            StructureType = structureType,
            SlotNumber = slotNumber,
            SlotType = slotType,
            VehicleType = vehicle.Type == VehicleType.Helicopter ? "helicopter" : "ship",
            VehicleUuid = vehicle.Uuid
        };

        Console.WriteLine("Solicitando permissão para o slot...");
        var slotResponse = await towerService.RequestSlotAsync(selectedTower.TowerAddress, slotRequest);

        if (slotResponse == null)
        {
            Console.WriteLine("Erro ao solicitar slot. Tentando novamente...");
            await Task.Delay(TimeSpan.FromSeconds(retryWaitSeconds));
            continue;
        }

        if (slotResponse.State == "free")
        {
            Console.WriteLine($"Slot concedido! Movendo-se para {structureType} {structureUuid}...");
            vehicle.Status = StatusMovimento.InTransit;

            var destino = new GeoPoint(structure.Latitude, structure.Longitude);
            bool arrived = false;

            while (!arrived)
            {
                arrived = vehicle.MoveTowardsDestination(destino);
                var distance = GeoHelper.HaversineDistance(vehicle.Position, destino);
                Console.WriteLine($"Posição atual: ({vehicle.Position.Latitude:F6}, {vehicle.Position.Longitude:F6}), " +
                                $"Distância ao destino: {distance:F2}m");

                if (!arrived)
                {
                    await Task.Delay(movementIntervalMs);
                }
            }

            Console.WriteLine("Chegou ao destino!");
            vehicle.Status = StatusMovimento.Parked;

            Console.WriteLine("Aguardando no destino...");
            await Task.Delay(TimeSpan.FromSeconds(30));
            vehicle.Status = StatusMovimento.Stationary;
        }
        else
        {
            Console.WriteLine($"Slot não disponível (state: {slotResponse.State}). Aguardando {retryWaitSeconds} segundos...");
            await Task.Delay(TimeSpan.FromSeconds(retryWaitSeconds));
        }
    }
    catch (Exception ex)
    {
        Console.WriteLine($"Erro no loop principal: {ex.Message}");
        Console.WriteLine($"Stack trace: {ex.StackTrace}");
        await Task.Delay(TimeSpan.FromSeconds(retryWaitSeconds));
    }
}