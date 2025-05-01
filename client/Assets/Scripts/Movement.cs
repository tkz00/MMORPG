using System;
using System.Collections;
using UnityEngine;

[RequireComponent(typeof(CharacterController))]
public class Movement : MonoBehaviour
{
    [SerializeField]
    Animator playerAnimator;
    private Coroutine movementCoroutine;
    private float playerSpeed = 10f;
    private CharacterController characterController;

    void Awake()
    {
        characterController = GetComponent<CharacterController>();
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
            yield return null;
        }
    }

    public void RotateTowards(PositionDTO direction)
    {
        Vector3 directionRotation = new Vector3(direction.x, 0, direction.z);
        if (directionRotation != Vector3.zero)
        {
            Quaternion targetRotation = Quaternion.LookRotation(directionRotation);
            transform.rotation = Quaternion.Euler(0, targetRotation.eulerAngles.y, 0);
        }
    }

    public void Respawn(Vector3 position)
    {
        transform.position = position;
        playerAnimator.SetTrigger("Spawn");
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

    public void DeathAnimation()
    {
        playerAnimator.SetTrigger("Death");
    }
}
