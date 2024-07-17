using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using Cinemachine;
using Unity.VisualScripting;
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

	bool hitboxOn = false;

	private string mainPlayerID;

	[SerializeField]
	public AbilitiesPanel abilitiesPanel;

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
		foreach (PlayerDTO playerDTO in playerDTOS)
		{
			Player player;
			bool playerExists = this.players.Any(playerMovement => playerMovement.Key == playerDTO.id);
			if (playerExists)
			{
				player = players[playerDTO.id];
				player.Movement.Move(new Vector3(playerDTO.position.x, 0, playerDTO.position.z));
			}
			else
			{
				Color playerColor = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);
				GameObject newPlayerGO = Instantiate(this.playerPrefab, new Vector3(playerDTO.position.x, 0, playerDTO.position.z), Quaternion.identity);
				player = newPlayerGO.GetComponent<Player>();
				player.id = playerDTO.id;
				player.SetPlayerName(player.id);
				player.SetHealthBarColor(playerColor);
            	players[playerDTO.id] = player;
				if (this.mainPlayerID == player.id)
				{
					MainPlayer mainPlayer = newPlayerGO.AddComponent<MainPlayer>();
					mainPlayer.abilitiesPanel = abilitiesPanel;
					this.cinemachineVirtualCamera.Follow = newPlayerGO.transform;
				}
			}
			player.UpdateHealth(playerDTO.currentHealth, playerDTO.maxHealth);
			
			switch(playerDTO.executingAction)
			{
				case ExecutingAction.Attacking:
					player.Movement.AttackAnimation();
					break;
				case ExecutingAction.CastingHeal:
					player.Movement.HealAnimation();
					break;
			}

			player.SetScale(playerDTO.radius * 2);
			player.SetHitbox(hitboxOn);
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
				newProjectileGO.GetComponent<MeshRenderer>().material.color = Color.blue;
			}
			projectiles[projectile.id].SetScale(projectile.radius * 2);

			if(projectile.state == State.Hit)
			{
				projectiles[projectile.id].TriggerHit();
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

	public void ToggleHitboxes()
	{
		hitboxOn = !hitboxOn;

		foreach(Player player in this.players.Values)
		{
			player.SetHitbox(hitboxOn);
		}
	}
}
