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
    GameObject npcPrefab;

    [SerializeField]
    CinemachineVirtualCamera cinemachineVirtualCamera;

    bool hitboxOn = false;

    private string mainPlayerID;
    private AbilityDTO[] mainPlayerAbilities; // Naughty naughty (in Borat's voice)

    [SerializeField]
    public List<Ability> availableAbilities;

    [SerializeField]
    public AbilitiesPanel abilitiesPanel;

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

        abilitiesPanel.UpdatePlayerPanel(
            gameState.players.Find(player => player.id == mainPlayerID)
        );
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
            Character player;
            bool playerExists = this.players.Any(
                playerMovement => playerMovement.Key == playerDTO.id
            );
            if (playerExists)
            {
                player = players[playerDTO.id];
                player.Movement.Move(new Vector3(playerDTO.position.x, 0, playerDTO.position.z));
            }
            else
            {
                player = CreatePlayer(playerDTO);
            }
            player.UpdateHealth(playerDTO.currentHealth, playerDTO.maxHealth);

            switch (playerDTO.executingAction)
            {
                case ExecutingAction.Attacking:
                    player.Movement.TriggerWalkingAnimation(false);
                    player.Movement.AttackAnimation();
                    break;
                case ExecutingAction.CastingHeal:
                    player.Movement.TriggerWalkingAnimation(false);
                    player.Movement.HealAnimation();
                    break;
                case ExecutingAction.Moving:
                    player.Movement.TriggerWalkingAnimation(true);
                    break;
                default:
                    player.Movement.TriggerWalkingAnimation(false);
                    break;
            }

            player.SetScale(playerDTO.radius * 2);
            player.SetHitbox(hitboxOn);
        }
    }

    private Character CreatePlayer(CharacterDTO playerDTO)
    {
        Character player;
        Color playerColor = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);
        GameObject newPlayerGO = Instantiate(
            this.playerPrefab,
            new Vector3(playerDTO.position.x, 0, playerDTO.position.z),
            Quaternion.identity
        );
        player = newPlayerGO.GetComponent<Character>();
        player.id = playerDTO.id;
        player.SetNpcName(player.id);
        player.SetHealthBarColor(playerColor);
        players[playerDTO.id] = player;
        if (this.mainPlayerID == player.id)
        {
            MainPlayer mainPlayer = newPlayerGO.AddComponent<MainPlayer>();
            mainPlayer.abilitiesPanel = abilitiesPanel;
            mainPlayer.InitAbilities(mainPlayerAbilities, availableAbilities);
            this.cinemachineVirtualCamera.Follow = newPlayerGO.transform;
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

            switch (npcDTO.executingAction)
            {
                case ExecutingAction.Attacking:
                    npc.Movement.TriggerWalkingAnimation(false);
                    npc.Movement.AttackAnimation();
                    break;
                case ExecutingAction.CastingHeal:
                    npc.Movement.TriggerWalkingAnimation(false);
                    npc.Movement.HealAnimation();
                    break;
                case ExecutingAction.Moving:
                    npc.Movement.TriggerWalkingAnimation(true);
                    break;
                default:
                    npc.Movement.TriggerWalkingAnimation(false);
                    break;
            }

            npc.SetScale(npcDTO.radius * 2);
            npc.SetHitbox(hitboxOn);
        }
    }

    private Character CreateNPC(CharacterDTO npcDTO)
    {
        Character npc;
        Color npcColor = UnityEngine.Random.ColorHSV(0f, 1f, 1f, 1f, 0.5f, 1f);
        GameObject newNpcGO = Instantiate(
            this.npcPrefab,
            new Vector3(npcDTO.position.x, 0, npcDTO.position.z),
            Quaternion.identity
        );
        npc = newNpcGO.GetComponent<Character>();
        npc.id = npcDTO.id;
        npc.SetNpcName(npc.id);
        npc.SetHealthBarColor(npcColor);
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
    }
}
