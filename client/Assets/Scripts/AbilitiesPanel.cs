using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class AbilitiesPanel : MonoBehaviour
{
    [SerializeField]
    GameObject abilityIconPrefab;

    [SerializeField]
    Transform abilityIconsContainer;

    Dictionary<string, GameObject> abilityIcons = new Dictionary<string, GameObject>();

    [SerializeField]
    List<Ability> abilities;

    const float CAST_ABILITY_ICON_ANIMATION_DURATION = 0.1f;

    void Start()
    {
        foreach(Ability ability in abilities)
        {
            GameObject abilityIcon = Instantiate(abilityIconPrefab, abilityIconsContainer);
            abilityIcon.name = ability.name;
            abilityIcon.GetComponent<Image>().sprite = ability.icon;
            abilityIcons.Add(ability.name, abilityIcon);
        }

        LayoutRebuilder.ForceRebuildLayoutImmediate(abilityIconsContainer.GetComponent<RectTransform>());
    }

    public void CastAbility(string abilityName)
    {
        StartCoroutine(AnimateAbility(abilityIcons[abilityName].GetComponent<Image>()));
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
}
