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

    public void SetUp(string key, int value, Sprite icon, Action<string> useItem)
    {
        id = key;
        quantity.text = value.ToString();
        image.sprite = icon;
        button.onClick.AddListener(() => useItem(id));
    }
}
