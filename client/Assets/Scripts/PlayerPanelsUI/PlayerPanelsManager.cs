using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class PlayerPanelsManager : MonoBehaviour
{
    [SerializeField]
    Camera gameCamera;

    [SerializeField]
    Canvas gameCanvas;

	// It would be nice to use a pooler here?

    [SerializeField]
    GameObject playerPanelPrefab;

    Dictionary<string, GameObject> playerPanels = new Dictionary<string, GameObject>();

    public void UpdatePanels(Dictionary<string, Player> players) {
        foreach (KeyValuePair<string, Player> player in players) {
            GameObject playerPanelGO;
            if (!playerPanels.TryGetValue(player.Key, out playerPanelGO)) {
                playerPanelGO = Instantiate(playerPanelPrefab, gameCanvas.transform);
                // var playerColor = player.Value.gameObject.GetComponent<MeshRenderer>().material.color;
                // nametagComponent.color = playerColor;
				if (playerPanelGO.TryGetComponent(out PlayerPanel playerPanel))
				{
					playerPanel.Initialize(player.Key, player.Value.Stats.CurrentHealth, player.Value.Stats.MaxHealth);
					player.Value.onHealthChanged += playerPanel.UpdateHealthBar;
				}
				else {
					Debug.LogError("No PlayerPanel script in instantiated player panel prefab");
				}
                playerPanels[player.Key] = playerPanelGO;
            }
            Vector3 playerPosition = player.Value.transform.position;
            playerPanelGO.transform.position = gameCamera.WorldToScreenPoint(playerPosition) + new Vector3(0, 75, 0);
        }
    }
}
