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

    public void SetUp(string key, int value, Sprite icon)
    {
        id = key;
        quantity.text = value.ToString();
        image.sprite = icon;
    }
}
