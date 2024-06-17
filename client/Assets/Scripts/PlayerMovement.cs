using System;
using System.Collections;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

[RequireComponent(typeof(CharacterController))]
public class PlayerMovement : MonoBehaviour
{
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
		movementCoroutine = StartCoroutine(MoveTowards(target));
	}

    private IEnumerator MoveTowards(Vector3 target) {
        while(Vector3.Distance(transform.position, target) > 0.1f) {
            Vector3 direction = target - transform.position;
            Vector3 movement = direction.normalized * playerSpeed * Time.deltaTime;
			movement.y = 0;
            characterController.Move(movement);
            // transform.rotation = Quaternion.Slerp(transform.rotation, Quaternion.LookRotation(direction.normalized), 1f * Time.deltaTime);
            yield return null;
        }
    }
}
