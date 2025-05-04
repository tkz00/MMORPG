using System;
using System.Collections.Generic;
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
        public List<MechanicDTO> mechanics = new();

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

        // Commented things are commented until projectile configuration is mature enough to support it as nested mechanics
        public class CreateProjectileParams : Params
        {
            [JsonProperty("on_hit_mechanics")] public List<MechanicDTO> onHitMechanics;
            // [JsonProperty("targeting_strategy")] public string targetingStrategy;
            // [JsonProperty("radius")] public float radius;
            // [JsonProperty("number")] public int number;
            // [JsonProperty("range")] public float range;

            public CreateProjectileParams()
            {
                onHitMechanics = new List<MechanicDTO>();
                // targetingStrategy = "arc";
                // radius = 0;
                // number = 1;
                // range = 8;
            }

            public override Params DeepCopy()
            {
                List<MechanicDTO> onHitMechanicsCopy = new List<MechanicDTO>();
                foreach (MechanicDTO mechanicDTO in this.onHitMechanics)
                {
                    onHitMechanicsCopy.Add(mechanicDTO.DeepCopy());
                }
                return new CreateProjectileParams
                {
                    onHitMechanics = onHitMechanicsCopy,
                    // targetingStrategy = this.targetingStrategy,
                };
            }
        }

        public class CreateAoEParams : Params
        {
            [JsonProperty("duration_ms")] public int durationMs;
            [JsonProperty("radius")] public float radius;
            [JsonProperty("on_hit_mechanics")] public List<MechanicDTO> onHitMechanics;

            public CreateAoEParams()
            {
                durationMs = 400;
                radius = 1;
                onHitMechanics = new List<MechanicDTO>();
            }

            public override Params DeepCopy()
            {
                List<MechanicDTO> onHitMechanicsCopy = new List<MechanicDTO>();
                foreach (MechanicDTO mechanicDTO in this.onHitMechanics)
                {
                    onHitMechanicsCopy.Add(mechanicDTO.DeepCopy());
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
            [JsonProperty("base_amount")] public int baseAmount;
            [JsonProperty("damage_stat_multiplier")] public float damageStatMultiplier;
            [JsonProperty("targeting_strategy")] public string targetingStrategy;
            [JsonProperty("on_hit_mechanics")] public List<MechanicDTO> onHitMechanics;

            public DamageParams()
            {
                baseAmount = 10;
                damageStatMultiplier = 0;
                targetingStrategy = "TO DO";
                onHitMechanics = new List<MechanicDTO>();
            }

            public override Params DeepCopy()
            {
                List<MechanicDTO> onHitMechanicsCopy = new List<MechanicDTO>();
                foreach (MechanicDTO mechanicDTO in this.onHitMechanics)
                {
                    onHitMechanicsCopy.Add(mechanicDTO.DeepCopy());
                }
                return new DamageParams
                {
                    baseAmount = this.baseAmount,
                    damageStatMultiplier = this.damageStatMultiplier,
                    targetingStrategy = this.targetingStrategy,
                    onHitMechanics = onHitMechanicsCopy
                };
            }
        }

        public class HealParams : Params
        {
            [JsonProperty("base_amount")] public int baseAmount;
            [JsonProperty("damage_stat_multiplier")] public float damageStatMultiplier;
            [JsonProperty("targeting_strategy")] public string targetingStrategy;
            [JsonProperty("on_hit_mechanics")] public List<MechanicDTO> onHitMechanics;

            public HealParams()
            {
                baseAmount = 10;
                damageStatMultiplier = 0;
                targetingStrategy = "TO DO";
                onHitMechanics = new List<MechanicDTO>();
            }

            public override Params DeepCopy()
            {
                List<MechanicDTO> onHitMechanicsCopy = new List<MechanicDTO>();
                foreach (MechanicDTO mechanicDTO in this.onHitMechanics)
                {
                    onHitMechanicsCopy.Add(mechanicDTO.DeepCopy());
                }
                return new HealParams
                {
                    baseAmount = this.baseAmount,
                    damageStatMultiplier = this.damageStatMultiplier,
                    targetingStrategy = this.targetingStrategy,
                    onHitMechanics = onHitMechanicsCopy
                };
            }
        }

        public class DelayParams : Params
        {
            [JsonProperty("delay_ms")] public int delayMs;
            [JsonProperty("execute_after_delay_mechanics")] public List<MechanicDTO> executeAfterDelayMechanics;

            public DelayParams()
            {
                delayMs = 100;
                executeAfterDelayMechanics = new List<MechanicDTO>();
            }
            public override Params DeepCopy()
            {
                List<MechanicDTO> executeAfterDelayMechanicsCopy = new List<MechanicDTO>();
                foreach (MechanicDTO mechanicDTO in this.executeAfterDelayMechanics)
                {
                    executeAfterDelayMechanicsCopy.Add(mechanicDTO.DeepCopy());
                }
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
                "heal" => jsonObject["params"]?.ToObject<AbilityDTO.HealParams>(serializer),
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