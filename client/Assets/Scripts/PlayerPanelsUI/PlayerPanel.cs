using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class PlayerPanel : MonoBehaviour
{
    Transform cameraTransform;

    void Start()
    {
        Camera camera = Camera.main;
        cameraTransform = camera.transform;
        // From this answer this should select the UI camera, but with the main camera it works (https://discussions.unity.com/t/world-space-canvas-on-top-of-everything/128165/3)
        GetComponent<Canvas>().worldCamera = camera;
    }

    void LateUpdate()
    {
        transform.LookAt(transform.position + cameraTransform.forward);
    }
}
