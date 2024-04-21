using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class PlayerPanel : MonoBehaviour
{
	[SerializeField]
	TMP_Text playerName;

	[SerializeField]
	HealthBar healthBar;

	public void SetPlayerName(string name) {
		this.playerName.text = name;
	}

	public void UpdateHealthBar(int currentHealth, int maxHealth) {
		this.healthBar.UpdateHealthBar(currentHealth, maxHealth);
	}

	// remove later
	public void SetHealthBarColor(Color color)
	{
		this.healthBar.SetColor(color);
	}
}
