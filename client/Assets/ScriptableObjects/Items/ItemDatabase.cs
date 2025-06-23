using System.Collections.Generic;
using UnityEngine;

[CreateAssetMenu(fileName = "ItemDatabase", menuName = "Item/ItemDatabase")]
public class ItemDatabase : ScriptableObject
{
    public List<Item> Items;
}