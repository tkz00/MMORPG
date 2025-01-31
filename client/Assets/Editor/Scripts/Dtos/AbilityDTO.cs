using System;
using Newtonsoft.Json;

namespace Configurator
{
    public class AbilityDTO
    {
        public string id;
        public string name;
        public int cooldown;
        public float range;
        public AbilityParameters targeting;
        [JsonProperty("character_state")] public CharacterAction characterAction;
        public MechanicDTO[] mechanics;

        [JsonConverter(typeof(MechanicDTOConverter))]
        public class MechanicDTO
        {
            [JsonProperty("mechanic_type")] public string mechanicType;
            public Params @params;

            public MechanicDTO DeepCopy()
            {
                return new MechanicDTO
                {
                    mechanicType = this.mechanicType,
                    @params = this.@params.DeepCopy()
                };
            }
        }

        public abstract class Params
        {
            public abstract Params DeepCopy();
        }

        public class CreateProjectileParams : Params
        {
            [JsonProperty("on_hit_mechanics")] public MechanicDTO[] onHitMechanics;

            public override Params DeepCopy()
            {
                MechanicDTO[] onHitMechanicsCopy = new MechanicDTO[this.onHitMechanics.Length];
                for (int i = 0; i < this.onHitMechanics.Length; i++)
                {
                    onHitMechanicsCopy[i] = this.onHitMechanics[i].DeepCopy();
                }
                return new CreateProjectileParams
                {
                    onHitMechanics = onHitMechanicsCopy
                };
            }
        }

        public class CreateAoEParams : Params
        {
            [JsonProperty("duration_ms")] public int durationMs;
            public float radius;
            [JsonProperty("on_hit_mechanics")] public MechanicDTO[] onHitMechanics;

            public override Params DeepCopy()
            {
                MechanicDTO[] onHitMechanicsCopy = new MechanicDTO[this.onHitMechanics.Length];
                for (int i = 0; i < this.onHitMechanics.Length; i++)
                {
                    onHitMechanicsCopy[i] = this.onHitMechanics[i].DeepCopy();
                }
                return new CreateAoEParams
                {
                    durationMs = this.durationMs,
                    radius = this.radius,
                    onHitMechanics = onHitMechanicsCopy
                };
            }
        }

        public class DamageParams : Params
        {
            public int amount;
            [JsonProperty("targeting_strategy")] public string targetingStrategy;

            public override Params DeepCopy()
            {
                return new DamageParams
                {
                    amount = this.amount,
                    targetingStrategy = this.targetingStrategy
                };
            }
        }

        public class DelayParams : Params
        {
            [JsonProperty("delay_ms")] public int delayMs;
            [JsonProperty("execute_after_delay_mechanics")] public MechanicDTO[] executeAfterDelayMechanics;

            public override Params DeepCopy()
            {
                return new DelayParams
                {
                    delayMs = this.delayMs,
                    executeAfterDelayMechanics = this.executeAfterDelayMechanics
                };
            }
        }
    }

    public class MechanicDTOConverter : JsonConverter<AbilityDTO.MechanicDTO>
    {
        public override AbilityDTO.MechanicDTO ReadJson(JsonReader reader, Type objectType, AbilityDTO.MechanicDTO existingValue, bool hasExistingValue, JsonSerializer serializer)
        {
            var jsonObject = Newtonsoft.Json.Linq.JObject.Load(reader);
            var mechanicType = jsonObject["mechanic_type"]?.ToString();

            if (string.IsNullOrWhiteSpace(mechanicType))
            {
                throw new JsonSerializationException("Missing or empty mechanic_type.");
            }

            var mechanic = new AbilityDTO.MechanicDTO
            {
                mechanicType = mechanicType
            };

            // Deserialize the "params" field based on the mechanic_type
            mechanic.@params = mechanicType switch
            {
                "create_projectile" => jsonObject["params"]?.ToObject<AbilityDTO.CreateProjectileParams>(serializer),
                "create_AoE" => jsonObject["params"]?.ToObject<AbilityDTO.CreateAoEParams>(serializer),
                "damage" => jsonObject["params"]?.ToObject<AbilityDTO.DamageParams>(serializer),
                "delay" => jsonObject["params"]?.ToObject<AbilityDTO.DelayParams>(serializer),
                _ => throw new JsonSerializationException($"Unknown mechanic_type: {mechanicType}")
            };

            return mechanic;
        }

        public override void WriteJson(JsonWriter writer, AbilityDTO.MechanicDTO value, JsonSerializer serializer)
        {
            var jsonObject = new Newtonsoft.Json.Linq.JObject
        {
            { "mechanic_type", value.mechanicType }
        };

            if (value.@params != null)
            {
                jsonObject.Add("params", Newtonsoft.Json.Linq.JToken.FromObject(value.@params, serializer));
            }

            jsonObject.WriteTo(writer);
        }
    }
}