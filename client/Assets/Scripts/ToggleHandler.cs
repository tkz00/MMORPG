using System.Collections;
using System.Collections.Generic;
using UnityEngine;
using UnityEngine.UI;

// eat my ass unity
public class ToggleHandler : MonoBehaviour
{
    bool toggleState;

    [SerializeField]
    Toggle toggle;

    void Start()
    {
        toggleState = toggle.isOn;
    }

    public void Toggle()
    {
        toggleState = !toggleState;
        toggle.isOn = toggleState;
    }
}
