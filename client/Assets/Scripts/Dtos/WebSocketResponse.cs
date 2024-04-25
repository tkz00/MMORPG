using System;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
public class WebSocketResponse {
	public DTO Body;
	public string ActionType;
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
        var actionType = jsonObject["actionType"]?.Value<string>();

        if (actionType == null)
        {
            throw new JsonSerializationException("Missing 'Type' field in JSON.");
        }

        var body = jsonObject["body"];

        switch (actionType)
        {
            case "Player":
                return new WebSocketResponse { ActionType = actionType, Body = body.ToObject<PlayerDTO>() };
            case "GameState":
                return new WebSocketResponse { ActionType = actionType, Body = body.ToObject<GameStateDTO>() };
            default:
                throw new JsonSerializationException($"Unknown type '{actionType}' in JSON.");
        }
    }

    public override void WriteJson(JsonWriter writer, object value, JsonSerializer serializer)
    {
        throw new NotImplementedException();
    }
}