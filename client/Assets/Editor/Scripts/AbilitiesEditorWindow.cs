using System;
using System.Collections.Generic;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;
using UnityEditor;
using UnityEngine;
using UnityEngine.Networking;

public class AbilitiesEditorWindow : EditorWindow
{
    const string BACKEND_URL = "http://0.0.0.0:8080";

    private string responseText = string.Empty;
    private List<Configurator.AbilityDTO> abilitiesList = new List<Configurator.AbilityDTO>();
    private Vector2 scrollPosition;
    private Configurator.AbilityDTO selectedAbility = null; // Tracks the selected ability

    #region Editable fields
    string abilityName; // Very wrong but this acts as a flag for if the editable fields have been initialized
    int abilityCooldown;
    float abilityRange;
    #endregion

    [MenuItem("Window/Abilities Editor Window")]
    public static void ShowWindow()
    {
        GetWindow<AbilitiesEditorWindow>("Abilities Editor Window");
    }

    private void OnGUI()
    {
        GUILayout.Space(10);
        GUILayout.Label("Abilities Editor", EditorStyles.boldLabel);

        if (selectedAbility == null)
        {
            DrawListView(); // Show list view if no ability is selected
        }
        else
        {
            DrawDetailView(); // Show detail view for the selected ability
        }
    }

    private void DrawListView()
    {
        if (GUILayout.Button("Fetch Abilities", GUILayout.Height(30)))
        {
            FetchAbilities();
        }

        if (responseText != string.Empty)
        {
            GUILayout.Space(10);
            GUILayout.Label("Response:", EditorStyles.boldLabel);
            GUILayout.TextArea(responseText, GUILayout.Height(50));
        }

        GUILayout.Space(10);
        GUILayout.Label("Abilities List:", EditorStyles.boldLabel);

        if (abilitiesList.Count == 0)
        {
            GUILayout.Label("No abilities loaded. Fetch abilities to see the list.");
        }
        else
        {
            scrollPosition = GUILayout.BeginScrollView(scrollPosition, GUILayout.Height(300));

            foreach (var ability in abilitiesList)
            {
                if (GUILayout.Button(ability.name, GUILayout.Height(25)))
                {
                    selectedAbility = ability; // Set selected ability to show detail view
                }
            }

            GUILayout.EndScrollView();
        }
    }

    private async Task DrawDetailView()
    {
        if (abilityName == null) // Initialize editable fields only once
        {
            InitializeEditableFields();
        }

        // Back button to return to the list view
        if (GUILayout.Button("Back to List", GUILayout.Height(30)))
        {
            selectedAbility = null; // Clear selected ability to return to the list view
            abilityName = null; // Reset editable fields
            FetchAbilities();
            return;
        }

        GUILayout.Space(10);
        GUILayout.Label($"Details for: {selectedAbility.name}", EditorStyles.boldLabel);

        GUILayout.Label($"ID: {selectedAbility.id}");
        abilityName = EditorGUILayout.TextField("Name", abilityName);
        abilityRange = EditorGUILayout.FloatField("Range", abilityRange);
        abilityCooldown = EditorGUILayout.IntField("Cooldown", abilityCooldown);
        GUILayout.Label($"Targeting: {selectedAbility.targeting}");
        GUILayout.Label($"Character State: {selectedAbility.characterAction}");

        GUILayout.Space(10);
        GUILayout.Label("Mechanics:", EditorStyles.boldLabel);

        // Display mechanics in a scrollable area
        scrollPosition = GUILayout.BeginScrollView(scrollPosition, GUILayout.Height(200));
        foreach (var mechanic in selectedAbility.mechanics)
        {
            GUILayout.Label($"Mechanic Type: {mechanic.mechanicType}");

            if (mechanic.mechanicType == "create_projectile")
            {
                GUILayout.Label("Projectile On Hit Mechanics:", EditorStyles.boldLabel);
                var projectileParams = mechanic.@params as Configurator.AbilityDTO.CreateProjectileParams;
                foreach (var projectileOnHitMechanic in projectileParams.onHitMechanics)
                {
                    GUILayout.Label($"{projectileOnHitMechanic.mechanicType}");

                    if (projectileOnHitMechanic.mechanicType == "create_AoE")
                    {
                        GUILayout.Label("AoE On Hit Mechanics:", EditorStyles.boldLabel);
                        var aoEParams = projectileOnHitMechanic.@params as Configurator.AbilityDTO.CreateAoEParams;
                        foreach (var aoEOnHitMechanic in aoEParams.onHitMechanics)
                        {
                            GUILayout.Label($"{aoEOnHitMechanic.mechanicType}");

                            if (aoEOnHitMechanic.mechanicType == "damage")
                            {
                                var damageParams = aoEOnHitMechanic.@params as Configurator.AbilityDTO.DamageParams;
                                GUILayout.Label($"Damage: {damageParams.amount}");
                            }
                        }
                        // GUILayout.Label($"{aoEParams.radius}");
                        // GUILayout.Label($"{aoEParams.durationMs}");
                    }
                }
            }

            GUILayout.Space(10);
        }
        GUILayout.EndScrollView();

        if (GUILayout.Button("Save Changes"))
        {
            await PatchAbility();
        }
    }

