using System;
using System.Diagnostics;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;
public class WebSocketMessage
{
    public DTO Body;
    public string ActionType;
}

public class WebSocketMessageConverter : JsonConverter
{
    public override bool CanConvert(Type objectType)
    {
        return objectType == typeof(WebSocketMessage);
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
                return new WebSocketMessage { ActionType = actionType, Body = body.ToObject<CharacterDTO>() };
            case "GameState":
                return new WebSocketMessage { ActionType = actionType, Body = body.ToObject<GameStateDTO>() };
            case "Respawn":
                return new WebSocketMessage { ActionType = actionType, Body = body.ToObject<CharacterDTO>() };
            default:
                throw new JsonSerializationException($"Unknown type '{actionType}' in JSON.");
        }
    }

    public override void WriteJson(JsonWriter writer, object value, JsonSerializer serializer)
    {
        throw new NotImplementedException();
    }
}
