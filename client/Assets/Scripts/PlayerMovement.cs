using System;
using System.Collections;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

[RequireComponent(typeof(CharacterController))]
public class PlayerMovement : MonoBehaviour
{
    [SerializeField]
    Animator playerAnimator;
    private Coroutine movementCoroutine;
    private float playerSpeed = 10f;
    private CharacterController characterController;

    void Awake() {
        characterController = GetComponent<CharacterController>();
    }

	public void Move(Vector3 target) {
		if(movementCoroutine != null) {
			StopCoroutine(movementCoroutine);
		}
		movementCoroutine = StartCoroutine(MoveTowards(new Vector3(target.x, this.transform.position.y, target.z)));
	}

    private IEnumerator MoveTowards(Vector3 target) {
        playerAnimator.SetBool("IsMoving", true);
        float playerDistanceToGround = transform.position.y - target.y;
        target.y += playerDistanceToGround;
        while(Vector3.Distance(transform.position, target) > 0.1f) {
            Vector3 direction = target - transform.position;
            Vector3 movement = direction.normalized * playerSpeed * Time.deltaTime;
            characterController.Move(movement);
            transform.rotation = Quaternion.Slerp(transform.rotation, Quaternion.LookRotation(direction.normalized), 1f * Time.deltaTime);
            yield return null;
        }
        playerAnimator.SetBool("IsMoving", false);
    }
}
