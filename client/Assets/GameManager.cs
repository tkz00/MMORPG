using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using Cinemachine;
using UnityEngine;
using UnityEngine.InputSystem;

public class GameManager : MonoBehaviour
{
	[SerializeField]
	GameObject playerPrefab;

	[SerializeField]
	CinemachineVirtualCamera cinemachineVirtualCamera;

	private string playerID;

	Dictionary<string, PlayerMovement> players = new Dictionary<string, PlayerMovement>();
	
	// despawnear

	async void Awake() {
		WebSocketConnection.SetHandler<string>(new Action<string>((playerId) => {
			this.playerID = playerId;
			Debug.Log(this.playerID);
		}));

		WebSocketConnection.SetHandler<GameStateDTO>(new Action<GameStateDTO>((gameState) => {
			string[] currentPlayerIds = this.players.Keys.ToArray();
			string[] playersToDestroy = currentPlayerIds.Except(gameState.Players.Select(player => player.Id).ToArray()).ToArray();
			foreach(string playerIdToDestroy in playersToDestroy) {
				Destroy(this.players[playerIdToDestroy].gameObject);
				this.players.Remove(playerIdToDestroy);
			}
			foreach(PlayerDTO player in gameState.Players) {
				if(this.players.Any(playerMovement => playerMovement.Key == player.Id)) {
					this.players[player.Id].Move(new Vector3(player.Position.x, 0, player.Position.z));
				} else {
					GameObject newPlayerGO = Instantiate(this.playerPrefab, new Vector3(player.Position.x, 1, player.Position.z), Quaternion.identity);
					PlayerMovement playerMovement = newPlayerGO.GetComponent<PlayerMovement>();
					players[player.Id] = playerMovement;
					newPlayerGO.GetComponent<MeshRenderer>().material.color = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);

					if(this.playerID == player.Id) {
						newPlayerGO.AddComponent<MainPlayer>();
						this.cinemachineVirtualCamera.Follow = newPlayerGO.transform;
					}
				}
			}
		}));

		await WebSocketConnection.Connect();
    }
}
