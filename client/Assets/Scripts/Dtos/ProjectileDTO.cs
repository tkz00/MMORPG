using System.Runtime.Serialization;
using Newtonsoft.Json;
using Newtonsoft.Json.Converters;

public class ProjectileDTO : DTO
{
   public string id;
   public string caster;
   public PositionDTO position;
   public float radius;

   [JsonProperty("state")]
   [JsonConverter(typeof(StringEnumConverter))]
   public State state;
}

public enum State
{
   [EnumMember(Value = "active")]
   Active,
   [EnumMember(Value = "hit")]
   Hit
}
