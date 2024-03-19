using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class NametagBehaviour : MonoBehaviour
{
    [SerializeField]
    Camera gameCamera;

    [SerializeField]
    Canvas gameCanvas;

	// It would be nice to use a pooler here

    [SerializeField]
    PlayerPanel playerPanelPrafab;

    Dictionary<string, GameObject> playerPanels = new Dictionary<string, GameObject>();

    public void AssignNametags(Dictionary<string, Player> players) {
        foreach (KeyValuePair<string, Player> player in players) {
            GameObject playerPanelGO;
            if (!playerPanels.TryGetValue(player.Key, out playerPanelGO)) {
                playerPanelGO = Instantiate(playerPanelPrafab.gameObject, gameCanvas.transform);
                // var playerColor = player.Value.gameObject.GetComponent<MeshRenderer>().material.color;
				PlayerPanel playerPanel = playerPanelGO.GetComponent<PlayerPanel>();
				playerPanel.Initialize(player.Key, player.Value.Stats.CurrentHealth, player.Value.Stats.MaxHealth);
                // nametagComponent.color = playerColor;
                playerPanels[player.Key] = playerPanelGO;
            }
            Vector3 playerPosition = player.Value.transform.position;
            playerPanelGO.transform.position = gameCamera.WorldToScreenPoint(playerPosition) + new Vector3(0, 75, 0);
			// playerPanelGO.transform.position = gameCamera.WorldToScreenPoint(playerPosition);
        }
    }
}
