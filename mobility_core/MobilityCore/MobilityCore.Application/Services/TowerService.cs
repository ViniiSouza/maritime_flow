using System.Net.Http.Json;
using System.Text.Json;
using MobilityCore.Shared.Models;

namespace MobilityCore.Application.Services;

public class TowerService
{
    private readonly HttpClient _httpClient;
    private readonly JsonSerializerOptions _jsonOptions;

    public TowerService(HttpClient httpClient)
    {
        _httpClient = httpClient;
        _jsonOptions = new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true
        };
    }

    public async Task<List<Tower>> GetTowersAsync(string towerAddress)
    {
        try
        {
            var response = await _httpClient.GetAsync($"http://{towerAddress}/towers");
            response.EnsureSuccessStatusCode();
            var towersResponse = await response.Content.ReadFromJsonAsync<TowersResponse>(_jsonOptions);
            return towersResponse?.Towers ?? new List<Tower>();
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Erro ao buscar torres: {ex.Message}");
            return new List<Tower>();
        }
    }

    public async Task<StructuresResponse?> GetStructuresAsync(string towerAddress)
    {
        try
        {
            var response = await _httpClient.GetAsync($"http://{towerAddress}/structures");
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadFromJsonAsync<StructuresResponse>(_jsonOptions);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Erro ao buscar estruturas: {ex.Message}");
            return null;
        }
    }

    public async Task<SlotResponse?> RequestSlotAsync(string towerAddress, SlotRequest request)
    {
        try
        {
            var response = await _httpClient.PostAsJsonAsync($"http://{towerAddress}/slots", request, _jsonOptions);
            response.EnsureSuccessStatusCode();
            return await response.Content.ReadFromJsonAsync<SlotResponse>(_jsonOptions);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Erro ao solicitar slot: {ex.Message}");
            return null;
        }
    }
}

