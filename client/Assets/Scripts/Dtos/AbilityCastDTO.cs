using System.Collections;
using System.Collections.Generic;

public class AbilityCastDTO : DTO
{
    public string name;
    public Dictionary<AbilityParameters, object> abilityParameters;
}

public enum AbilityParameters
{
    TargetPosition,
    TargetId
}

