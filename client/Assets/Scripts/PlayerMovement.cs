using System;
using System.Collections;
using Newtonsoft.Json;
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
    private string playerID;

    async void Awake() {
        mainCamera = Camera.main;
        characterController = GetComponent<CharacterController>();

		WebSocketConnection.SetHandler<string>(new Action<string>((playerId) => {
			this.playerID = playerId;
			Debug.Log(this.playerID);
		}));

		WebSocketConnection.SetHandler<GameStateDTO>(new Action<GameStateDTO>((gameState) => {
			PlayerDTO player = gameState.Players.Find(player => player.Id == this.playerID);

			if(movementCoroutine != null) {
                StopCoroutine(movementCoroutine);
            }
            movementCoroutine = StartCoroutine(MoveTowards(new Vector3(player.Position.x, this.transform.position.y, player.Position.z)));
		}));

        await WebSocketConnection.Connect();
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

			float x = hit.point.x, z = hit.point.z;
			PositionDTO inputPosition = new PositionDTO{x = x, z = z};
			string message = JsonConvert.SerializeObject(inputPosition);
			WebSocketConnection.SendMessage(message);
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
