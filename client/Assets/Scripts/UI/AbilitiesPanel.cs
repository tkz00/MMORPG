using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

public class AbilitiesPanel : MonoBehaviour
{
    [SerializeField]
    GameObject abilityIconPrefab;

    [SerializeField]
    Transform abilityIconsContainer;

    [SerializeField]
    AbilitiesContainer abilitiesContainer;

    Dictionary<string, AbilityIcon> abilityIcons = new Dictionary<string, AbilityIcon>();

    const float CAST_ABILITY_ICON_ANIMATION_DURATION = 0.1f;

    public void Init(AbilityDTO[] abilitiesDTOs)
    {
        foreach (AbilityDTO abilityDTO in abilitiesDTOs)
        {
            Ability ability = abilitiesContainer.GetAvailableAbilities().Find(ability => ability.id == abilityDTO.id);
            AbilityIcon abilityIcon = Instantiate(abilityIconPrefab, abilityIconsContainer).GetComponent<AbilityIcon>();
            abilityIcon.SetAbility(ability);
            abilityIcons.Add(ability.id, abilityIcon);
        }

        LayoutRebuilder.ForceRebuildLayoutImmediate(abilityIconsContainer.GetComponent<RectTransform>());
    }

    public bool CanCast(string abilityId)
    {
        AbilityIcon abilityIcon = abilityIcons[abilityId];
        if (abilityIcons[abilityId].AbilityCooldown == 0)
        {
            return true;
        }
        else
        {
            abilityIcon.FailedCast();
            return false;
        }
    }

    public void CastAbility(string abilityId)
    {
        // Will this ever be needed? Will leave like this for now
        // StartCoroutine(AnimateAbility(abilityIcons[abilityName].GetComponent<Image>()));
    }

    private IEnumerator AnimateAbility(Image abilityImage)
    {
        Color originalColor = abilityImage.color;
        Color targetColor = originalColor;
        targetColor.a = 0.5f;

        float elapsedTime = 0f;

        while (elapsedTime < CAST_ABILITY_ICON_ANIMATION_DURATION)
        {
            abilityImage.color = Color.Lerp(originalColor, targetColor, elapsedTime / CAST_ABILITY_ICON_ANIMATION_DURATION);
            elapsedTime += Time.deltaTime;
            yield return null;
        }

        abilityImage.color = targetColor;

        yield return new WaitForSeconds(CAST_ABILITY_ICON_ANIMATION_DURATION);

        elapsedTime = 0f;

        while (elapsedTime < CAST_ABILITY_ICON_ANIMATION_DURATION)
        {
            abilityImage.color = Color.Lerp(targetColor, originalColor, elapsedTime / CAST_ABILITY_ICON_ANIMATION_DURATION);
            elapsedTime += Time.deltaTime;
            yield return null;
        }

        originalColor.a = 1;
        abilityImage.color = originalColor;
    }

    public void UpdatePlayerPanel(CharacterDTO playerDTO)
    {
        if (playerDTO.abilities == null)
            return;

        foreach (AbilityDTO abilityDTO in playerDTO.abilities)
        {
            if (abilityIcons[abilityDTO.id] != null)
            {
                abilityIcons[abilityDTO.id].UpdateCooldown(abilityDTO.remainingCooldown);
            }
        }
    }
}
