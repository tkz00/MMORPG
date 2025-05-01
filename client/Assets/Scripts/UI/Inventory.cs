using System;
using System.Collections.Generic;
using System.Linq;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

public class Inventory : MonoBehaviour
{
    [SerializeField] GameObject inventoryMenu;

    [SerializeField] GameObject itemUIPrefab;
    [SerializeField] Transform itemsContainer;
    [SerializeField] List<Item> availableItems;

    readonly Dictionary<string, int> items = new Dictionary<string, int>();


    public InputAction openCloseInventoryAction = new InputAction(binding: "<Keyboard>/i");

    private void OnEnable()
    {
        openCloseInventoryAction.Enable();
        openCloseInventoryAction.performed += ToggleInventory;
    }

    private void OnDisable()
    {
        openCloseInventoryAction.performed -= ToggleInventory;
        openCloseInventoryAction.Disable();
    }

    private void ToggleInventory(InputAction.CallbackContext context)
    {
        inventoryMenu.SetActive(!inventoryMenu.activeSelf);
    }

    public void UpdateInventory(InventoryDTO inventory)
    {
        if (inventory?.items == null)
            return;

        if (inventory.items.Length == 0)
            return;

        foreach ((string id, int quantity) in inventory.items.Select(item => (item.id, item.quantity)))
        {
            items[id] = quantity;

            if (items[id] == 0)
            {
                items.Remove(id);
            }
        }

        DrawInventoryIcons();
    }

    private void DrawInventoryIcons()
    {
        foreach (Transform child in itemsContainer.transform)
        {
            Destroy(child.gameObject);
        }

        foreach (var item in items)
        {
            Item itemSO = availableItems.SingleOrDefault(i => i.id == item.Key);
            if (itemSO == null)
            {
                Debug.LogError($"Item with ID: {item.Key} not found in available items collection");
                continue;
            }
            InventoryItem itemUI = Instantiate(itemUIPrefab, itemsContainer).GetComponent<InventoryItem>();
            itemUI.SetUp(item.Key, item.Value, itemSO.icon, itemSO.isConsumible, UseItem);
        }
    }

    private void UseItem(string itemId)
    {
        if (!GameManager.GetPlayer(GameManager.MainPlayerID).IsAlive) return;
        UseItemDTO useItem = new UseItemDTO { itemId = itemId, targetId = GameManager.MainPlayerID };
        WebSocketMessage response = new WebSocketMessage
        {
            Body = useItem,
            ActionType = "use_item"
        };
        string message = JsonConvert.SerializeObject(response);
        WebSocketConnection.SendMessage(message);
    }
}
