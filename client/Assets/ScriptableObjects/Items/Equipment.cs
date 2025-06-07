using UnityEngine;

[CreateAssetMenu(menuName = "Item/Equipment")]
public class Equipment : Item
{
    public enum EquipmentType
    {
        Helmet,
        Chest,
        Boots
    }

    public EquipmentType equipmentType;
}
