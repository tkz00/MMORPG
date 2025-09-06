using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;

public class EquipmentWindow : MonoBehaviour
{
    [SerializeField] InventoryItem helmetSlot;
    [SerializeField] InventoryItem chestSlot;
    [SerializeField] InventoryItem bootsSlot;
    [SerializeField] ItemDatabase itemDatabase;

    private Dictionary<EquipmentType, InventoryItem> equipmentSlots;
    private Dictionary<EquipmentType, InventoryItem> equippedItems = new Dictionary<EquipmentType, InventoryItem>();

    private void Awake()
    {
        // Initialize equipment slots dictionary
        equipmentSlots = new Dictionary<EquipmentType, InventoryItem>
        {
            { EquipmentType.Helmet, helmetSlot },
            { EquipmentType.Chest, chestSlot },
            { EquipmentType.Boots, bootsSlot }
        };
    }

    public void UpdateEquippedItem(string itemId, int quantity, bool isEquipped, System.Action<string> onUseItem, System.Action<string> onEquipItem, System.Action<string> onUnequipItem)
    {
        // Find the equipment data
        Equipment equipmentData = itemDatabase.Items
            .OfType<Equipment>()
            .SingleOrDefault(e => e.id == itemId);

        if (equipmentData == null)
            return;

        EquipmentType equipmentType = equipmentData.equipmentType;

        // Remove existing item from this slot if any
        if (equippedItems.ContainsKey(equipmentType))
        {
            if (equippedItems[equipmentType] != null)
            {
                EmptySlot(equipmentType);
            }
            equippedItems.Remove(equipmentType);
        }

        // Handle the item based on equipped status
        if (isEquipped && quantity > 0)
        {
            // Equip the item
            InventoryItem itemUI = equipmentSlots[equipmentType];
            itemUI.transform.GetChild(2).gameObject.SetActive(false);
            itemUI.transform.GetChild(0).gameObject.SetActive(true);
            itemUI.SetUp(itemId, quantity, equipmentData.icon, equipmentData.isConsumible, onUseItem, true, onUnequipItem);
            equippedItems[equipmentType] = itemUI;
        }
        else
        {
            // Unequip the item - ensure slot is empty
            EmptySlot(equipmentType);
        }
    }

    public void ClearAllSlots()
    {
        foreach (var equippedItem in equippedItems)
        {
            if (equippedItem.Value != null)
            {
                EmptySlot(equippedItem.Key);
            }
        }
        equippedItems.Clear();
    }

    public bool IsSlotEmpty(EquipmentType equipmentType)
    {
        return !equippedItems.ContainsKey(equipmentType) || equippedItems[equipmentType] == null;
    }

    public void EmptySlot(EquipmentType equipmentType)
    {
        equipmentSlots[equipmentType].SetUp("", 0, null, false, null, false, null);
        equipmentSlots[equipmentType].transform.GetChild(0).gameObject.SetActive(false);
        equipmentSlots[equipmentType].transform.GetChild(2).gameObject.SetActive(true);
    }
}
