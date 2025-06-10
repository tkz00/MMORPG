using UnityEngine;

[CreateAssetMenu(menuName = "Item/Equipment")]
public class Equipment : Item
{

    public EquipmentType equipmentType;
}

public enum EquipmentType
{
    Helmet,
    Chest,
    Boots
}
