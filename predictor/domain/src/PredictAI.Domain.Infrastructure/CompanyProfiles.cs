using System.Text.Json;
using PredictAI.Domain.Core;

namespace PredictAI.Domain.Infrastructure;

public static class CompanyProfiles
{
    public static void Apply(Company company, JsonElement profile)
    {
        if (profile.ValueKind != JsonValueKind.Object || !profile.EnumerateObject().Any())
            throw new ArgumentException("profile must be a non-empty JSON object");
        var json = JsonSerializer.Serialize(profile);
        if (System.Text.Encoding.UTF8.GetByteCount(json) > 1_000_000)
            throw new ArgumentException("profile exceeds 1 MB");
        company.ProfileJson = json;
    }

    public static object ToResponse(Company c) => new
    {
        companyId = c.Id, c.ProfileRevision, c.ValeraDeliveryStatus,
        profile = JsonDocument.Parse(c.ProfileJson).RootElement.Clone()
    };

    public static JsonElement ToProfile(Company c) => JsonDocument.Parse(c.ProfileJson).RootElement.Clone();
}
