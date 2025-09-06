using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;
using UnityEngine.UI;

public class InventoryItem : MonoBehaviour
{
    string id;
    [SerializeField] TMP_Text quantity;
    [SerializeField] Image image;
    [SerializeField] Button button;

    public void SetUp(string id, int quantity, Sprite icon, bool isConsumible = false, Action<string> useItem = null, bool isEquippable = false, Action<string> equipItem = null)
    {
        this.id = id;
        this.quantity.text = quantity.ToString();
        image.sprite = icon;

        button.onClick.RemoveAllListeners();

        if (isConsumible)
        {
            button.onClick.AddListener(() => useItem(id));
            button.enabled = true;
        }
        else if (isEquippable && equipItem != null)
        {
            button.onClick.AddListener(() => equipItem(id));
            button.enabled = true;
        }
        else
        {
            button.enabled = false;
        }
    }
}
