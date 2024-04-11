using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class Player : MonoBehaviour
{
    [SerializeField]
    private PlayerMovement movement;
    public PlayerStats Stats { get; set; }

    public PlayerMovement Movement {
        get { return movement; }
    }
}
