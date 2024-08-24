using System;
using System.Collections;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

[RequireComponent(typeof(CharacterController))]
public class Movement : MonoBehaviour
{
    [SerializeField]
    Animator playerAnimator;
    private Coroutine movementCoroutine;
    private float playerSpeed = 10f;
    private CharacterController characterController;

    private static Movement instance;

    public static Movement Instance()
    {
        return instance;
    }

    void Awake()
    {
        characterController = GetComponent<CharacterController>();
        instance = this;
    }

    public void Move(Vector3 target)
    {
        if (movementCoroutine != null)
        {
            StopCoroutine(movementCoroutine);
        }
        movementCoroutine = StartCoroutine(MoveTowards(target));
    }

    private IEnumerator MoveTowards(Vector3 target)
    {
        float playerDistanceToGround = transform.position.y - target.y;
        target.y += playerDistanceToGround;
        while (Vector3.Distance(transform.position, target) > 0.1f)
        {
            Vector3 direction = target - transform.position;
            Vector3 movement = direction.normalized * playerSpeed * Time.deltaTime;
            movement.y = 0;
            characterController.Move(movement);
            transform.rotation = Quaternion.Slerp(
                transform.rotation,
                Quaternion.LookRotation(direction.normalized),
                5f * Time.deltaTime
            );
            yield return null;
        }
    }

    // this shouldn't be handled in the movement script, there should be an animations controller or smth like that
    public void AttackAnimation()
    {
        playerAnimator.SetTrigger("Attack");
    }

    public void HealAnimation()
    {
        playerAnimator.SetTrigger("Heal");
    }

    public void TriggerWalkingAnimation(bool isMoving)
    {
        playerAnimator.SetBool("IsMoving", isMoving);
    }
}
