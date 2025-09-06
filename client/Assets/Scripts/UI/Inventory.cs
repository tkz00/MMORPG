using System;
using System.Collections.Generic;
using System.Linq;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.InputSystem;

public class Inventory : MonoBehaviour
{
    [SerializeField] GameObject inventoryMenu;
    [SerializeField] GameObject equipmentMenu;

    [SerializeField] GameObject itemUIPrefab;
    [SerializeField] Transform inventoryContainer;
    [SerializeField] EquipmentWindow equipmentWindow;
    [SerializeField] ItemDatabase itemDatabase;

    readonly Dictionary<string, int> items = new Dictionary<string, int>();
    readonly HashSet<string> equippedItems = new();


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
        equipmentMenu.SetActive(!equipmentMenu.activeSelf);
    }

    public void UpdateInventory(InventoryDTO inventory)
    {
        if (inventory?.items == null)
            return;

        if (inventory.items.Length == 0)
            return;

        foreach ((string id, int quantity, bool isEquipped) in inventory.items.Select(item => (item.id, item.quantity, item.isEquipped)))
        {
            items[id] = quantity;

            if (items[id] == 0)
                items.Remove(id);

            if (isEquipped)
                equippedItems.Add(id);
            else
                equippedItems.Remove(id);
        }

        DrawInventoryIcons();
    }

    private void DrawInventoryIcons()
    {
        // Clear inventory container
        foreach (Transform child in inventoryContainer.transform)
            Destroy(child.gameObject);

        // Clear equipment window slots
        if (equipmentWindow != null)
            equipmentWindow.ClearAllSlots();

        // Draw inventory items
        foreach (var item in items)
        {
            Item itemSO = itemDatabase.Items.SingleOrDefault(i => i.id == item.Key);
            if (itemSO == null)
            {
                Debug.LogError($"Item with ID: {item.Key} not found in item database");
                continue;
            }

            bool isEquippable = itemSO is Equipment;
            bool isEquipped = equippedItems.Contains(item.Key);

            if (isEquipped && isEquippable)
            {
                // Handle equipped items through EquipmentWindow
                equipmentWindow.UpdateEquippedItem(item.Key, item.Value, true, UseItem, EquipItem, UnequipItem);
            }
            else
            {
                // Handle non-equipped items in inventory container
                InventoryItem itemUI = Instantiate(itemUIPrefab).GetComponent<InventoryItem>();
                itemUI.transform.SetParent(inventoryContainer);
                itemUI.SetUp(item.Key, item.Value, itemSO.icon, itemSO.isConsumible, UseItem, isEquippable, isEquippable ? EquipItem : null);
            }
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

    void EquipItem(string itemId)
    {
        if (!GameManager.GetPlayer(GameManager.MainPlayerID).IsAlive) return;
        var equipItem = new EquipItemDTO { itemId = itemId };
        WebSocketMessage response = new WebSocketMessage
        {
            Body = equipItem,
            ActionType = "equip_item"
        };
        string message = JsonConvert.SerializeObject(response);
        WebSocketConnection.SendMessage(message);
    }

    void UnequipItem(string itemId)
    {
        if (!GameManager.GetPlayer(GameManager.MainPlayerID).IsAlive) return;
        var unequipItem = new UnequipItemDTO { itemId = itemId };
        WebSocketMessage response = new WebSocketMessage
        {
            Body = unequipItem,
            ActionType = "unequip_item"
        };
        string message = JsonConvert.SerializeObject(response);
        WebSocketConnection.SendMessage(message);
    }
}
