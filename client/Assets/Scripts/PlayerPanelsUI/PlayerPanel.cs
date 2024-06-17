using System;
using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class PlayerPanel : MonoBehaviour
{
	Transform cameraTransform;

	void Start() {
		cameraTransform = Camera.main.transform;
	}

	void LateUpdate() {
		transform.LookAt(transform.position + cameraTransform.forward);
	}
}
