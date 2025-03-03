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

    string responseText = string.Empty;
    List<Configurator.AbilityDTO> abilitiesList = new List<Configurator.AbilityDTO>();
    Vector2 scrollPosition;
    Configurator.AbilityDTO selectedAbility = null; // Tracks the selected ability

    static readonly string[] TargetingStrategies = { "caster", "character_hit" };
    static readonly string[] MechanicTypes = { "create_projectile", "create_AoE", "delay", "damage" };

    #region Editable fields
    string abilityName; // Very wrong but this acts as a flag for if the editable fields have been initialized
    Sprite abilityIcon;
    int abilityCooldown;
    float abilityRange;
    AbilityParameters abilityTargeting;
    CharacterAction characterState;
    Configurator.AbilityDTO.MechanicDTO[] abilityMechanics;
    #endregion

    bool isEditing = false;
    Configurator.AbilityDTO abilityToDelete = null;

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

        // Handle deletion outside of OnGUI
        if (abilityToDelete != null)
        {
            DeleteAbility(abilityToDelete);
            abilityToDelete = null; // Reset after handling
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
                EditorGUILayout.BeginHorizontal();
                if (GUILayout.Button(ability.name, GUILayout.ExpandWidth(true), GUILayout.Height(25)))
                {
                    isEditing = true;
                    selectedAbility = ability; // Set selected ability to show detail view
                }
                if (GUILayout.Button("Delete", GUILayout.Width(50), GUILayout.Height(25)))
                {
                    abilityToDelete = ability; // Mark for deletion outside OnGUI
                }
                EditorGUILayout.EndHorizontal();
            }

            EditorGUILayout.BeginHorizontal();
            GUILayout.FlexibleSpace();
            if (GUILayout.Button("+", GUILayout.Width(30), GUILayout.Height(30)))
            {
                isEditing = false;
                selectedAbility = new() { mechanics = new Configurator.AbilityDTO.MechanicDTO[0] };
            }
            EditorGUILayout.EndHorizontal();

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
        abilityIcon = EditorGUILayout.ObjectField(abilityIcon, typeof(Sprite), false, GUILayout.Width(64f), GUILayout.Height(64f)) as Sprite;
        abilityRange = EditorGUILayout.FloatField("Range", abilityRange);
        abilityCooldown = EditorGUILayout.IntField("Cooldown", abilityCooldown);
        abilityTargeting = (AbilityParameters)EditorGUILayout.Popup("Targeting", (int)abilityTargeting, Enum.GetNames(typeof(AbilityParameters)));
        characterState = (CharacterAction)EditorGUILayout.Popup("Character State", (int)characterState, Enum.GetNames(typeof(CharacterAction)));

        GUILayout.Label("Mechanics:", EditorStyles.boldLabel);
        abilityMechanics = DisplayMechanics(abilityMechanics);

        if (abilityIcon == null)
        {
            GUILayout.Label("An icon is required", EditorStyles.boldLabel);
            return;
        }

        if (isEditing)
        {
            if (GUILayout.Button("Save Changes"))
            {
                if (await PatchAbility())
                {
                    UpdateAbilityScriptableObject(selectedAbility);
                }
                UpdateAbilityIcon(selectedAbility.id);
            }
        }
        else
        {
            if (GUILayout.Button("Create Ability"))
            {
                if (await CreateAbility())
                {
                    CreateAbilityScriptableObject(selectedAbility);
                }
            }
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
                var damageParams = mechanic.@params as Configurator.AbilityDTO.DamageParams;
                damageParams.amount = EditorGUILayout.IntField("Amount", damageParams.amount);
                int targetingStrategySelectedIndex = Mathf.Max(0, Array.IndexOf(TargetingStrategies, damageParams.targetingStrategy));
                int newTargetingStrategySelectedIndex = EditorGUILayout.Popup("Targeting Strategy", targetingStrategySelectedIndex, TargetingStrategies);
                damageParams.targetingStrategy = TargetingStrategies[newTargetingStrategySelectedIndex];
                GUILayout.Label("On Hit Mechanics:", EditorStyles.boldLabel);
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
        characterState = selectedAbility.characterAction;

        if (isEditing)
        {
            Ability abilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(selectedAbility.id, "Assets/ScriptableObjects/Abilities");
            abilityIcon = abilitySO.icon;
        }

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

    #region Create Ability

    async Task<bool> CreateAbility()
    {
        // Create a dictionary to hold only the modified fields
        var abilityFields = new Dictionary<string, object>
        {
            ["name"] = abilityName,
            ["cooldown"] = abilityCooldown,
            ["range"] = abilityRange,
            ["targeting"] = abilityTargeting,
            ["character_state"] = characterState,
            ["mechanics"] = abilityMechanics
        };

        // Serialize the modified fields to JSON
        string jsonPayload = JsonConvert.SerializeObject(abilityFields);

        using (UnityWebRequest request = new UnityWebRequest($"{BACKEND_URL}/ability/{selectedAbility.id}", "POST"))
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
                Debug.Log($"POST successful: {request.downloadHandler.text}");
                selectedAbility = JsonConvert.DeserializeObject<Configurator.AbilityDTO>(request.downloadHandler.text);
                return true;
            }
            else
            {
                Debug.LogError($"POST failed: {request.error}");
            }

            return false;
        }
    }

    void CreateAbilityScriptableObject(Configurator.AbilityDTO selectedAbility)
    {
        Ability abilitySO = null;

        switch (selectedAbility.targeting)
        {
            case AbilityParameters.TargetId:
                abilitySO = CreateInstance<TargetedAbility>();
                abilitySO.targetLayer = LayerMask.GetMask("Players");
                break;
            case AbilityParameters.TargetPosition:
                abilitySO = CreateInstance<DirectionalAbility>();
                abilitySO.targetLayer = LayerMask.GetMask("Ground");
                break;
            default:
                Debug.LogWarning("Unexpected targeting type.");
                break;

        }

        abilitySO.id = selectedAbility.id;
        abilitySO.name = selectedAbility.name;
        abilitySO.icon = abilityIcon;

        string name = AssetDatabase.GenerateUniqueAssetPath($"Assets/ScriptableObjects/Abilities/{abilitySO.name}.asset");
        AssetDatabase.CreateAsset(abilitySO, name);
        AssetDatabase.SaveAssets();
        AddAbilitiesReferences(abilitySO);
        EditorUtility.FocusProjectWindow();
    }

    #endregion

    #region Modify ability

    async Task<bool> PatchAbility()
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
        if (characterState != selectedAbility.characterAction)
            modifiedFields["character_state"] = characterState;

        string serializedAbilityMechanics = JsonConvert.SerializeObject(abilityMechanics);
        string serializedOriginalMechanics = JsonConvert.SerializeObject(selectedAbility.mechanics);
        if (serializedAbilityMechanics != serializedOriginalMechanics)
            modifiedFields["mechanics"] = abilityMechanics;

        // If no fields were modified, exit early
        if (modifiedFields.Count == 0)
        {
            Debug.Log("No changes detected. Skipping PATCH request.");
            return false;
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
                return true;
            }
            else
            {
                Debug.LogError($"PATCH failed: {request.error}");
            }

            return false;
        }
    }

    void UpdateAbilityScriptableObject(Configurator.AbilityDTO selectedAbility)
    {
        Ability oldAbilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(selectedAbility.id, "Assets/ScriptableObjects/Abilities");
        if (oldAbilitySO == null)
        {
            Debug.LogError($"Ability scriptable object not found for ability id: {selectedAbility.id}");
            return;
        }

        Ability abilitySO = null;

        switch (selectedAbility.targeting)
        {
            case AbilityParameters.TargetId:
                abilitySO = CreateInstance<TargetedAbility>();
                abilitySO.targetLayer = LayerMask.GetMask("Players");
                break;
            case AbilityParameters.TargetPosition:
                abilitySO = CreateInstance<DirectionalAbility>();
                abilitySO.targetLayer = LayerMask.GetMask("Ground");
                break;
            default:
                Debug.LogWarning("Unexpected targeting type.");
                break;

        }

        abilitySO.id = selectedAbility.id;
        abilitySO.name = selectedAbility.name;
        abilitySO.icon = abilityIcon;

        string name = AssetDatabase.GenerateUniqueAssetPath($"Assets/ScriptableObjects/Abilities/{abilitySO.name}.asset");
        AssetDatabase.CreateAsset(abilitySO, name);
        AssetDatabase.SaveAssets();

        ReplaceAbilitiesReferences(abilitySO, oldAbilitySO);

        string oldAbilityPath = AssetDatabase.GetAssetPath(oldAbilitySO);
        AssetDatabase.DeleteAsset(oldAbilityPath);

        EditorUtility.FocusProjectWindow();
    }

    void UpdateAbilityIcon(string abilityId)
    {
        Ability abilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(abilityId, "Assets/ScriptableObjects/Abilities");
        if (abilitySO == null)
        {
            Debug.LogError($"Ability scriptable object not found for ability id: {abilityId}");
            return;
        }
        abilitySO.icon = abilityIcon;
        EditorUtility.SetDirty(abilitySO);
        AssetDatabase.SaveAssets();
        EditorUtility.FocusProjectWindow();
    }

    #endregion

    #region Delete Ability

    async Task<bool> DeleteAbility(Configurator.AbilityDTO ability)
    {
        using (UnityWebRequest request = new UnityWebRequest($"{BACKEND_URL}/ability/{ability.id}", "DELETE"))
        {
            // Download handler to handle response
            request.downloadHandler = new DownloadHandlerBuffer();

            // Set headers
            request.SetRequestHeader("Content-Type", "application/json");

            // Send the request and await response
            await SendWebRequestAsync(request);

            // Handle response
            if (request.result == UnityWebRequest.Result.Success)
            {
                Debug.Log($"DELETE successful: {request.downloadHandler.text}");
                abilitiesList.Remove(ability);
                DeleteAbilityScriptableObject(ability.id);
                return true;
            }
            else
            {
                Debug.LogError($"DELETE failed: {request.error}");
            }

            return false;
        }
    }

    void DeleteAbilityScriptableObject(string abilityId)
    {
        Ability abilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(abilityId, "Assets/ScriptableObjects/Abilities");
        if (abilitySO == null)
        {
            Debug.LogError($"Ability scriptable object not found for ability id: {abilityId}");
            return;
        }

        RemoveAbilitiesReferences(abilitySO);

        string oldAbilityPath = AssetDatabase.GetAssetPath(abilitySO);
        AssetDatabase.DeleteAsset(oldAbilityPath);
    }

    #endregion

    void AddAbilitiesReferences(Ability abilitySO)
    {
        string prefabPath = $"Assets/Prefabs/AbilitiesContainer.prefab";
        GameObject prefab = AssetDatabase.LoadAssetAtPath<GameObject>(prefabPath);
        if (prefab == null)
        {
            Debug.LogError($"Prefab not found at path: {prefabPath}");
            return;
        }
        AbilitiesContainer abilitiesContainer = prefab.GetComponent<AbilitiesContainer>();
        if (abilitiesContainer != null)
        {
            abilitiesContainer.availableAbilities.Add(abilitySO);
            EditorUtility.SetDirty(prefab);
            AssetDatabase.SaveAssets();
        }
        else
        {
            Debug.LogWarning($"Component 'AbilitiesContainer' not found on prefab at {prefabPath}");
        }
    }

    void ReplaceAbilitiesReferences(Ability abilitySO, Ability oldAbilitySO)
    {
        string prefabPath = $"Assets/Prefabs/AbilitiesContainer.prefab";
        GameObject prefab = AssetDatabase.LoadAssetAtPath<GameObject>(prefabPath);
        if (prefab == null)
        {
            Debug.LogError($"Prefab not found at path: {prefabPath}");
            return;
        }
        AbilitiesContainer abilitiesContainer = prefab.GetComponent<AbilitiesContainer>();
        if (abilitiesContainer != null)
        {
            abilitiesContainer.availableAbilities.Remove(oldAbilitySO);
            abilitiesContainer.availableAbilities.Add(abilitySO);
            EditorUtility.SetDirty(prefab);
            AssetDatabase.SaveAssets();
        }
        else
        {
            Debug.LogWarning($"Component 'AbilitiesContainer' not found on prefab at {prefabPath}");
        }
    }

    void RemoveAbilitiesReferences(Ability abilitySO)
    {
        string prefabPath = $"Assets/Prefabs/AbilitiesContainer.prefab";
        GameObject prefab = AssetDatabase.LoadAssetAtPath<GameObject>(prefabPath);
        if (prefab == null)
        {
            Debug.LogError($"Prefab not found at path: {prefabPath}");
            return;
        }
        AbilitiesContainer abilitiesContainer = prefab.GetComponent<AbilitiesContainer>();
        if (abilitiesContainer != null)
        {
            abilitiesContainer.availableAbilities.Remove(abilitySO);
            EditorUtility.SetDirty(prefab);
            AssetDatabase.SaveAssets();
        }
        else
        {
            Debug.LogWarning($"Component 'AbilitiesContainer' not found on prefab at {prefabPath}");
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
