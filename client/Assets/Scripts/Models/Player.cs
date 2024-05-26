using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Player : MonoBehaviour
{
    [SerializeField]
    private PlayerMovement movement;
    public PlayerMovement Movement {
        get { return movement; }
    }

	private PlayerStats stats = new PlayerStats();
	public PlayerStats Stats {
        get { return stats; }
    }

	public Action<int, int> onHealthChanged; 

	public void UpdateHealth(int currentHealth, int maxHealth) {
		this.Stats.UpdateHealth(currentHealth, maxHealth);
		onHealthChanged?.Invoke(currentHealth, maxHealth);
	}

    public void SetMovement(PlayerMovement movement) {
        this.movement = movement;
    }
}
