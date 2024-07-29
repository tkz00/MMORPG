using System.Runtime.Serialization;

public class PlayerDTO : DTO{
	public string id;
	public PositionDTO position;
	public float radius;
	public int maxHealth;
	public int currentHealth;
	public ExecutingAction executingAction;
}

public enum ExecutingAction
{
   [EnumMember(Value = "idle")]
   Idle,
   [EnumMember(Value = "attacking")]
   Attacking,
   [EnumMember(Value = "castingHeal")]
   CastingHeal
}
