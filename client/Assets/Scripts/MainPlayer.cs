using System.Collections;
using System.Collections.Generic;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

public class MainPlayer : MonoBehaviour
{
    private Camera mainCamera;

    public InputAction mouseClickAction = new InputAction(binding: "<Mouse>/rightButton");
    public InputAction projectileAbilityAction = new InputAction(binding: "<Keyboard>/q");
    public InputAction healAbilityAction = new InputAction(binding: "<Keyboard>/w");

    private LayerMask groundLayer;
    private LayerMask playersLayer;

    void Awake() {
        mainCamera = Camera.main;
    }

    void Start() {
        // The type of shit Unity makes you do:
        groundLayer = (1 << LayerMask.NameToLayer("Ground"));
        playersLayer = (1 << LayerMask.NameToLayer("Players"));
    }

    private void OnEnable() {
        mouseClickAction.Enable();
        mouseClickAction.performed += SendMovementMessage;

        projectileAbilityAction.Enable();
        projectileAbilityAction.performed += CastProjectileAbility;

        healAbilityAction.Enable();
        healAbilityAction.performed += CastHealAbility;
    }

    private void OnDisable() {
        mouseClickAction.performed -= SendMovementMessage;
        mouseClickAction.Disable();

        projectileAbilityAction.Disable();
        projectileAbilityAction.performed -= CastProjectileAbility;

        healAbilityAction.Disable();
        healAbilityAction.performed -= CastHealAbility;
    }

    private void SendMovementMessage(InputAction.CallbackContext context)
    {
        Ray ray = mainCamera.ScreenPointToRay(Mouse.current.position.ReadValue());
        if(Physics.Raycast(ray, out RaycastHit hit, Mathf.Infinity, groundLayer) && hit.collider) {
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

    private void CastProjectileAbility(InputAction.CallbackContext context)
    {
        Ray ray = mainCamera.ScreenPointToRay(Mouse.current.position.ReadValue());
        if(Physics.Raycast(ray, out RaycastHit hit, Mathf.Infinity, groundLayer) && hit.collider) {

            float x = hit.point.x, z = hit.point.z;
            PositionDTO inputPosition = new PositionDTO{x = x, z = z};
            Dictionary<AbilityParameters, object> abilityParameters =
                new Dictionary<AbilityParameters, object>{
                    {AbilityParameters.TargetPosition, inputPosition}
                };

            WebSocketMessage response = new WebSocketMessage {
                Body = new AbilityCastDTO {
                    abilityParameters = abilityParameters,
                    name = "projectile"
                },
                ActionType = "AbilityCast"
            };
            string message = JsonConvert.SerializeObject(response);
            WebSocketConnection.SendMessage(message);
            PlayerMovement.Instance().AttackAnimation();
        }
    }

    private void CastHealAbility(InputAction.CallbackContext context)
    {
        Ray ray = mainCamera.ScreenPointToRay(Mouse.current.position.ReadValue());
        if(Physics.Raycast(ray, out RaycastHit hit, Mathf.Infinity, playersLayer) && hit.collider) {
            
            string targetId = hit.collider.GetComponent<ITargeteable>().GetTargetId();
            Dictionary<AbilityParameters, object> abilityParameters =
                new Dictionary<AbilityParameters, object>{
                    {AbilityParameters.TargetId, targetId}
                };

            WebSocketMessage response = new WebSocketMessage {
                Body = new AbilityCastDTO {
                    abilityParameters = abilityParameters,
                    name = "heal"
                },
                ActionType = "AbilityCast"
            };
            string message = JsonConvert.SerializeObject(response);
            WebSocketConnection.SendMessage(message);
        }
    }
}

