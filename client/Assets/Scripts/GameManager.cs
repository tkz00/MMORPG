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
	NametagBehaviour nametagBehaviour;

	private string mainPlayerID;

	Dictionary<string, PlayerMovement> players = new Dictionary<string, PlayerMovement>();
	
	async void Awake() {
		WebSocketConnection.SetHandler<string>(new Action<string>((playerId) => {
			this.mainPlayerID = playerId;
			Debug.Log(this.mainPlayerID);
		}));

		WebSocketConnection.SetHandler<GameStateDTO>(new Action<GameStateDTO>((gameState) => onGameStateUpdate(gameState)));

		await WebSocketConnection.Connect();
    }

	void Update() {
		nametagBehaviour.AssignNametags(players);
	}
	void movePlayers(List<PlayerDTO> playerDTOS) {
        foreach(PlayerDTO player in playerDTOS) {
			bool playerExists = this.players.Any(playerMovement => playerMovement.Key == player.Id);
            if(playerExists) {
            	this.players[player.Id].Move(new Vector3(player.Position.x, 0, player.Position.z));
            } else {
            	GameObject newPlayerGO = Instantiate(this.playerPrefab, new Vector3(player.Position.x, 1, player.Position.z), Quaternion.identity);
            	players[player.Id] = newPlayerGO.GetComponent<PlayerMovement>();
            	newPlayerGO.GetComponent<MeshRenderer>().material.color = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);

            	if(this.mainPlayerID == player.Id) {
            		newPlayerGO.AddComponent<MainPlayer>();
            		this.cinemachineVirtualCamera.Follow = newPlayerGO.transform;
            	}
            }
        }
	}

	void destroyPlayers(string[] playerIdsToDestroy) {
        foreach(string playerIdToDestroy in playerIdsToDestroy) {
            Destroy(this.players[playerIdToDestroy].gameObject);
            this.players.Remove(playerIdToDestroy);
        }
	}

	void onGameStateUpdate(GameStateDTO gameState) {
        string[] currentPlayerIds = this.players.Keys.ToArray();
		HashSet<string> playersToDestroy = new HashSet<string>(currentPlayerIds.Except(gameState.Players.Select(player => player.Id)));

		destroyPlayers(playersToDestroy.ToArray());
		movePlayers(gameState.Players);
	}

}
