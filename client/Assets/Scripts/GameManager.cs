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
	[SerializeField]
	PlayerPanelsManager PlayerPanelsManager;

	private string mainPlayerID;

	Dictionary<string, Player> players = new Dictionary<string, Player>();

	async void Awake() {
		WebSocketConnection.SetHandler<PlayerDTO>(new Action<PlayerDTO>((playerDTO) => {
			this.mainPlayerID = playerDTO.id;
			Debug.Log(this.mainPlayerID);
		}), "Player");

		WebSocketConnection.SetHandler<GameStateDTO>(new Action<GameStateDTO>((gameState) => OnGameStateUpdate(gameState)), "GameState");

		await WebSocketConnection.Connect();
    }
	
	void Update() {
		PlayerPanelsManager.UpdatePanels(players);
	}
	
	void OnDestroy() {
		WebSocketConnection.ClearHandlers();
	}

	async void OnApplicationQuit() {
		await WebSocketConnection.Disconnect();
	}

	void OnGameStateUpdate(GameStateDTO gameState) {
        string[] currentPlayerIds = this.players.Keys.ToArray();
		HashSet<string> playersToDestroy = new HashSet<string>(currentPlayerIds.Except(gameState.players.Select(player => player.id)));

		DestroyPlayers(playersToDestroy.ToArray());
		UpdatePlayersPositions(gameState.players);
	}

	void UpdatePlayersPositions(List<PlayerDTO> playerDTOS) {
        foreach(PlayerDTO player in playerDTOS) {
			bool playerExists = this.players.Any(playerMovement => playerMovement.Key == player.id);
            if(playerExists) {
            	this.players[player.id].Movement.Move(new Vector3(player.position.x, 0, player.position.z));
            } else {
            	GameObject newPlayerGO = Instantiate(this.playerPrefab, new Vector3(player.position.x, 1, player.position.z), Quaternion.identity);
            	players[player.id] = newPlayerGO.GetComponent<Player>();
            	newPlayerGO.GetComponent<MeshRenderer>().material.color = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);

            	if(this.mainPlayerID == player.id) {
            		newPlayerGO.AddComponent<MainPlayer>();
            		this.cinemachineVirtualCamera.Follow = newPlayerGO.transform;
            	}
            }
			players[player.id].UpdateHealth(player.currentHealth, player.maxHealth);
        }
	}

	void DestroyPlayers(string[] playerIdsToDestroy) {
        foreach(string playerIdToDestroy in playerIdsToDestroy) {
            Destroy(this.players[playerIdToDestroy].gameObject);
            this.players.Remove(playerIdToDestroy);
        }
	}
}
