using System.Collections.Generic;
using UnityEngine;

public class AbilitiesContainer : MonoBehaviour
{
    [SerializeField] public List<Ability> availableAbilities;

    public List<Ability> GetAvailableAbilities()
    {
        return availableAbilities;
    }
}
