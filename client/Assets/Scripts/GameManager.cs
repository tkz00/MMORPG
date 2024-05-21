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
	GameObject projectilePrefab;

	[SerializeField]
	CinemachineVirtualCamera cinemachineVirtualCamera;
	// [SerializeField]
	// PlayerPanelsManager PlayerPanelsManager;

	private string mainPlayerID;

	Dictionary<string, Player> players = new Dictionary<string, Player>();
	Dictionary<string, Projectile> projectiles = new Dictionary<string, Projectile>();

	async void Awake()
	{
		WebSocketConnection.SetHandler<PlayerDTO>(new Action<PlayerDTO>((playerDTO) =>
		{
			this.mainPlayerID = playerDTO.id;
			Debug.Log(this.mainPlayerID);
		}), "Player");

		WebSocketConnection.SetHandler<GameStateDTO>(new Action<GameStateDTO>((gameState) => OnGameStateUpdate(gameState)), "GameState");

		await WebSocketConnection.Connect();
	}

	void Update()
	{
		// PlayerPanelsManager.UpdatePanels(players);
	}

	void OnDestroy()
	{
		WebSocketConnection.ClearHandlers();
	}

	async void OnApplicationQuit()
	{
		await WebSocketConnection.Disconnect();
	}

	void OnGameStateUpdate(GameStateDTO gameState)
	{
		string[] currentPlayerIds = this.players.Keys.ToArray();
		string[] currentProjectileIds = this.projectiles.Keys.ToArray();
		HashSet<string> playersToDestroy = new HashSet<string>(currentPlayerIds.Except(gameState.players.Select(player => player.id)));
		HashSet<string> projectilesToDestroy = new HashSet<string>(currentProjectileIds.Except(gameState.projectiles.Select(projectile => projectile.id)));

		DestroyPlayers(playersToDestroy.ToArray());
		DestroyProjectiles(projectilesToDestroy.ToArray());
		UpdatePlayersPositions(gameState.players);
		UpdateProjectilesPositions(gameState.projectiles);
	}

	void UpdatePlayersPositions(List<PlayerDTO> playerDTOS)
	{
		foreach (PlayerDTO player in playerDTOS)
		{
			bool playerExists = this.players.Any(playerMovement => playerMovement.Key == player.id);
			if (playerExists)
			{
				this.players[player.id].Movement.Move(new Vector3(player.position.x, 0, player.position.z));
			}
			else
			{
				GameObject newPlayerGO = Instantiate(this.playerPrefab, new Vector3(player.position.x, 0, player.position.z), Quaternion.identity);
				players[player.id] = newPlayerGO.GetComponentInChildren<Player>();
				Color playerColor = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);
				newPlayerGO.GetComponentInChildren<MeshRenderer>().material.color = playerColor;
				players[player.id].SetPlayerName(player.id);
				players[player.id].SetHealthBarColor(playerColor);
				if (this.mainPlayerID == player.id)
				{
					newPlayerGO.AddComponent<MainPlayer>();
					this.cinemachineVirtualCamera.Follow = newPlayerGO.transform;
				}
			}
			players[player.id].UpdateHealth(player.currentHealth, player.maxHealth);
		}
	}

	void UpdateProjectilesPositions(List<ProjectileDTO> projectileDTOS)
	{
		foreach (ProjectileDTO projectile in projectileDTOS)
		{
			bool projectileExists = this.projectiles.Any(p => p.Key == projectile.id);
			if (projectileExists)
			{
				this.projectiles[projectile.id].Move(new Vector3(projectile.position.x, 0, projectile.position.z));
			}
			else
			{
				GameObject newProjectileGO = Instantiate(this.projectilePrefab, new Vector3(projectile.position.x, 1, projectile.position.z), Quaternion.identity);
				projectiles[projectile.id] = newProjectileGO.GetComponent<Projectile>();
				newProjectileGO.GetComponent<MeshRenderer>().material.color = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);
			}
		}
	}

	void DestroyPlayers(string[] playerIdsToDestroy)
	{
		foreach (string playerIdToDestroy in playerIdsToDestroy)
		{
			Destroy(this.players[playerIdToDestroy].gameObject);
			this.players.Remove(playerIdToDestroy);
		}
	}

	void DestroyProjectiles(string[] projectilesToDestroy)
	{
		foreach (string projectileId in projectilesToDestroy)
		{
			Destroy(this.projectiles[projectileId].gameObject);
			this.projectiles.Remove(projectileId);
		}
	}
}
