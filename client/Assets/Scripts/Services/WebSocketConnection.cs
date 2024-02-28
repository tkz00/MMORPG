using System;
using System.Text;
using System.Threading;
using System.Threading.Tasks;
using System.Net.WebSockets;
using System.Collections.Generic;
using UnityEngine;

public static class WebSocketConnection 
{
    static ClientWebSocket webSocket;
    static String URL = "ws://localhost:3009/ws";

    public static async Task Connect() {
        // handshake
        webSocket = new ClientWebSocket();
        Uri wsUri = new Uri(URL);
        try {
            webSocket.Options.SetRequestHeader("Origin", "http://example.com");
            await webSocket.ConnectAsync(wsUri, CancellationToken.None);
            Debug.Log("conexion satisfactoria");
        } catch(Exception ex) {
            Debug.Log(ex.Message);
        }
    } 

    public static async void SendPosition(Vector3 position) {
        Debug.Log("sape");
        float x = position.x, z = position.z;
        string message = $"{x}, {z}";
        var encodedMessage = Encoding.UTF8.GetBytes(message);
        var wsBuffer = new ArraySegment<Byte>(encodedMessage, 0, encodedMessage.Length);

        try {
            await webSocket.SendAsync(wsBuffer, WebSocketMessageType.Text, true, CancellationToken.None);
            Debug.Log("Posicion enviada");
        } catch(Exception ex) {
            Debug.Log(ex.Message);
        }
    }

    public static async Task<string> ReceivePlayerID() {
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
            
            if (result.MessageType == WebSocketMessageType.Text)
            {
                string playerID = Encoding.UTF8.GetString(messageBytes.ToArray());
                messageBytes.Clear();

                return playerID;
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
        return "Error receiving message";
    }

/*
    static async Task ReceiveMessages()
    {
        List<byte> messageBytes = new List<byte>();
        byte[] receiveBuffer = new byte[1024];
        while (webSocketConnection.State == WebSocketState.Open)
        {
            WebSocketReceiveResult result = null;
            do
            {
                result = await webSocketConnection.ReceiveAsync(new ArraySegment<byte>(receiveBuffer), CancellationToken.None);
                for (int i = 0; i < result.Count; i++)
                {
                    messageBytes.Add(receiveBuffer[i]);
                }
            }
            while (!result.EndOfMessage);
            
            if (result.MessageType == WebSocketMessageType.Binary)
            {
                var webSocketResponse = WebSocketResponse.Parser.ParseFrom(messageBytes.ToArray());
                Debug.Log("received user: " + webSocketResponse.User.Username);
                
                messageBytes.Clear();
            }
            else if (result.MessageType == WebSocketMessageType.Close)
            {
                await webSocketConnection.CloseAsync(WebSocketCloseStatus.NormalClosure, string.Empty, CancellationToken.None);
            }
            else
            {
                Debug.LogError("Message from server is not binary");
            }
        }
    }
    */
}