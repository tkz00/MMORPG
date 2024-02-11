using System;
using System.Collections;
using UnityEngine;
using UnityEngine.InputSystem;

[RequireComponent(typeof(CharacterController))]
public class PlayerMovement : MonoBehaviour
{
    [SerializeField]
    private InputAction mouseClickAction;

    private Camera mainCamera;

    private Coroutine movementCoroutine;
    private float playerSpeed = 10f;
    private CharacterController characterController;

    void Awake() {
        mainCamera = Camera.main;
        characterController = GetComponent<CharacterController>();
    }

    private void OnEnable() {
        mouseClickAction.Enable();
        mouseClickAction.performed += Move;
    }

    private void OnDisable() {
        mouseClickAction.performed -= Move;
        mouseClickAction.Disable();
    }

    private void Move(InputAction.CallbackContext context)
    {
        Ray ray = mainCamera.ScreenPointToRay(Mouse.current.position.ReadValue());
        if(Physics.Raycast(ray: ray, hitInfo: out RaycastHit hit) && hit.collider) {
            if(movementCoroutine != null) {
                StopCoroutine(movementCoroutine);
            }
            movementCoroutine = StartCoroutine(MoveTowards(hit.point));
        }
    }

    private IEnumerator MoveTowards(Vector3 target) {
        float playerDistanceToGround = transform.position.y - target.y;
        target.y += playerDistanceToGround;
        while(Vector3.Distance(transform.position, target) > 0.1f) {
            Vector3 direction = target - transform.position;
            Vector3 movement = direction.normalized * playerSpeed * Time.deltaTime;
            characterController.Move(movement);
            transform.rotation = Quaternion.Slerp(transform.rotation, Quaternion.LookRotation(direction.normalized), 1f * Time.deltaTime);
            yield return null;
        }
    }
}
