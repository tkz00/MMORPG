using System;
using System.Collections;
using System.Collections.Generic;
using System.Linq;
using Cinemachine;
using Newtonsoft.Json;
using UnityEngine;

public class GameManager : MonoBehaviour
{
    [SerializeField] GameObject playerPrefab;

    [SerializeField] GameObject projectilePrefab;

    [SerializeField] GameObject npcPrefab;

    [SerializeField] CinemachineVirtualCamera cinemachineVirtualCamera;

    bool hitboxOn = false;

    private string mainPlayerID;
    private AbilityDTO[] mainPlayerAbilities; // Naughty naughty (in Borat's voice)

    [SerializeField] public List<Ability> availableAbilities;

    [SerializeField] public AbilitiesPanel abilitiesPanel;

    [SerializeField] GameObject deathSplash;

    [SerializeField] Inventory inventory;


    Dictionary<string, Character> players = new Dictionary<string, Character>();
    Dictionary<string, Projectile> projectiles = new Dictionary<string, Projectile>();
    Dictionary<string, Character> npcs = new Dictionary<string, Character>();

    async void Awake()
    {
        WebSocketConnection.SetHandler<CharacterDTO>(
            new Action<CharacterDTO>(
                (playerDTO) =>
                {
                    this.mainPlayerID = playerDTO.id;
                    Debug.Log($"Player connected, id: {this.mainPlayerID}");
                    abilitiesPanel.Init(playerDTO.abilities);
                    mainPlayerAbilities = playerDTO.abilities;
                }
            ),
            "Player"
        );

        WebSocketConnection.SetHandler<GameStateDTO>(
            new Action<GameStateDTO>((gameState) => OnGameStateUpdate(gameState)),
            "GameState"
        );

        WebSocketConnection.SetHandler(
            new Action<CharacterDTO>(
                (playerDTO) =>
                {
                    Character player = players[playerDTO.id];
                    player.Respawn(new Vector3(playerDTO.position.x, 0, playerDTO.position.z));
                }
            ),
            "respawn"
        );

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
        UpdatePlayers(gameState.players);
        UpdateProjectiles(gameState.projectiles);
        UpdateNPCs(gameState.npcs);

        CharacterDTO mainPlayer = gameState.players.Find(player => player.id == mainPlayerID);

        abilitiesPanel.UpdatePlayerPanel(mainPlayer);

        inventory.UpdateInventory(mainPlayer.inventory);
    }

    #region Players

    private void UpdatePlayers(List<CharacterDTO> players)
    {
        string[] currentPlayerIds = this.players.Keys.ToArray();
        HashSet<string> playersToDestroy = new HashSet<string>(
            currentPlayerIds.Except(players.Select(player => player.id))
        );
        DestroyPlayers(playersToDestroy.ToArray());
        UpdatePlayersPositions(players);
    }

    void UpdatePlayersPositions(List<CharacterDTO> playerDTOS)
    {
        foreach (CharacterDTO playerDTO in playerDTOS)
        {
            Character player = GetPlayer(playerDTO);
            player.Movement.Move(new Vector3(playerDTO.position.x, 0, playerDTO.position.z));
            player.UpdateHealth(playerDTO.currentHealth, playerDTO.maxHealth);
            UpdateCharacterAnimations(playerDTO, player);

            player.SetScale(playerDTO.radius * 2);

            CheckDeathSplash(playerDTO);
        }
    }

    private Character GetPlayer(CharacterDTO playerDTO)
    {
        Character player;
        bool playerExists = this.players.Any(
            playerMovement => playerMovement.Key == playerDTO.id
        );
        if (playerExists)
        {
            player = players[playerDTO.id];
        }
        else
        {
            player = CreatePlayer(playerDTO);
        }

        return player;
    }

    private void CheckDeathSplash(CharacterDTO playerDTO)
    {
        if (this.mainPlayerID == playerDTO.id)
        {
            if (playerDTO.currentHealth <= 0 && !deathSplash.activeSelf)
            {
                deathSplash.SetActive(true);
            }
        }
    }

    private void UpdateCharacterAnimations(CharacterDTO characterDTO, Character character)
    {
        if (characterDTO.currentHealth > 0)
        {
            character.HandleActionFeedback(characterDTO.executingAction);

            return;
        }

        if (character.IsAlive)
        {
            character.TriggerDeath();
        }
    }

    private Character CreatePlayer(CharacterDTO playerDTO)
    {
        Character player;
        GameObject newPlayerGO = Instantiate(
            this.playerPrefab,
            new Vector3(playerDTO.position.x, 0, playerDTO.position.z),
            Quaternion.identity
        );
        player = newPlayerGO.GetComponent<Character>();
        player.id = playerDTO.id;
        player.SetCharacterName(player.id);
        player.SetHitbox(hitboxOn);
        players[playerDTO.id] = player;
        if (this.mainPlayerID == player.id)
        {
            MainPlayer mainPlayer = newPlayerGO.AddComponent<MainPlayer>();
            mainPlayer.abilitiesPanel = abilitiesPanel;
            mainPlayer.InitAbilities(mainPlayerAbilities, availableAbilities);
            this.cinemachineVirtualCamera.Follow = newPlayerGO.transform;
            Color playerColor = Color.green;
            player.SetHealthBarColor(playerColor);
        }
        else
        {
            Color playerColor = Color.red;
            player.SetHealthBarColor(playerColor);
        }
        return player;
    }

