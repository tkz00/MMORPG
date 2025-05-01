using System;
using System.Threading;
using System.Threading.Tasks;
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

    private static readonly Color onCooldownOverlayColor = new Color(0.4622642f, 0.4622642f, 0.4622642f, 0.7843137f);
    private static readonly Color failedCastOverlayColor = new Color(0.5490196f, 0.1254902f, 0.1254902f, 0.5f);
    private Color cooldownOverlayColor = onCooldownOverlayColor;
    private CancellationTokenSource colorLerpCancellationTokenSource;



    // This should not be here
    private float abilityCooldown;

    public float AbilityCooldown
    {
        get { return abilityCooldown; }
        private set { abilityCooldown = value; }
    }
    // End of this should not be here

    public void SetAbility(Ability ability)
    {
        this.name = ability.name;
        image.sprite = ability.icon;
    }

    public void UpdateCooldown(float remainingCooldown)
    {
        if (remainingCooldown > 0)
        {
            cooldownOverlay.color = cooldownOverlayColor;
            cooldownOverlay.enabled = true;
            remainingCooldownText.text = remainingCooldown.ToString("F1");
            abilityCooldown = remainingCooldown;
        }
        else
        {
            colorLerpCancellationTokenSource?.Cancel();
            cooldownOverlayColor = onCooldownOverlayColor;
            cooldownOverlay.enabled = false;
            remainingCooldownText.text = String.Empty;
            abilityCooldown = 0f;
        }
    }

    public void FailedCast()
    {
        cooldownOverlay.color = failedCastOverlayColor;
        cooldownOverlay.enabled = true;

        colorLerpCancellationTokenSource?.Cancel();
        colorLerpCancellationTokenSource = new CancellationTokenSource();

        LerpColorAsync(failedCastOverlayColor, onCooldownOverlayColor, .4f, colorLerpCancellationTokenSource.Token);
    }

    private async void LerpColorAsync(Color startColor, Color endColor, float duration, CancellationToken cancellationToken)
    {
        float elapsedTime = 0f;

        while (elapsedTime < duration)
        {
            if (cancellationToken.IsCancellationRequested)
            {
                break;
            }

            cooldownOverlayColor = Color.Lerp(startColor, endColor, elapsedTime / duration);
            elapsedTime += Time.deltaTime;
            await Task.Yield();
        }

        cooldownOverlayColor = endColor;
    }
}
