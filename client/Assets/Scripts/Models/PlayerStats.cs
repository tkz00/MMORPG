using System.Collections;
using System.Collections.Generic;

public class PlayerStats
{
    public int MaxHealth { private set; get; }
    public int CurrentHealth { private set; get; }

    public void UpdateHealth(int currentHealth, int maxHealth)
    {
        this.CurrentHealth = currentHealth;
        this.MaxHealth = maxHealth;
    }
}
