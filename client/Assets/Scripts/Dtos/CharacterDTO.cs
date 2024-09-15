using System.Runtime.Serialization;

public class CharacterDTO : DTO
{
    public string id;
    public PositionDTO position;
    public float radius;
    public int maxHealth;
    public int currentHealth;
    public ExecutingAction executingAction;
    public AbilityDTO[] abilities;
}

public class ExecutingAction
{
    public CharacterAction action;
    public PositionDTO direction;
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
