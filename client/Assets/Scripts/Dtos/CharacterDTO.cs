using System.Collections.Generic;
using System.Runtime.Serialization;
using Newtonsoft.Json;

public class CharacterDTO : DTO
{
    public string id;
    public PositionDTO position;
    public float? radius;
    public int? maxHealth;
    public int? currentHealth;
    public Dictionary<string, int> stats;
    public CharacterAction? action;
    public PositionDTO direction;
    public AbilityDTO[] abilities;
    public InventoryDTO inventory;
}

public enum CharacterAction
{
    [EnumMember(Value = "idle")]
    Idle,

    [EnumMember(Value = "moving")]
    Moving,

    [EnumMember(Value = "attacking")]
    Attacking,

    [EnumMember(Value = "castingHeal")]
    CastingHeal
}

public class AbilityDTO
{
    public string id;
    public string name;
    public float range;
    public float remainingCooldown;
}

public class InventoryDTO
{
    public ItemDTO[] items;
}

public class ItemDTO
{
    public string id;
    public int quantity;
}

public class UseItemDTO : DTO
{
    [JsonProperty("item_id")]
    public string itemId;
    [JsonProperty("target_id")]
    public string targetId;
}
