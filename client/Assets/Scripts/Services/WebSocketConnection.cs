using System;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Net.WebSockets;
using System.Collections.Generic;
using UnityEngine;
using Newtonsoft.Json;
using Newtonsoft.Json.Linq;

public static class WebSocketConnection 
{
    static ClientWebSocket webSocket;
    static String URL = "ws://localhost:3009/ws";
	static Dictionary<Type, Delegate> _responseHandlers = new Dictionary<Type, Delegate>();
	
    public static async Task Connect() {
        // handshake
        webSocket = new ClientWebSocket();
        Uri wsUri = new Uri(URL);
        try {
            webSocket.Options.SetRequestHeader("Origin", "http://example.com");
            await webSocket.ConnectAsync(wsUri, CancellationToken.None);
			ReadLoopAsync();
            // Debug.Log("conexion satisfactoria");
        } catch(Exception ex) {
            Debug.Log(ex.Message);
        }
    }

	public static void SetHandler<TResponse>(Action<TResponse> responseHandler) where TResponse : DTO
    {
        _responseHandlers[typeof(TResponse)] = responseHandler;
    }

	public static async void SendMessage(string messageJson)
    {
		var encodedMessage = Encoding.UTF8.GetBytes(messageJson);
        var wsBuffer = new ArraySegment<Byte>(encodedMessage, 0, encodedMessage.Length);

        try {
            await webSocket.SendAsync(wsBuffer, WebSocketMessageType.Text, true, CancellationToken.None);
        } catch(Exception ex) {
            Debug.Log(ex.Message);
        }
    }

	private static async void ReadLoopAsync() {
		List<byte> messageBytes = new List<byte>();
        byte[] receiveBuffer = new byte[1024];

        while (webSocket.State == WebSocketState.Open)
        {
            WebSocketReceiveResult result = null;
            do
            {
                result = await webSocket.ReceiveAsync(new ArraySegment<byte>(receiveBuffer), CancellationToken.None);
                for (int i = 0; i < result.Count; i++)
                {
                    messageBytes.Add(receiveBuffer[i]);
                }
            }
            while (!result.EndOfMessage);

            if (result.MessageType == WebSocketMessageType.Binary)
            {
				string responseJson = Encoding.UTF8.GetString(messageBytes.ToArray());
                JObject jsonObject = JObject.Parse(responseJson);
                string responseTypeString = (string)jsonObject["type"];
                switch (responseTypeString) {
                    case "PositionDTO":
                        runHandler<PositionDTO>(responseJson);
                        break;
                    case "PlayerDTO":
                        runHandler<PlayerDTO>(responseJson);
                        break;
                    case "GameStateDTO":
                        runHandler<GameStateDTO>(responseJson);
                        break;
                }
				messageBytes.Clear();
            }
            else if (result.MessageType == WebSocketMessageType.Close)
            {
                await webSocket.CloseAsync(WebSocketCloseStatus.NormalClosure, string.Empty, CancellationToken.None);
            }
            else
            {
                Debug.LogError("Message from server is not text");
            }
        }
	}

    static void runHandler<T>(string responseJson) where T : DTO {
        WebSocketResponse<T> response = JsonConvert.DeserializeObject<WebSocketResponse<T>>(responseJson);
        if(_responseHandlers.ContainsKey(typeof(T))) {
            var handler = _responseHandlers[typeof(T)];
            handler.DynamicInvoke(response.Body);
        }
    }
}
