using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

public class MainPlayer : MonoBehaviour
{
    private Camera mainCamera;

    public InputAction mouseClickAction = new InputAction(binding: "<Mouse>/rightButton");
    public InputAction firstAbilityAction = new InputAction(binding: "<Keyboard>/q");
    public InputAction secondAbilityAction = new InputAction(binding: "<Keyboard>/w");

    private Ability[] abilities = new Ability[4];

    private LayerMask groundLayer;

    public AbilitiesPanel abilitiesPanel;

    void Awake() {
        mainCamera = Camera.main;
    }

    void Start() {
        // The type of shit Unity makes you do:
        groundLayer = (1 << LayerMask.NameToLayer("Ground"));
    }

    private void OnEnable() {
        mouseClickAction.Enable();
        mouseClickAction.performed += SendMovementMessage;

        firstAbilityAction.Enable();
        firstAbilityAction.performed += context => CastAbility(context, abilities[0]);

        secondAbilityAction.Enable();
        secondAbilityAction.performed += context => CastAbility(context, abilities[1]);
    }

    private void OnDisable() {
        mouseClickAction.performed -= SendMovementMessage;
        mouseClickAction.Disable();

        firstAbilityAction.Disable();
        firstAbilityAction.performed -= context => CastAbility(context, abilities[0]);

        secondAbilityAction.Disable();
        secondAbilityAction.performed -= context => CastAbility(context, abilities[1]);
    }

    public void InitAbilities(AbilityDTO[] abilitiesDTOs, List<Ability> availableAbilities)
    {
        int index = 0;
        foreach(AbilityDTO abilityDTO in abilitiesDTOs)
        {
            this.abilities[index] = availableAbilities.Find(ability => ability.id == abilityDTO.id);
            index++;
        }
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

    private void CastAbility(InputAction.CallbackContext context, Ability ability)
    {
        if(ability != null)
        {
            // abilities panel is doing the work of what should be an abilities manager, it should be refactored
            if(abilitiesPanel.CanCast(ability.id))
            {
                Ray ray = mainCamera.ScreenPointToRay(Mouse.current.position.ReadValue());
                if (Physics.Raycast(ray, out RaycastHit hit, Mathf.Infinity, ability.GetLayerMask()) && hit.collider)
                {
                    Dictionary<AbilityParameters, object> abilityParameters = ability.GetAbilityParameters(hit);

                    WebSocketMessage response = new WebSocketMessage
                    {
                        Body = new AbilityCastDTO
                        {
                            abilityParameters = abilityParameters,
                            id = ability.id
                        },
                        ActionType = "AbilityCast"
                    };
                    string message = JsonConvert.SerializeObject(response);
                    WebSocketConnection.SendMessage(message);

                    abilitiesPanel.CastAbility(ability.id);
                }
            }
        }
    }
}

