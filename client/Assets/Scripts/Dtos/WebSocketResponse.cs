using System;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
public class WebSocketResponse {
	public DTO Body;
	public string Type;
}

public class WebSocketResponseConverter : JsonConverter
{
    public override bool CanConvert(Type objectType)
    {
        return objectType == typeof(WebSocketResponse);
    }

    public override object ReadJson(JsonReader reader, Type objectType, object existingValue, JsonSerializer serializer)
    {
        var jsonObject = JObject.Load(reader);
        var type = jsonObject["type"]?.Value<string>();

        if (type == null)
        {
            throw new JsonSerializationException("Missing 'Type' field in JSON.");
        }

        var body = jsonObject["body"];

        switch (type)
        {
            case "PlayerDTO":
                return new WebSocketResponse { Type = type, Body = body.ToObject<PlayerDTO>() };
            case "GameStateDTO":
                return new WebSocketResponse { Type = type, Body = body.ToObject<GameStateDTO>() };
            default:
                throw new JsonSerializationException($"Unknown type '{type}' in JSON.");
        }
    }

    public override void WriteJson(JsonWriter writer, object value, JsonSerializer serializer)
    {
        throw new NotImplementedException();
    }
}