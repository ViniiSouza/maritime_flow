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
    Console.WriteLine("Uso: MobilityCore.Application <vehicle_uuid> <vehicle_type> <latitude> <longitude> [towers_discovery_address]");
    Console.WriteLine("vehicle_type: Ship ou Helicopter");
    Console.WriteLine("\nVariáveis de ambiente opcionais:");
    Console.WriteLine("  BASE_DNS - DNS base para concatenar com o Id da Torre (padrão: tower.svc.cluster.local)");
    Console.WriteLine("  RABBITMQ_HOST - Host do RabbitMQ (padrão: localhost)");
    Console.WriteLine("  RABBITMQ_PORT - Porta do RabbitMQ (padrão: 5672)");
    Console.WriteLine("  RABBITMQ_USERNAME - Usuário do RabbitMQ (padrão: guest)");
    Console.WriteLine("  RABBITMQ_PASSWORD - Senha do RabbitMQ (padrão: guest)");
    Environment.Exit(1);
}

var uuid = args[0];
var typeStr = args[1];
var vehicleType = typeStr.Equals("Helicopter", StringComparison.OrdinalIgnoreCase)
    ? VehicleType.Helicopter
    : VehicleType.Ship;

var lat = Convert.ToDouble(args[2]);
var lon = Convert.ToDouble(args[3]);
var towersDiscoveryAddress = args.Length > 4 ? args[4] : "towers-svc.tower.svc.cluster.local";

Vehicle vehicle = new(vehicleType, lat, lon, uuid);
Console.WriteLine($"Veículo criado: UUID={vehicle.Uuid}, Tipo={vehicle.Type}, Posição=({lat}, {lon})");

var httpClient = new HttpClient();
var towerService = new TowerService(httpClient);

var baseDns = Environment.GetEnvironmentVariable("BASE_DNS") ?? "tower.svc.cluster.local";
var rabbitmqHost = Environment.GetEnvironmentVariable("RABBITMQ_HOST") ?? "localhost";
var rabbitmqPortStr = Environment.GetEnvironmentVariable("RABBITMQ_PORT") ?? "5672";
if (!int.TryParse(rabbitmqPortStr, out var rabbitmqPort))
{
    Console.WriteLine($"ERRO: Porta do RabbitMQ inválida: {rabbitmqPortStr}");
    Environment.Exit(1);
}
var rabbitmqUsername = Environment.GetEnvironmentVariable("RABBITMQ_USERNAME") ?? "guest";
var rabbitmqPassword = Environment.GetEnvironmentVariable("RABBITMQ_PASSWORD") ?? "guest";

RabbitMQService? rabbitMQService = null;
try
{
    rabbitMQService = new RabbitMQService(rabbitmqHost, rabbitmqPort, rabbitmqUsername, rabbitmqPassword);
}
catch (Exception ex)
{
    Console.WriteLine($"Aviso: Não foi possível conectar ao RabbitMQ: {ex.Message}");
    Console.WriteLine("A aplicação continuará sem publicar métricas e eventos de audit.");
}

AppDomain.CurrentDomain.ProcessExit += (sender, e) =>
{
    rabbitMQService?.Dispose();
};

const int initialWaitSeconds = 5;
const int retryWaitSeconds = 10;
const int movementIntervalMs = 1000;
const int metricsIntervalMs = 2000;

var random = new Random();

