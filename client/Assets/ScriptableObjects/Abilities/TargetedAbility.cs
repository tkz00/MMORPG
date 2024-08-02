using System.Collections.Generic;
using UnityEngine;

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
