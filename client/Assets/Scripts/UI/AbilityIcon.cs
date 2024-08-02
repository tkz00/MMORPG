using System;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

public class AbilityIcon : MonoBehaviour
{
    [SerializeField]
    Image image;

    [SerializeField]
    Image cooldownOverlay;

    [SerializeField]
    TextMeshProUGUI remainingCooldownText;

    public void SetAbility(Ability ability)
    {
        this.name = ability.name;
        image.sprite = ability.icon;
    }

    public void UpdateCooldown(float remainingCooldown)
    {
        Debug.Log(remainingCooldown);

        if(remainingCooldown > 0)
        {
            cooldownOverlay.enabled = true;
            remainingCooldownText.text = remainingCooldown.ToString("F1");
        }
        else
        {
            remainingCooldownText.text = String.Empty;
            cooldownOverlay.enabled = false;
        }
    }
}
