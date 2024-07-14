using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class PlayerVFXsHandler : MonoBehaviour
{
    [SerializeField]
    GameObject healVFX;

    public void TriggerHealingVFX()
    {
        if (healVFX != null)
        {
            healVFX.SetActive(true);
            StartCoroutine(DeactivateAfterDelay(0.5f));
        }
    }

    private IEnumerator DeactivateAfterDelay(float delay)
    {
        yield return new WaitForSeconds(delay);
        healVFX.SetActive(false);
    }
}
