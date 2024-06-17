using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class Player : MonoBehaviour
{
    [SerializeField]
    private Transform model;

    [SerializeField]
    private PlayerMovement movement;
    public PlayerMovement Movement {
        get { return movement; }
    }

	[SerializeField]
	TMP_Text playerNameUI;

	[SerializeField]
	HealthBar healthBar;

	private PlayerStats stats = new PlayerStats();
	public PlayerStats Stats {
        get { return stats; }
    }

	public void SetPlayerName(string playerName) {
		playerNameUI.text = playerName;
	}

	public void UpdateHealth(int currentHealth, int maxHealth) {
		this.Stats.UpdateHealth(currentHealth, maxHealth);
		this.healthBar.UpdateHealthBar(currentHealth, maxHealth);
	}

	// remove later
	public void SetHealthBarColor(Color color)
	{
		this.healthBar.SetColor(color);
	}

    public void SetScale(float scale)
    {
        model.localScale = Vector3.one * scale;
    }
}