while (true)
{
    try
    {
        Console.WriteLine($"Aguardando {initialWaitSeconds} segundos antes de solicitar slot...");
        await Task.Delay(TimeSpan.FromSeconds(initialWaitSeconds));

        Console.WriteLine($"Buscando torres em {towersDiscoveryAddress}...");
        var towers = await towerService.GetTowersAsync(towersDiscoveryAddress);

        if (towers.Count == 0)
        {
            Console.WriteLine("Nenhuma torre encontrada. Tentando novamente...");
            await Task.Delay(TimeSpan.FromSeconds(retryWaitSeconds));
            continue;
        }

        var selectedTower = towers[0];
        var towerAddress = $"{selectedTower.TowerUuid}.{baseDns}";
        Console.WriteLine($"Torre selecionada: {selectedTower.TowerUuid} em {towerAddress}");

        Console.WriteLine("Buscando estruturas disponíveis...");
        var structuresResponse = await towerService.GetStructuresAsync(towerAddress);

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
        var slotResponse = await towerService.RequestSlotAsync(towerAddress, slotRequest);

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

            if (rabbitMQService != null)
            {
                var departedEvent = new AuditMessage
                {
                    VehicleType = vehicle.Type == VehicleType.Helicopter ? "helicopter" : "ship",
                    VehicleUuid = vehicle.Uuid,
                    StructureType = structureType,
                    StructureUuid = structureUuid,
                    Timestamp = DateTimeOffset.UtcNow.ToUnixTimeSeconds(),
                    Event = "departed",
                    SlotNumber = slotNumber,
                    TowerId = selectedTower.TowerUuid
                };
                rabbitMQService.PublishAudit(departedEvent);
            }

            var destino = new GeoPoint(structure.Latitude, structure.Longitude);
            bool arrived = false;
            var lastMetricsTime = DateTime.UtcNow;

            while (!arrived)
            {
                arrived = vehicle.MoveTowardsDestination(destino);
                var distance = GeoHelper.HaversineDistance(vehicle.Position, destino);
                Console.WriteLine($"Posição atual: ({vehicle.Position.Latitude:F6}, {vehicle.Position.Longitude:F6}), " +
                                $"Distância ao destino: {distance:F2}m");

                if (rabbitMQService != null && (DateTime.UtcNow - lastMetricsTime).TotalMilliseconds >= metricsIntervalMs)
                {
                    var metrics = GenerateMetrics(vehicle, random);
                    rabbitMQService.PublishMetrics(metrics);
                    lastMetricsTime = DateTime.UtcNow;
                }

                if (!arrived)
                {
                    await Task.Delay(movementIntervalMs);
                }
            }

            Console.WriteLine("Chegou ao destino!");
            vehicle.Status = StatusMovimento.Parked;

            if (rabbitMQService != null)
            {
                var arrivedEvent = new AuditMessage
                {
                    VehicleType = vehicle.Type == VehicleType.Helicopter ? "helicopter" : "ship",
                    VehicleUuid = vehicle.Uuid,
                    StructureType = structureType,
                    StructureUuid = structureUuid,
                    Timestamp = DateTimeOffset.UtcNow.ToUnixTimeSeconds(),
                    Event = "arrived",
                    SlotNumber = slotNumber,
                    TowerId = selectedTower.TowerUuid
                };
                rabbitMQService.PublishAudit(arrivedEvent);
            }

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

static MetricsMessage GenerateMetrics(Vehicle vehicle, Random random)
{
    var fuelLevel = Math.Max(0.5, 1.0 - (random.NextDouble() * 0.5));
    var baseTemp = vehicle.Type == VehicleType.Helicopter ? 40.0 : 35.0;
    var temperature = baseTemp + (random.NextDouble() * 20.0); // 35-55 para navio, 40-60 para helicóptero
    
    // CPU e memória variam entre 0.2 e 0.8
    var cpuUsage = 0.2 + (random.NextDouble() * 0.6);
    var memUsage = 0.2 + (random.NextDouble() * 0.6);
    
    // Memória em bytes (simula entre 2MB e 8MB)
    var memUsageBytes = (long)(2_000_000 + (random.NextDouble() * 6_000_000));

    return new MetricsMessage
    {
        Latitude = vehicle.Position.Latitude,
        Longitude = vehicle.Position.Longitude,
        FuelLevel = fuelLevel,
        Temperature = temperature,
        CpuUsage = cpuUsage,
        MemUsage = memUsage,
        MemUsageBytes = memUsageBytes
    };
}