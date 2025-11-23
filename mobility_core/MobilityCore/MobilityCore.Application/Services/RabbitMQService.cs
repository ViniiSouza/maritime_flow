using System.Text;
using System.Text.Json;
using RabbitMQ.Client;
using MobilityCore.Shared.Models;

namespace MobilityCore.Application.Services;

public class RabbitMQService : IDisposable
{
    private readonly IConnection _connection;
    private readonly IModel _channel;
    private readonly JsonSerializerOptions _jsonOptions;
    private readonly string _metricsQueue;
    private readonly string _auditQueue;

    public RabbitMQService(string host, int port, string username, string password, string metricsQueue = "metrics", string auditQueue = "audit")
    {
        _metricsQueue = metricsQueue;
        _auditQueue = auditQueue;
        
        _jsonOptions = new JsonSerializerOptions
        {
            PropertyNamingPolicy = null
        };

        var factory = new ConnectionFactory
        {
            HostName = host,
            Port = port,
            UserName = username,
            Password = password,
            RequestedHeartbeat = TimeSpan.FromSeconds(60),
            RequestedConnectionTimeout = TimeSpan.FromSeconds(30),
            AutomaticRecoveryEnabled = true,
            NetworkRecoveryInterval = TimeSpan.FromSeconds(10)
        };

        _connection = factory.CreateConnection();
        _channel = _connection.CreateModel();

        _channel.QueueDeclare(queue: _metricsQueue, durable: false, exclusive: false, autoDelete: false, arguments: null);
        _channel.QueueDeclare(queue: _auditQueue, durable: false, exclusive: false, autoDelete: false, arguments: null);
    }

    public void PublishMetrics(MetricsMessage message)
    {
        try
        {
            var json = JsonSerializer.Serialize(message, _jsonOptions);
            var body = Encoding.UTF8.GetBytes(json);

            _channel.BasicPublish(
                exchange: "",
                routingKey: _metricsQueue,
                basicProperties: null,
                body: body
            );

        }
        catch (Exception ex)
        {
            Console.WriteLine($"Erro ao publicar métricas: {ex.Message}");
        }
    }

    public void PublishAudit(AuditMessage message)
    {
        try
        {
            var json = JsonSerializer.Serialize(message, _jsonOptions);
            var body = Encoding.UTF8.GetBytes(json);

            _channel.BasicPublish(
                exchange: "",
                routingKey: _auditQueue,
                basicProperties: null,
                body: body
            );

        }
        catch (Exception ex)
        {
            Console.WriteLine($"Erro ao publicar evento de audit: {ex.Message}");
        }
    }

    public void Dispose()
    {
        _channel?.Close();
        _channel?.Dispose();
        _connection?.Close();
        _connection?.Dispose();
    }
}

