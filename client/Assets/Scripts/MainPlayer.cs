using System.Collections;
using System.Collections.Generic;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

public class MainPlayer : MonoBehaviour
{
	private Camera mainCamera;

    public InputAction mouseClickAction = new InputAction(binding: "<Mouse>/rightButton");
    public InputAction abilityAction = new InputAction(binding: "<Keyboard>/q");

	void Awake() {
        mainCamera = Camera.main;
	}

	private void OnEnable() {
        mouseClickAction.Enable();
        mouseClickAction.performed += SendMovementMessage;

        abilityAction.Enable();
        abilityAction.performed += CastAbilityMessage;
    }

    private void OnDisable() {
        mouseClickAction.performed -= SendMovementMessage;
        mouseClickAction.Disable();

        abilityAction.Disable();
        abilityAction.performed -= CastAbilityMessage;
    }

    private void SendMovementMessage(InputAction.CallbackContext context)
    {
        Ray ray = mainCamera.ScreenPointToRay(Mouse.current.position.ReadValue());
        if(Physics.Raycast(ray: ray, hitInfo: out RaycastHit hit) && hit.collider) {

			float x = hit.point.x, z = hit.point.z;
			PositionDTO inputPosition = new PositionDTO{x = x, z = z};
            WebSocketMessage response = new WebSocketMessage {
                Body = inputPosition,
                ActionType = "Position"
            };
			string message = JsonConvert.SerializeObject(response);
			WebSocketConnection.SendMessage(message);
        }
    }

    private void CastAbilityMessage(InputAction.CallbackContext context) {
        Ray ray = mainCamera.ScreenPointToRay(Mouse.current.position.ReadValue());
        if(Physics.Raycast(ray: ray, hitInfo: out RaycastHit hit) && hit.collider) {

			float x = hit.point.x, z = hit.point.z;
			PositionDTO inputPosition = new PositionDTO{x = x, z = z};
            WebSocketMessage response = new WebSocketMessage {
                Body = new AbilityCastDTO {
                    direction = inputPosition,
                    name = "first-ability"
                },
                ActionType = "AbilityCast"
            };
			string message = JsonConvert.SerializeObject(response);
			WebSocketConnection.SendMessage(message);
        }
    }
}
