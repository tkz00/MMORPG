using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class Player : MonoBehaviour, ITargeteable
{
    // this field shouldn't be public, it should be set on creation of the player and be readonly from then on
    public string id;

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

    public string GetTargetId()
    {
        return id;
	}
    public void SetMovement(PlayerMovement movement) {
        this.movement = movement;
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
