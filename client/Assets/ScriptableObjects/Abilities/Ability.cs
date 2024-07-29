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

public enum AbilityParameters
{
    TargetPosition,
    TargetId
}

[CreateAssetMenu(menuName = "Abilities/DirectionalAbility")]
public class DirectionalAbility : Ability
{
    public override Dictionary<AbilityParameters, object> GetAbilityParameters(RaycastHit hit)
    {
        float x = hit.point.x, z = hit.point.z;
        PositionDTO inputPosition = new PositionDTO { x = x, z = z };
        return new Dictionary<AbilityParameters, object>
        {
            { AbilityParameters.TargetPosition, inputPosition }
        };
    }
}

[CreateAssetMenu(menuName = "Abilities/TargetedAbility")]
public class TargetedAbility : Ability
{
    public override Dictionary<AbilityParameters, object> GetAbilityParameters(RaycastHit hit)
    {
        string targetId = hit.collider.GetComponent<ITargeteable>().GetTargetId();
        return new Dictionary<AbilityParameters, object>
        {
            { AbilityParameters.TargetId, targetId }
        };
    }
}