    private void InitializeEditableFields()
    {
        abilityName = selectedAbility.name;
        abilityCooldown = selectedAbility.cooldown;
        abilityRange = selectedAbility.range;
    }

    private void FetchAbilities()
    {
        var request = UnityWebRequest.Get(BACKEND_URL + "/abilities");
        request.SendWebRequest().completed += (asyncOperation) =>
        {
            if (request.result == UnityWebRequest.Result.Success)
            {
                string jsonResponse = request.downloadHandler.text;

                try
                {
                    // Deserialize JSON and extract abilities
                    var abilitiesDictionary = JsonConvert.DeserializeObject<Dictionary<string, Configurator.AbilityDTO>>(
                        JsonConvert.DeserializeObject<dynamic>(jsonResponse).abilities.ToString()
                    );

                    abilitiesList = new List<Configurator.AbilityDTO>(abilitiesDictionary.Values);
                    responseText = string.Empty;
                }
                catch (Exception ex)
                {
                    responseText = "Error deserializing abilities: " + ex.Message;
                }
            }
            else
            {
                responseText = "Error fetching abilities: " + request.error;
            }

            Repaint(); // Update the editor window to reflect changes
        };
    }

    async Task PatchAbility()
    {
        // Create a dictionary to hold only the modified fields
        var modifiedFields = new Dictionary<string, object>();

        // Compare fields and add only modified ones to the dictionary
        if (abilityName != selectedAbility.name)
        {
            modifiedFields["name"] = abilityName;
        }
        if (abilityCooldown != selectedAbility.cooldown)
        {
            modifiedFields["cooldown"] = abilityCooldown;
        }
        if (abilityRange != selectedAbility.range)
        {
            modifiedFields["range"] = abilityRange;
        }

        // If no fields were modified, exit early
        if (modifiedFields.Count == 0)
        {
            Debug.Log("No changes detected. Skipping PATCH request.");
            return;
        }

        // Serialize the modified fields to JSON
        string jsonPayload = JsonConvert.SerializeObject(modifiedFields);

        using (UnityWebRequest request = new UnityWebRequest($"{BACKEND_URL}/ability/{selectedAbility.id}", "PATCH"))
        {
            // Attach JSON payload
            byte[] bodyRaw = Encoding.UTF8.GetBytes(jsonPayload);
            request.uploadHandler = new UploadHandlerRaw(bodyRaw);

            // Download handler to handle response
            request.downloadHandler = new DownloadHandlerBuffer();

            // Set headers
            request.SetRequestHeader("Content-Type", "application/json");

            // Send the request and await response
            await SendWebRequestAsync(request);

            // Handle response
            if (request.result == UnityWebRequest.Result.Success)
            {
                Debug.Log($"PATCH successful: {request.downloadHandler.text}");
                selectedAbility = JsonConvert.DeserializeObject<Configurator.AbilityDTO>(request.downloadHandler.text);

            }
            else
            {
                Debug.LogError($"PATCH failed: {request.error}");
            }
        }
    }

    async Task SendWebRequestAsync(UnityWebRequest request)
    {
        var operation = request.SendWebRequest();

        while (!operation.isDone)
        {
            await Task.Yield(); // Yield execution to avoid blocking the main thread
        }

        if (request.result != UnityWebRequest.Result.Success)
        {
            Debug.LogError($"HTTP Error: {request.responseCode}, Message: {request.error}");
        }
    }

}
