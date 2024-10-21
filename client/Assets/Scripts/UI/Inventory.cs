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

    Dictionary<string, int> items = new Dictionary<string, int>();

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

    public void UpdateInventory(IEnumerable<(string id, int quantity)> itemVariations)
    {
        foreach ((string id, int quantity) in itemVariations)
        {
            if (items.ContainsKey(id))
            {
                items[id] += quantity;
            }
            else
            {
                items[id] = quantity;
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
            itemUI.SetUp(item.Key, item.Value, itemSO.icon, UseItem);
        }
    }

    private void UseItem(string itemId)
    {
        UseItemDTO useItem = new UseItemDTO { id = itemId };
        WebSocketMessage response = new WebSocketMessage
        {
            Body = useItem,
            ActionType = "use_item"
        };
        string message = JsonConvert.SerializeObject(response);
        WebSocketConnection.SendMessage(message);
    }
}
