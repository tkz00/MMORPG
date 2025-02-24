using System;
using System.Collections.Generic;
using System.Linq;
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

    private static readonly string[] MechanicTypes = { "create_projectile", "create_AoE", "delay", "damage" };

    #region Editable fields
    string abilityName; // Very wrong but this acts as a flag for if the editable fields have been initialized
    int abilityCooldown;
    float abilityRange;
    AbilityParameters abilityTargeting;
    Configurator.AbilityDTO.MechanicDTO[] abilityMechanics;
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

        // reverse next two lines commenting to enable targeting editing
        GUILayout.Label($"Targeting: {selectedAbility.targeting}");
        // abilityTargeting = (AbilityParameters)EditorGUILayout.Popup("Targeting", (int)abilityTargeting, Enum.GetNames(typeof(AbilityParameters)));

        GUILayout.Label($"Character State: {selectedAbility.characterAction}");

        GUILayout.Label("Mechanics:", EditorStyles.boldLabel);
        abilityMechanics = DisplayMechanics(abilityMechanics);

        if (GUILayout.Button("Save Changes"))
        {
            await PatchAbility();
        }
    }

    Configurator.AbilityDTO.MechanicDTO[] DisplayMechanics(Configurator.AbilityDTO.MechanicDTO[] mechanics)
    {
        GUILayout.Space(10);

        for (int i = 0; i < mechanics.Count(); i++)
        {
            Configurator.AbilityDTO.MechanicDTO mechanic = mechanics[i];
            EditorGUILayout.BeginVertical("box");

            EditorGUILayout.BeginHorizontal();
            // Display mechanic type dropdown
            int selectedIndex = Mathf.Max(0, Array.IndexOf(MechanicTypes, mechanics[i].mechanicType));
            int newSelectedIndex = EditorGUILayout.Popup("Mechanic Type", selectedIndex, MechanicTypes);
            if (GUILayout.Button("-", GUILayout.Width(20), GUILayout.Height(20)))
            {
                var auxList = new List<Configurator.AbilityDTO.MechanicDTO>(mechanics);
                auxList.RemoveAt(i);
                mechanics = auxList.ToArray();
            }

            if (newSelectedIndex != selectedIndex)
            {
                mechanics[i].mechanicType = MechanicTypes[newSelectedIndex];
                switch (mechanics[i].mechanicType)
                {
                    case "create_projectile":
                        mechanics[i].@params = new Configurator.AbilityDTO.CreateProjectileParams();
                        break;
                    case "create_AoE":
                        mechanics[i].@params = new Configurator.AbilityDTO.CreateAoEParams();
                        break;
                    case "damage":
                        mechanics[i].@params = new Configurator.AbilityDTO.DamageParams();
                        break;
                    case "delay":
                        mechanics[i].@params = new Configurator.AbilityDTO.DelayParams();
                        break;
                    default:
                        GUILayout.Label($"Error, mechanic {mechanic.mechanicType} not found");
                        break;
                }
            }

            EditorGUILayout.EndHorizontal();

            DisplayMechanic(mechanic);

            EditorGUILayout.EndVertical();
        }

        if (GUILayout.Button("+", GUILayout.Width(20), GUILayout.Height(20)))
        {
            var auxList = new List<Configurator.AbilityDTO.MechanicDTO>(mechanics)
                {
                    new Configurator.AbilityDTO.MechanicDTO{
                        mechanicType = "create_projectile",
                        @params = new Configurator.AbilityDTO.CreateProjectileParams(),
                    }
                };
            mechanics = auxList.ToArray();
        }

        return mechanics;
    }

    void DisplayMechanic(Configurator.AbilityDTO.MechanicDTO mechanic)
    {
        switch (mechanic.mechanicType)
        {
            case "create_projectile":
                GUILayout.Label("On Hit Mechanics:", EditorStyles.boldLabel);
                var projectileParams = mechanic.@params as Configurator.AbilityDTO.CreateProjectileParams;
                projectileParams.onHitMechanics = DisplayMechanics(projectileParams.onHitMechanics);
                mechanic.@params = projectileParams;
                break;
            case "create_AoE":
                GUILayout.Label("On Hit Mechanics:", EditorStyles.boldLabel);
                var aoEParams = mechanic.@params as Configurator.AbilityDTO.CreateAoEParams;
                aoEParams.onHitMechanics = DisplayMechanics(aoEParams.onHitMechanics);
                mechanic.@params = aoEParams;
                break;
            case "damage":
                GUILayout.Label("On Hit Mechanics:", EditorStyles.boldLabel);
                var damageParams = mechanic.@params as Configurator.AbilityDTO.DamageParams;
                damageParams.amount = EditorGUILayout.IntField("Amount", damageParams.amount);
                damageParams.onHitMechanics = DisplayMechanics(damageParams.onHitMechanics);
                mechanic.@params = damageParams;
                break;
            case "delay":
                var delayParams = mechanic.@params as Configurator.AbilityDTO.DelayParams;
                delayParams.delayMs = EditorGUILayout.IntField("Delay", delayParams.delayMs);
                delayParams.executeAfterDelayMechanics = DisplayMechanics(delayParams.executeAfterDelayMechanics);
                mechanic.@params = delayParams;
                break;
            default:
                GUILayout.Label($"Error, mechanic {mechanic.mechanicType} not found");
                break;
        }
    }

    private void InitializeEditableFields()
    {
        abilityName = selectedAbility.name;
        abilityCooldown = selectedAbility.cooldown;
        abilityRange = selectedAbility.range;
        abilityTargeting = selectedAbility.targeting;

        abilityMechanics = new Configurator.AbilityDTO.MechanicDTO[selectedAbility.mechanics.Length];
        for (int i = 0; i < selectedAbility.mechanics.Length; i++)
        {
            abilityMechanics[i] = selectedAbility.mechanics[i].DeepCopy();
        }
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
            modifiedFields["name"] = abilityName;
        if (abilityCooldown != selectedAbility.cooldown)
            modifiedFields["cooldown"] = abilityCooldown;
        if (abilityRange != selectedAbility.range)
            modifiedFields["range"] = abilityRange;
        if (abilityTargeting != selectedAbility.targeting)
            modifiedFields["targeting"] = abilityTargeting;

        string serializedAbilityMechanics = JsonConvert.SerializeObject(abilityMechanics);
        string serializedOriginalMechanics = JsonConvert.SerializeObject(selectedAbility.mechanics);
        if (serializedAbilityMechanics != serializedOriginalMechanics)
            modifiedFields["mechanics"] = abilityMechanics;

        // If no fields were modified, exit early
        if (modifiedFields.Count == 0)
        {
            Debug.Log("No changes detected. Skipping PATCH request.");
            return;
        }

        // Serialize the modified fields to JSON
        string jsonPayload = JsonConvert.SerializeObject(modifiedFields);
        Debug.Log(jsonPayload);

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
                UpdateAbilityScriptableObject(selectedAbility);
            }
            else
            {
                Debug.LogError($"PATCH failed: {request.error}");
            }
        }
    }

    void UpdateAbilityScriptableObject(Configurator.AbilityDTO selectedAbility)
    {
        Ability abilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(selectedAbility.id, "Assets/ScriptableObjects/Abilities");
        if (abilitySO == null)
        {
            Debug.LogError($"Ability scriptable object not found for ability id: {selectedAbility.id}");
            return;
        }

        abilitySO.name = selectedAbility.name;
        switch (selectedAbility.targeting)
        {
            case AbilityParameters.TargetId:
                abilitySO.targetLayer = LayerMask.GetMask("Players");
                break;
            case AbilityParameters.NoTarget:
                abilitySO.targetLayer = ~0;
                break;
            case AbilityParameters.TargetPosition:
                abilitySO.targetLayer = LayerMask.GetMask("Ground");
                break;
        }

        EditorUtility.SetDirty(abilitySO);
        AssetDatabase.SaveAssets();
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

public static class ScriptableObjectFinder
{
    public static T FindScriptableObjectByID<T>(string id, string folder) where T : ScriptableObject
    {
        // Search for all assets of type T
        string[] guids = AssetDatabase.FindAssets($"t:{typeof(T).Name}", new[] { folder });

        foreach (string guid in guids)
        {
            string path = AssetDatabase.GUIDToAssetPath(guid);
            T so = AssetDatabase.LoadAssetAtPath<T>(path);

            // Assuming your ScriptableObject has a public 'ID' field
            var idField = typeof(T).GetField("id");
            if (idField != null)
            {
                string soID = idField.GetValue(so) as string;
                if (soID == id)
                {
                    return so;
                }
            }
        }

        Debug.LogWarning($"ScriptableObject with ID {id} not found.");
        return null;
    }
}
