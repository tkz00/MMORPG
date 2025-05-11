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

    public Movement Movement
    {
        get { return movement; }
    }

    [SerializeField]
    TMP_Text playerNameUI;

    [SerializeField]
    HealthBar healthBar;

    [SerializeField]
    PlayerVFXsHandler playerVFXsHandler;

    private PlayerStats stats = new PlayerStats();

    public PlayerStats Stats
    {
        get { return stats; }
    }

    public bool IsAlive { get; private set; } = true;

    public void SetCharacterName(string playerName)
    {
        playerNameUI.text = playerName;
    }

    public void UpdateHealth(int? currentHealth, int? maxHealth)
    {
        // this logic should be changed, "states" should come from the backend and them trigger specific feedbacks, i.e.: the "healed" state should trigger the respective healing vfx
        if (currentHealth != null && this.Stats.CurrentHealth < currentHealth)
        {
            this.playerVFXsHandler.TriggerHealingVFX();
        }

        this.stats.UpdateHealth(currentHealth, maxHealth);
        this.healthBar.UpdateHealthBar(this.stats.CurrentHealth, this.stats.MaxHealth);
    }

    public string GetTargetId()
    {
        return id;
    }

    public void SetMovement(Movement movement)
    {
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

    public void Respawn(Vector3 position)
    {
        IsAlive = true;
        this.Movement.Respawn(position);
        healthBar.gameObject.SetActive(true);
        playerNameUI.transform.parent.gameObject.SetActive(true);
    }

    public void TriggerDeath()
    {
        IsAlive = false;
        Movement.DeathAnimation();
        healthBar.gameObject.SetActive(false);
        playerNameUI.transform.parent.gameObject.SetActive(false);
    }

    public void HandleActionFeedback(CharacterAction action, int? animationDuration)
    {
        switch (action)
        {
            case CharacterAction.Attacking:
                Movement.TriggerWalkingAnimation(false);
                Movement.AttackAnimation(animationDuration);
                break;
            case CharacterAction.CastingHeal:
                Movement.TriggerWalkingAnimation(false);
                Movement.HealAnimation(animationDuration);
                break;
            case CharacterAction.Moving:
                Movement.TriggerWalkingAnimation(true);
                break;
            default:
                Movement.TriggerWalkingAnimation(false);
                break;
        }
    }
}
