using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class Character : MonoBehaviour, ITargeteable
{
    // this field shouldn't be public, it should be set on creation of the player and be readonly from then on
    public string id;

    [SerializeField]
    private Transform model;

	[SerializeField]
    private Transform hitbox;

    [SerializeField]
    private Movement movement;

    public Movement Movement {
        get { return movement; }
    }

	[SerializeField]
	TMP_Text playerNameUI;

	[SerializeField]
	HealthBar healthBar;

    [SerializeField]
    PlayerVFXsHandler playerVFXsHandler;

	private PlayerStats stats = new PlayerStats();

	public PlayerStats Stats {
        get { return stats; }
    }

	public void SetNpcName(string playerName) {
		playerNameUI.text = playerName;
	}

	public void UpdateHealth(int currentHealth, int maxHealth) {
        // this logic should be changed, "states" should come from the backend and them trigger specific feedbacks, i.e.: the "healed" state should trigger the respective healing vfx
        if(this.Stats.CurrentHealth < currentHealth)
        {
            this.playerVFXsHandler.TriggerHealingVFX();
        }

		this.Stats.UpdateHealth(currentHealth, maxHealth);
		this.healthBar.UpdateHealthBar(currentHealth, maxHealth);
	}

    public string GetTargetId()
    {
        return id;
	}

    public void SetMovement(Movement movement) {
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
		hitbox.localScale = Vector3.one * scale;
    }

    public void SetHitbox(bool hitboxOn)
    {
        hitbox.gameObject.SetActive(hitboxOn);
    }
}
