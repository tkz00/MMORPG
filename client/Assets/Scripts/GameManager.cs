using System;
using System.Collections.Generic;
using System.Linq;
using Cinemachine;
using Newtonsoft.Json;
using UnityEngine;
using UnityEngine.SceneManagement;

public class GameManager : MonoBehaviour
{
    [SerializeField] GameObject playerPrefab;

    [SerializeField] GameObject projectilePrefab;
    [SerializeField] GameObject aoePrefab;

    [SerializeField] GameObject npcPrefab;

    [SerializeField] CinemachineVirtualCamera cinemachineVirtualCamera;

    static GameManager gameManager;

    bool hitboxOn = false;

    private static string mainPlayerID;
    public static string MainPlayerID { get { return mainPlayerID; } }
    private AbilityDTO[] mainPlayerAbilities; // Naughty naughty (in Borat's voice)

    // [SerializeField] public static List<Ability> availableAbilities;
    [SerializeField] AbilitiesContainer abilitiesContainer;

    [SerializeField] public AbilitiesPanel abilitiesPanel;

    [SerializeField] GameObject deathSplash;

    [SerializeField] Inventory inventory;
    [SerializeField] StatsPanel statsPanel;


    Dictionary<string, Character> players = new Dictionary<string, Character>();
    Dictionary<string, Projectile> projectiles = new Dictionary<string, Projectile>();
    Dictionary<string, Character> npcs = new Dictionary<string, Character>();
    Dictionary<string, AoE> areaEffects = new Dictionary<string, AoE>();

