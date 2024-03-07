using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using UnityEngine;
using UnityEngine.InputSystem;

public class GameManager : MonoBehaviour
{
	[SerializeField]
	GameObject playerPrefab;

	private string playerID;

	List<PlayerMovement> players = new List<PlayerMovement>();

	// los colores
	// despawnear
	// seguir con la camara
	
	async void Awake() {
		WebSocketConnection.SetHandler<string>(new Action<string>((playerId) => {
			this.playerID = playerId;
			Debug.Log(this.playerID);
		}));

		WebSocketConnection.SetHandler<GameStateDTO>(new Action<GameStateDTO>((gameState) => {
			foreach(PlayerDTO player in gameState.Players) {
				if(this.players.Any(playerMovement => playerMovement.playerID == player.Id)) {
					this.players.Find(playerMovement => playerMovement.playerID == player.Id).Move(new Vector3(player.Position.x, 0, player.Position.z));
				} else {
					GameObject newPlayerGO = Instantiate(this.playerPrefab, new Vector3(player.Position.x, 1, player.Position.z), Quaternion.identity);
					PlayerMovement playerMovement = newPlayerGO.GetComponent<PlayerMovement>();
					playerMovement.playerID = player.Id;
					players.Add(playerMovement);
					newPlayerGO.GetComponent<MeshRenderer>().material.color = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);

					if(this.playerID == player.Id) {
						newPlayerGO.AddComponent<MainPlayer>();
					}
				}
			}
		}));

		await WebSocketConnection.Connect();
    }
}
