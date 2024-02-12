using System;
using System.Text;
using System.Threading;
using System.Net.WebSockets;
using UnityEngine;

public static class WebSocketConnection 
{
    static ClientWebSocket webSocket;
    static String URL = "ws://localhost:3009/ws";

    public static async void Connect() {
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
}