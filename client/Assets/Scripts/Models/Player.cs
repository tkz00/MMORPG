using System;
using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Player : MonoBehaviour
{
    // this field shouldn't be public, it should be set on creation of the player and be readonly from then on
    public string id;

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
}
