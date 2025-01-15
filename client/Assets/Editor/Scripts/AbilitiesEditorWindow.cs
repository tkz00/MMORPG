using System.Collections.Generic;
using Newtonsoft.Json;
using UnityEditor;
using UnityEngine;

public class AbilitiesEditorWindow : EditorWindow
{
    private string responseText = "Response will appear here";

    [MenuItem("Window/Abilities Editor Window")]
    public static void ShowWindow()
    {
        GetWindow<AbilitiesEditorWindow>("Abilities Editor Window");
    }

    private void OnGUI()
    {
        if (GUILayout.Button("Make GET Request"))
        {
            MakeGetRequest();
        }

        GUILayout.Label("Response:");
        GUILayout.TextArea(responseText, GUILayout.Height(200));
    }

    private void MakeGetRequest()
    {
        // Example URL - replace with your actual backend URL
        string url = "http://0.0.0.0:8080/abilities";

        // Start the request
        var request = UnityEngine.Networking.UnityWebRequest.Get(url);
        request.SendWebRequest().completed += (asyncOperation) =>
        {
            if (request.result == UnityEngine.Networking.UnityWebRequest.Result.Success)
            {
                string jsonResponse = request.downloadHandler.text;
                // Deserialize JSON string to object
                Dictionary<string, AbilityDTO> abilities = JsonConvert.DeserializeObject<Dictionary<string, AbilityDTO>>(
                    JsonConvert.DeserializeObject<dynamic>(jsonResponse).abilities.ToString()
                );

                foreach (KeyValuePair<string, AbilityDTO> kyp in abilities)
                {
                    Debug.Log(kyp.Value.name);
                }
            }
            else
            {
                responseText = "Error: " + request.error;
            }
            Repaint(); // Update the window to show the new response
        };
    }
}
