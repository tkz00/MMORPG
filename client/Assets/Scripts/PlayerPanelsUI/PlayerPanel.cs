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

	public void Initialize(string name, int currentHealth, int maxHealth) {
		this.playerName.text = name;
		this.healthBar.UpdateHealthBar(currentHealth, maxHealth);
	}
}
