using System;
using Newtonsoft.Json;
using UnityEngine;
public class WebSocketResponse <T> where T : DTO {
	public T Body;
	public string Type;
}
