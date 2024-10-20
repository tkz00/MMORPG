using System;
using System.Collections.Generic;
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
            InventoryItem itemUI = Instantiate(itemUIPrefab, itemsContainer).GetComponent<InventoryItem>();
            itemUI.SetUp(item.Key, item.Value, availableItems.Find(i => i.id == item.Key).icon);
        }
    }
}
