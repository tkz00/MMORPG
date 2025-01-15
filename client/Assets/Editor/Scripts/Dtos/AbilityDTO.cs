public class AbilityDTO
{
    public string id;
    public string name;
    public int cooldown;
    public float range;
    public AbilityParameters targeting;
    public CharacterAction characterAction;
    public MechanicDTO[] mechanics;

    public class MechanicDTO
    {
        public string mechanicType;
    }
}
