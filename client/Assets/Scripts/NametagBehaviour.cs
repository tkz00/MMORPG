using System.Collections;
using System.Collections.Generic;
using TMPro;
using UnityEngine;

public class NametagBehaviour : MonoBehaviour
{
    [SerializeField]
    Camera gameCamera;
    [SerializeField]
    Canvas gameCanvas;
    [SerializeField]
    GameObject nametagPrefab;
    Dictionary<string, GameObject> nametags = new Dictionary<string, GameObject>();

    public void AssignNametags(Dictionary<string, PlayerMovement> players) {
        foreach (KeyValuePair<string, PlayerMovement> player in players) {
            GameObject nametagGO;
            if (!nametags.TryGetValue(player.Key, out nametagGO)) {
                nametagGO = Instantiate(nametagPrefab, gameCanvas.transform);
                var playerColor = player.Value.gameObject.GetComponent<MeshRenderer>().material.color;
                var nametagComponent = nametagGO.GetComponent<TMP_Text>();
                nametagComponent.text = player.Key + "\n" + player.Value.currentHealth + "/" + player.Value.maxHealth;
                nametagComponent.color = playerColor;
                nametags[player.Key] = nametagGO;
            }
            Vector3 playerPosition = player.Value.transform.position;
            nametagGO.transform.position = gameCamera.WorldToScreenPoint(playerPosition) + new Vector3(0, 75, 0);
        }
    }
}