    async void Awake()
    {
        if (gameManager == null)
            gameManager = this;

        WebSocketConnection.SetHandler<CharacterDTO>(
            new Action<CharacterDTO>(
                (playerDTO) =>
                {
                    mainPlayerID = playerDTO.id;
                    Debug.Log($"Player connected, id: {mainPlayerID}");
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
        if (WebSocketConnection.IsConnected)
            await WebSocketConnection.Disconnect();
        gameManager = null;
    }

    public async void Disconnect()
    {
        if (WebSocketConnection.IsConnected)
            await WebSocketConnection.Disconnect();
        gameManager = null;
        SceneManager.LoadScene("AuthenticationScene", LoadSceneMode.Single);
    }

    void OnGameStateUpdate(GameStateDTO gameState)
    {
        UpdatePlayers(gameState.players);
        UpdateProjectiles(gameState.projectiles);
        UpdateNPCs(gameState.npcs);
        UpdateAreaEffects(gameState.aoEs);

        CharacterDTO mainPlayer = gameState.players.Find(player => player.id == mainPlayerID);
        abilitiesPanel.UpdatePlayerPanel(mainPlayer);
        inventory.UpdateInventory(mainPlayer.inventory);
        if (mainPlayer.stats != null) statsPanel.UpdatePanel(mainPlayer.stats);
    }

    private void UpdateAreaEffects(List<AreaEffectDTO> AoEs)
    {
        if (AoEs == null || AoEs.Count == 0)
        {
            DestroyAoEs(areaEffects.Keys.ToArray());
            return;
        }
        string[] currentAoEsIds = areaEffects.Keys.ToArray();
        HashSet<string> AoEsToDestroy = new HashSet<string>(
            currentAoEsIds.Except(AoEs.Select(AoE => AoE.id))
        );
        DestroyAoEs(AoEsToDestroy.ToArray());

        foreach (AreaEffectDTO areaEffectDTO in AoEs)
        {
            bool areaEffectExists = areaEffects.Any(aoe => aoe.Key == areaEffectDTO.id);
            if (!areaEffectExists)
            {
                GameObject newAoEGO = Instantiate(
                                    this.aoePrefab,
                                    new Vector3(areaEffectDTO.position.x, this.aoePrefab.transform.position.y, areaEffectDTO.position.z),
                                    Quaternion.identity
                                );
                areaEffects[areaEffectDTO.id] = newAoEGO.GetComponent<AoE>();
                newAoEGO.GetComponentInChildren<MeshRenderer>().material.color = Color.red;
                areaEffects[areaEffectDTO.id].transform.localScale = Vector3.one * (areaEffectDTO.radius * 2);
            }
        }
    }

    void DestroyAoEs(string[] AoEsToDestroy)
    {
        foreach (string AoEId in AoEsToDestroy)
        {
            Destroy(areaEffects[AoEId].gameObject);
            areaEffects.Remove(AoEId);
        }
    }

    #region Players

    void UpdatePlayers(List<CharacterDTO> playerDTOS)
    {
        string[] currentPlayerIds = players.Keys.ToArray();
        HashSet<string> playersToDestroy = new HashSet<string>(
            currentPlayerIds.Except(playerDTOS.Select(player => player.id))
        );
        DestroyPlayers(playersToDestroy.ToArray());
        UpdatePlayersPositions(playerDTOS);
        UpdatePlayersEquipment(playerDTOS);
    }

    void UpdatePlayersPositions(List<CharacterDTO> playerDTOS)
    {
        foreach (CharacterDTO playerDTO in playerDTOS)
        {
            Character player;
            bool playerExists = players.Any(character => character.Key == playerDTO.id);
            if (playerExists)
                player = GetPlayer(playerDTO.id);
            else
                player = CreatePlayer(playerDTO);
            if (playerDTO.position != null) player.Movement.Move(new Vector3(playerDTO.position.x, 0, playerDTO.position.z));
            player.UpdateHealth(playerDTO.currentHealth, playerDTO.maxHealth);
            UpdateCharacterAnimations(playerDTO, player);

            if (playerDTO.radius != null) player.SetScale((float)(playerDTO.radius * 2));

            CheckDeathSplash(playerDTO);
        }
    }

    void UpdatePlayersEquipment(List<CharacterDTO> playerDTOS)
    {
        foreach (CharacterDTO playerDTO in playerDTOS)
        {
            if (playerDTO.inventory?.items == null)
                return;

            if (playerDTO.inventory?.items.Length == 0)
                return;

            Character player = GetPlayer(playerDTO.id);
            player.UpdateEquipmentVisuals(playerDTO.inventory.items);
        }
    }

    Character GetPlayer(string playerId)
    {
        return players[playerId];
    }

    public static bool MainPlayerIsAlive()
    {
        return gameManager.players[MainPlayerID].IsAlive;
    }

    private void CheckDeathSplash(CharacterDTO playerDTO)
    {
        if (mainPlayerID == playerDTO.id)
        {
            if (playerDTO.currentHealth <= 0 && !deathSplash.activeSelf)
                deathSplash.SetActive(true);
            else if (playerDTO.currentHealth > 0 && deathSplash.activeSelf)
                deathSplash.SetActive(false);
        }
    }

    private void UpdateCharacterAnimations(CharacterDTO characterDTO, Character character)
    {
        if (character.Stats.CurrentHealth > 0)
        {
            if (characterDTO.action != null) character.HandleActionFeedback((CharacterAction)characterDTO.action, characterDTO.executionDurationMs);
            if (characterDTO.direction != null) character.Movement.RotateTowards(characterDTO.direction);
            return;
        }

        if (character.IsAlive)
            character.TriggerDeath();
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
        if (mainPlayerID == player.id)
        {
            MainPlayer mainPlayer = newPlayerGO.AddComponent<MainPlayer>();
            mainPlayer.abilitiesPanel = abilitiesPanel;
            mainPlayer.InitAbilities(mainPlayerAbilities, abilitiesContainer.GetAvailableAbilities());
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
            Destroy(players[playerIdToDestroy].gameObject);
            players.Remove(playerIdToDestroy);
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

    private void UpdateProjectiles(List<ProjectileDTO> projectilesDTOS)
    {
        if (projectilesDTOS == null || projectilesDTOS.Count == 0)
        {
            DestroyProjectiles(projectiles.Keys.ToArray());
            return;
        }
        string[] currentProjectileIds = projectiles.Keys.ToArray();
        HashSet<string> projectilesToDestroy = new HashSet<string>(
            currentProjectileIds.Except(projectilesDTOS.Select(projectile => projectile.id))
        );
        DestroyProjectiles(projectilesToDestroy.ToArray());
        UpdateProjectilesPositions(projectilesDTOS);
    }

    void UpdateProjectilesPositions(List<ProjectileDTO> projectileDTOS)
    {
        foreach (ProjectileDTO projectile in projectileDTOS)
        {
            bool projectileExists = projectiles.Any(p => p.Key == projectile.id);
            if (projectileExists)
            {
                if (projectile.position != null)
                    projectiles[projectile.id].Move(
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

            if (projectile.radius != null)
                projectiles[projectile.id].SetScale((float)(projectile.radius * 2));

            if (projectile.state == State.Hit)
                projectiles[projectile.id].TriggerHit();
        }
    }

    void DestroyProjectiles(string[] projectilesToDestroy)
    {
        foreach (string projectileId in projectilesToDestroy)
        {
            Destroy(projectiles[projectileId].gameObject);
            projectiles.Remove(projectileId);
        }
    }

    #endregion

    #region NPCs

    private void UpdateNPCs(List<CharacterDTO> npcsDTOS)
    {
        if (npcsDTOS == null || npcsDTOS.Count == 0)
        {
            DestroyNPCs(npcs.Keys.ToArray());
            return;
        }
        string[] currentNPCsIds = npcs.Keys.ToArray();
        HashSet<string> npcsToDestroy = new HashSet<string>(
            currentNPCsIds.Except(npcsDTOS.Select(npc => npc.id))
        );
        DestroyNPCs(npcsToDestroy.ToArray());
        UpdateNPCsPositions(npcsDTOS);
    }

    void UpdateNPCsPositions(List<CharacterDTO> npcDTOS)
    {
        foreach (CharacterDTO npcDTO in npcDTOS)
        {
            Character npc;
            bool npcExists = npcs.Any(npcMovement => npcMovement.Key == npcDTO.id);
            if (npcExists)
            {
                npc = npcs[npcDTO.id];
                if (npcDTO.position != null) npc.Movement.Move(new Vector3(npcDTO.position.x, 0, npcDTO.position.z));
            }
            else
            {
                npc = CreateNPC(npcDTO);
            }
            npc.UpdateHealth(npcDTO.currentHealth, npcDTO.maxHealth);

            if (npcDTO.action != null) npc.HandleActionFeedback((CharacterAction)npcDTO.action, npcDTO.executionDurationMs);
            if (npcDTO.direction != null) npc.Movement.RotateTowards(npcDTO.direction);

            if (npcDTO.radius != null) npc.SetScale((float)(npcDTO.radius * 2));
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
            Destroy(npcs[npcIdToDestroy].gameObject);
            npcs.Remove(npcIdToDestroy);
        }
    }

    #endregion

    public void ToggleHitboxes()
    {
        hitboxOn = !hitboxOn;

        foreach (Character player in players.Values)
        {
            player.SetHitbox(hitboxOn);
        }

        foreach (Character player in npcs.Values)
        {
            player.SetHitbox(hitboxOn);
        }
    }
}
