using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public abstract class Ability : ScriptableObject
{
    public string id;
    public new string name;
    public Sprite icon;
    public abstract Dictionary<AbilityParameters, object> GetAbilityParameters(RaycastHit hit);
    public LayerMask targetLayer;
    public LayerMask GetLayerMask() => targetLayer;
}

public enum AbilityParameters // WTF is this name, change
{
    TargetId = 0,
    NoTarget = 1,
    TargetPosition = 2,
}
