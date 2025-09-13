using System;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Net.WebSockets;
using System.Collections.Generic;
using UnityEngine;
using Newtonsoft.Json;
using UnityEditor;

public static class WebSocketConnection
{
    static ClientWebSocket webSocket;
    static String URL = "ws://localhost:3009/ws";
    static Dictionary<string, Delegate> _responseHandlers = new Dictionary<string, Delegate>();

    public static async Task Connect()
    {
        if (string.IsNullOrEmpty(AuthenticationManager.CharacterName))
        {
            Debug.LogError("No character name");
            if (Application.isEditor)
                EditorApplication.ExitPlaymode();
            else
                Application.Quit();

            return;
        }

        webSocket = new ClientWebSocket();
        Uri wsUri = new Uri(URL + $"?character={AuthenticationManager.CharacterName}");
        try
        {
            webSocket.Options.SetRequestHeader("Origin", "http://example.com");
            await webSocket.ConnectAsync(wsUri, CancellationToken.None);
            ReadLoopAsync();
        }
        catch (Exception ex)
        {
            Debug.Log(ex.Message);
        }
    }

    public static void SetHandler<TResponse>(Action<TResponse> responseHandler, string type)
        where TResponse : DTO
    {
        _responseHandlers[type] = responseHandler;
    }

    public static async void SendMessage(string messageJson)
    {
        var encodedMessage = Encoding.UTF8.GetBytes(messageJson);
        var wsBuffer = new ArraySegment<Byte>(encodedMessage, 0, encodedMessage.Length);

        try
        {
            await webSocket.SendAsync(
                wsBuffer,
                WebSocketMessageType.Text,
                true,
                CancellationToken.None
            );
        }
        catch (Exception ex)
        {
            Debug.Log(ex.Message);
        }
    }

    private static async void ReadLoopAsync()
    {
        List<byte> messageBytes = new List<byte>();
        byte[] receiveBuffer = new byte[1024];

        JsonSerializerSettings settings = new JsonSerializerSettings();
        settings.Converters.Add(new WebSocketMessageConverter());

        while (webSocket.State == WebSocketState.Open)
        {
            WebSocketReceiveResult result = null;
            do
            {
                result = await webSocket.ReceiveAsync(
                    new ArraySegment<byte>(receiveBuffer),
                    CancellationToken.None
                );
                for (int i = 0; i < result.Count; i++)
                {
                    messageBytes.Add(receiveBuffer[i]);
                }
            } while (!result.EndOfMessage);

            if (result.MessageType == WebSocketMessageType.Binary)
            {
                string responseJson = Encoding.UTF8.GetString(messageBytes.ToArray());
                WebSocketMessage response = JsonConvert.DeserializeObject<WebSocketMessage>(
                    responseJson,
                    settings
                );
                messageBytes.Clear();

                if (_responseHandlers.ContainsKey(response.ActionType))
                {
                    var handler = _responseHandlers[response.ActionType];
                    handler.DynamicInvoke(response.Body);
                }
            }
            else if (result.MessageType == WebSocketMessageType.Close)
            {
                Debug.Log("You have been kicked from the server");
                break;
            }
            else
            {
                Debug.LogError("Message from server is not text");
            }
        }
    }

    public static void ClearHandlers()
    {
        _responseHandlers = new Dictionary<string, Delegate>();
    }

    public static async Task Disconnect()
    {
        await webSocket.CloseAsync(
            WebSocketCloseStatus.NormalClosure,
            string.Empty,
            CancellationToken.None
        );
    }
}