    void DestroyPlayers(string[] playerIdsToDestroy)
    {
        foreach (string playerIdToDestroy in playerIdsToDestroy)
        {
            Destroy(this.players[playerIdToDestroy].gameObject);
            this.players.Remove(playerIdToDestroy);
        }
    }

    // this should defenitely not be here
    public void Respawn()
    {
        WebSocketMessage message = new WebSocketMessage
        {
            ActionType = "respawn"
        };
        string jsonMessage = JsonConvert.SerializeObject(message);
        WebSocketConnection.SendMessage(jsonMessage);
    }
    #endregion

    #region Projectiles

    private void UpdateProjectiles(List<ProjectileDTO> projectiles)
    {
        string[] currentProjectileIds = this.projectiles.Keys.ToArray();
        HashSet<string> projectilesToDestroy = new HashSet<string>(
            currentProjectileIds.Except(projectiles.Select(projectile => projectile.id))
        );
        DestroyProjectiles(projectilesToDestroy.ToArray());
        UpdateProjectilesPositions(projectiles);
    }

    void UpdateProjectilesPositions(List<ProjectileDTO> projectileDTOS)
    {
        foreach (ProjectileDTO projectile in projectileDTOS)
        {
            bool projectileExists = this.projectiles.Any(p => p.Key == projectile.id);
            if (projectileExists)
            {
                this.projectiles[projectile.id].Move(
                    new Vector3(projectile.position.x, 0, projectile.position.z)
                );
            }
            else
            {
                GameObject newProjectileGO = Instantiate(
                    this.projectilePrefab,
                    new Vector3(projectile.position.x, 1, projectile.position.z),
                    Quaternion.identity
                );
                projectiles[projectile.id] = newProjectileGO.GetComponent<Projectile>();
                newProjectileGO.GetComponent<MeshRenderer>().material.color = Color.blue;
            }
            projectiles[projectile.id].SetScale(projectile.radius * 2);

            if (projectile.state == State.Hit)
            {
                projectiles[projectile.id].TriggerHit();
            }
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

    #endregion

    #region NPCs

    private void UpdateNPCs(List<CharacterDTO> npcs)
    {
        string[] currentNPCsIds = this.npcs.Keys.ToArray();
        HashSet<string> npcsToDestroy = new HashSet<string>(
            currentNPCsIds.Except(npcs.Select(npc => npc.id))
        );
        DestroyNPCs(npcsToDestroy.ToArray());
        UpdateNPCsPositions(npcs);
    }

    void UpdateNPCsPositions(List<CharacterDTO> npcDTOS)
    {
        foreach (CharacterDTO npcDTO in npcDTOS)
        {
            Character npc;
            bool npcExists = this.npcs.Any(npcMovement => npcMovement.Key == npcDTO.id);
            if (npcExists)
            {
                npc = npcs[npcDTO.id];
                npc.Movement.Move(new Vector3(npcDTO.position.x, 0, npcDTO.position.z));
            }
            else
            {
                npc = CreateNPC(npcDTO);
            }
            npc.UpdateHealth(npcDTO.currentHealth, npcDTO.maxHealth);

            npc.HandleActionFeedback(npcDTO.executingAction);

            npc.SetScale(npcDTO.radius * 2);
        }
    }

    private Character CreateNPC(CharacterDTO npcDTO)
    {
        Character npc;
        GameObject newNpcGO = Instantiate(
            this.npcPrefab,
            new Vector3(npcDTO.position.x, 0, npcDTO.position.z),
            Quaternion.identity
        );
        npc = newNpcGO.GetComponent<Character>();
        npc.id = npcDTO.id;
        npc.SetCharacterName(npc.id);
        Color npcColor = Color.yellow;
        npc.SetHealthBarColor(npcColor);
        npc.SetHitbox(hitboxOn);
        npcs[npcDTO.id] = npc;
        return npc;
    }

    void DestroyNPCs(string[] npcIdsToDestroy)
    {
        foreach (string npcIdToDestroy in npcIdsToDestroy)
        {
            Destroy(this.npcs[npcIdToDestroy].gameObject);
            this.npcs.Remove(npcIdToDestroy);
        }
    }

    #endregion

    public void ToggleHitboxes()
    {
        hitboxOn = !hitboxOn;

        foreach (Character player in this.players.Values)
        {
            player.SetHitbox(hitboxOn);
        }

        foreach (Character player in this.npcs.Values)
        {
            player.SetHitbox(hitboxOn);
        }
    }
}
