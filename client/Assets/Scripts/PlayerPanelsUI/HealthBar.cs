using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

public class HealthBar : MonoBehaviour
{
	[SerializeField]
	Slider healthSlider;

	// this has to be called by an event
	public void UpdateHealthBar(int currentHealth, int maxHealth)
	{
		float fillAmount = currentHealth / (float)maxHealth;
        healthSlider.value = fillAmount;
	}
}
