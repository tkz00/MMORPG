using System.Collections.Generic;
using UnityEngine;

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
