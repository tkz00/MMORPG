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
    const string ABILITIES_FOLDER = "Assets/ScriptableObjects/Abilities";
    const string ABILITIES_CONTAINER_PREFAB = "Assets/Prefabs/AbilitiesContainer.prefab";

    readonly AbilityHttpClient httpClient;
    readonly AbilityUIManager uiManager;
    readonly AbilityMechanicManager mechanicManager;
    readonly AbilityScriptableObjectManager soManager;

    string responseText = string.Empty;
    List<Configurator.AbilityDTO> abilitiesList = new();
    Vector2 scrollPosition;
    Configurator.AbilityDTO selectedAbility;
    bool isEditing;
    bool isSaving;
    bool isCreating;
    Configurator.AbilityDTO abilityToDelete;

    public AbilitiesEditorWindow()
    {
        httpClient = new AbilityHttpClient(BACKEND_URL);
        soManager = new AbilityScriptableObjectManager(ABILITIES_FOLDER, ABILITIES_CONTAINER_PREFAB);
        uiManager = new AbilityUIManager(soManager);
        mechanicManager = new AbilityMechanicManager();
    }

    [MenuItem("Window/Abilities Editor Window")]
    public static void ShowWindow()
    {
        GetWindow<AbilitiesEditorWindow>("Abilities Editor Window");
    }

    void OnGUI()
    {
        GUILayout.Space(10);
        GUILayout.Label("Abilities Editor", EditorStyles.boldLabel);

        if (selectedAbility == null)
        {
            DrawListView();
        }
        else
        {
            DrawDetailView();
        }

        HandlePendingOperations();
    }

    void HandlePendingOperations()
    {
        if (abilityToDelete != null)
        {
            DeleteAbility(abilityToDelete);
            abilityToDelete = null;
        }

        if (isSaving)
        {
            PatchAbility();
            isSaving = false;
        }

        if (isCreating)
        {
            CreateAbility();
            isCreating = false;
        }
    }

    void DrawListView()
    {
        if (GUILayout.Button("Fetch Abilities", GUILayout.Height(30)))
        {
            FetchAbilities();
        }

        uiManager.DrawResponseText(responseText);

        GUILayout.Space(10);
        GUILayout.Label("Abilities List:", EditorStyles.boldLabel);

        if (abilitiesList.Count == 0)
        {
            GUILayout.Label("No abilities loaded. Fetch abilities to see the list.");
            return;
        }

        DrawAbilitiesList();
    }

    void DrawAbilitiesList()
    {
        scrollPosition = GUILayout.BeginScrollView(scrollPosition, GUILayout.Height(300));

        foreach (var ability in abilitiesList)
        {
            DrawAbilityListItem(ability);
        }

        DrawCreateAbilityButton();

        GUILayout.EndScrollView();
    }

    void DrawAbilityListItem(Configurator.AbilityDTO ability)
    {
        EditorGUILayout.BeginHorizontal();
        if (GUILayout.Button(ability.name, GUILayout.ExpandWidth(true), GUILayout.Height(25)))
        {
            isEditing = true;
            selectedAbility = ability;
        }
        if (GUILayout.Button("Delete", GUILayout.Width(50), GUILayout.Height(25)))
        {
            abilityToDelete = ability;
        }
        EditorGUILayout.EndHorizontal();
    }

    void DrawCreateAbilityButton()
    {
        EditorGUILayout.BeginHorizontal();
        GUILayout.FlexibleSpace();
        if (GUILayout.Button("+", GUILayout.Width(30), GUILayout.Height(30)))
        {
            isEditing = false;
            selectedAbility = new() { mechanics = new Configurator.AbilityDTO.MechanicDTO[0] };
        }
        EditorGUILayout.EndHorizontal();
    }

    void DrawDetailView()
    {
        if (!uiManager.IsInitialized)
        {
            uiManager.InitializeFields(selectedAbility, isEditing);
        }

        if (GUILayout.Button("Back to List", GUILayout.Height(30)))
        {
            selectedAbility = null;
            uiManager.Reset();
            FetchAbilities();
            return;
        }

        uiManager.DrawAbilityDetails(selectedAbility);

        if (uiManager.pendingIcon == null)
        {
            GUILayout.Label("An icon is required", EditorStyles.boldLabel);
            return;
        }

        DrawSaveOrCreateButton();
    }

    void DrawSaveOrCreateButton()
    {
        if (isEditing)
        {
            if (GUILayout.Button("Save Changes"))
            {
                isSaving = true;
            }
        }
        else
        {
            if (GUILayout.Button("Create Ability"))
            {
                isCreating = true;
            }
        }
    }

    async void FetchAbilities()
    {
        try
        {
            var response = await httpClient.GetAbilities();
            abilitiesList = new List<Configurator.AbilityDTO>(response.Values);
            responseText = string.Empty;
        }
        catch (Exception ex)
        {
            responseText = $"Error fetching abilities: {ex.Message}";
        }

        Repaint();
    }

    async void CreateAbility()
    {
        try
        {
            var abilityFields = uiManager.GetAbilityFields();
            var response = await httpClient.CreateAbility(selectedAbility.id, abilityFields);
            selectedAbility = response;
            soManager.CreateAbilityScriptableObject(selectedAbility, uiManager.pendingIcon);
        }
        catch (Exception ex)
        {
            Debug.LogError($"Failed to create ability: {ex.Message}");
        }
    }

    async void PatchAbility()
    {
        try
        {
            var modifiedFields = uiManager.GetModifiedFields(selectedAbility);
            if (modifiedFields.Count == 0 && !uiManager.HasIconChanges())
            {
                Debug.Log("No changes detected. Skipping PATCH request.");
                return;
            }

            var response = await httpClient.PatchAbility(selectedAbility.id, modifiedFields);
            selectedAbility = response;
            uiManager.ApplyChanges(selectedAbility);
            soManager.UpdateAbilityScriptableObject(selectedAbility, uiManager.Icon);
        }
        catch (Exception ex)
        {
            Debug.LogError($"Failed to patch ability: {ex.Message}");
        }
    }

    async void DeleteAbility(Configurator.AbilityDTO ability)
    {
        try
        {
            await httpClient.DeleteAbility(ability.id);
            abilitiesList.Remove(ability);
            soManager.DeleteAbilityScriptableObject(ability.id);
        }
        catch (Exception ex)
        {
            Debug.LogError($"Failed to delete ability: {ex.Message}");
        }
    }
}

public class AbilityHttpClient
{
    readonly string baseUrl;

    public AbilityHttpClient(string baseUrl)
    {
        this.baseUrl = baseUrl;
    }

    public async Task<Dictionary<string, Configurator.AbilityDTO>> GetAbilities()
    {
        using var request = UnityWebRequest.Get($"{baseUrl}/abilities");
        await SendWebRequestAsync(request);

        if (request.result != UnityWebRequest.Result.Success)
            throw new Exception(request.error);

        var jsonResponse = request.downloadHandler.text;
        var response = JsonConvert.DeserializeObject<dynamic>(jsonResponse);
        return JsonConvert.DeserializeObject<Dictionary<string, Configurator.AbilityDTO>>(response.abilities.ToString());
    }

    public async Task<Configurator.AbilityDTO> CreateAbility(string id, Dictionary<string, object> fields)
    {
        using var request = CreateWebRequest($"{baseUrl}/ability/{id}", "POST", fields);
        await SendWebRequestAsync(request);

        if (request.result != UnityWebRequest.Result.Success)
            throw new Exception(request.error);

        return JsonConvert.DeserializeObject<Configurator.AbilityDTO>(request.downloadHandler.text);
    }

    public async Task<Configurator.AbilityDTO> PatchAbility(string id, Dictionary<string, object> fields)
    {
        using var request = CreateWebRequest($"{baseUrl}/ability/{id}", "PATCH", fields);
        await SendWebRequestAsync(request);

        if (request.result != UnityWebRequest.Result.Success)
            throw new Exception(request.error);

        return JsonConvert.DeserializeObject<Configurator.AbilityDTO>(request.downloadHandler.text);
    }

    public async Task DeleteAbility(string id)
    {
        using var request = new UnityWebRequest($"{baseUrl}/ability/{id}", "DELETE");
        request.downloadHandler = new DownloadHandlerBuffer();
        request.SetRequestHeader("Content-Type", "application/json");

        await SendWebRequestAsync(request);

        if (request.result != UnityWebRequest.Result.Success)
            throw new Exception(request.error);
    }

    UnityWebRequest CreateWebRequest(string url, string method, Dictionary<string, object> fields)
    {
        var request = new UnityWebRequest(url, method);
        var jsonPayload = JsonConvert.SerializeObject(fields);
        var bodyRaw = Encoding.UTF8.GetBytes(jsonPayload);
        request.uploadHandler = new UploadHandlerRaw(bodyRaw);
        request.downloadHandler = new DownloadHandlerBuffer();
        request.SetRequestHeader("Content-Type", "application/json");
        return request;
    }

    async Task SendWebRequestAsync(UnityWebRequest request)
    {
        var operation = request.SendWebRequest();
        while (!operation.isDone)
        {
            await Task.Yield();
        }
    }
}

public class AbilityUIManager
{
    readonly AbilityScriptableObjectManager soManager;
    bool isEditing;

    public bool IsInitialized { get; set; }
    public Sprite Icon { get; set; }
    public Sprite pendingIcon;

    string abilityName;
    int abilityCooldown;
    float abilityRange;
    AbilityParameters abilityTargeting;
    CharacterAction characterState;
    Configurator.AbilityDTO.MechanicDTO[] abilityMechanics;

    public AbilityUIManager(AbilityScriptableObjectManager soManager)
    {
        this.soManager = soManager;
    }

    public void InitializeFields(Configurator.AbilityDTO ability, bool isEditing)
    {
        this.isEditing = isEditing;
        abilityName = ability.name;
        abilityCooldown = ability.cooldown;
        abilityRange = ability.range;
        abilityTargeting = ability.targeting;
        characterState = ability.characterAction;

        if (isEditing)
        {
            var abilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(ability.id, "Assets/ScriptableObjects/Abilities");
            Icon = abilitySO.icon;
            pendingIcon = Icon;
        }

        abilityMechanics = new Configurator.AbilityDTO.MechanicDTO[ability.mechanics.Length];
        for (int i = 0; i < ability.mechanics.Length; i++)
        {
            abilityMechanics[i] = ability.mechanics[i].DeepCopy();
        }

        IsInitialized = true;
    }

    public void Reset()
    {
        IsInitialized = false;
        abilityName = null;
        Icon = null;
        pendingIcon = null;
        isEditing = false;
    }

    public void DrawResponseText(string responseText)
    {
        if (string.IsNullOrEmpty(responseText)) return;

        GUILayout.Space(10);
        GUILayout.Label("Response:", EditorStyles.boldLabel);
        GUILayout.TextArea(responseText, GUILayout.Height(50));
    }

    public void DrawAbilityDetails(Configurator.AbilityDTO ability)
    {
        GUILayout.Space(10);
        GUILayout.Label($"Details for: {ability.name}", EditorStyles.boldLabel);
        GUILayout.Label($"ID: {ability.id}");

        abilityName = EditorGUILayout.TextField("Name", abilityName);
        pendingIcon = EditorGUILayout.ObjectField(pendingIcon, typeof(Sprite), false, GUILayout.Width(64f), GUILayout.Height(64f)) as Sprite;
        abilityRange = EditorGUILayout.FloatField("Range", abilityRange);
        abilityCooldown = EditorGUILayout.IntField("Cooldown", abilityCooldown);
        abilityTargeting = (AbilityParameters)EditorGUILayout.Popup("Targeting", (int)abilityTargeting, Enum.GetNames(typeof(AbilityParameters)));
        characterState = (CharacterAction)EditorGUILayout.Popup("Character State", (int)characterState, Enum.GetNames(typeof(CharacterAction)));

        GUILayout.Label("Mechanics:", EditorStyles.boldLabel);
        abilityMechanics = AbilityMechanicManager.DisplayMechanics(abilityMechanics);
    }

    public Dictionary<string, object> GetAbilityFields()
    {
        return new Dictionary<string, object>
        {
            ["name"] = abilityName,
            ["cooldown"] = abilityCooldown,
            ["range"] = abilityRange,
            ["targeting"] = abilityTargeting,
            ["character_state"] = characterState,
            ["mechanics"] = abilityMechanics
        };
    }

    public Dictionary<string, object> GetModifiedFields(Configurator.AbilityDTO originalAbility)
    {
        var modifiedFields = new Dictionary<string, object>();

        if (abilityName != originalAbility.name)
            modifiedFields["name"] = abilityName;
        if (abilityCooldown != originalAbility.cooldown)
            modifiedFields["cooldown"] = abilityCooldown;
        if (abilityRange != originalAbility.range)
            modifiedFields["range"] = abilityRange;
        if (abilityTargeting != originalAbility.targeting)
            modifiedFields["targeting"] = abilityTargeting;
        if (characterState != originalAbility.characterAction)
            modifiedFields["character_state"] = characterState;

        var serializedAbilityMechanics = JsonConvert.SerializeObject(abilityMechanics);
        var serializedOriginalMechanics = JsonConvert.SerializeObject(originalAbility.mechanics);
        if (serializedAbilityMechanics != serializedOriginalMechanics)
            modifiedFields["mechanics"] = abilityMechanics;

        return modifiedFields;
    }

    public bool HasIconChanges()
    {
        return pendingIcon != Icon;
    }

    public void ApplyChanges(Configurator.AbilityDTO ability)
    {
        if (isEditing && HasIconChanges())
        {
            Icon = pendingIcon;
            soManager.UpdateAbilityIcon(ability.id, Icon);
        }
    }
}

public class AbilityMechanicManager
{
    static readonly string[] TargetingStrategies = { "caster", "character_hit" };
    static readonly string[] MechanicTypes = { "create_projectile", "create_AoE", "delay", "damage" };

    public static Configurator.AbilityDTO.MechanicDTO[] DisplayMechanics(Configurator.AbilityDTO.MechanicDTO[] mechanics)
    {
        GUILayout.Space(10);

        var mechanicsList = new List<Configurator.AbilityDTO.MechanicDTO>(mechanics);
        for (int i = mechanicsList.Count - 1; i >= 0; i--)
        {
            var mechanic = DisplayMechanic(mechanicsList[i], i, mechanicsList);
            if (mechanic == null)
            {
                mechanicsList.RemoveAt(i);
            }
        }

        if (GUILayout.Button("+", GUILayout.Width(20), GUILayout.Height(20)))
        {
            mechanicsList.Add(new Configurator.AbilityDTO.MechanicDTO
            {
                mechanicType = "create_projectile",
                @params = new Configurator.AbilityDTO.CreateProjectileParams()
            });
        }

        return mechanicsList.ToArray();
    }

    static Configurator.AbilityDTO.MechanicDTO DisplayMechanic(Configurator.AbilityDTO.MechanicDTO mechanic, int index, List<Configurator.AbilityDTO.MechanicDTO> mechanics)
    {
        EditorGUILayout.BeginVertical("box");

        EditorGUILayout.BeginHorizontal();
        int selectedIndex = Mathf.Max(0, Array.IndexOf(MechanicTypes, mechanic.mechanicType));
        int newSelectedIndex = EditorGUILayout.Popup("Mechanic Type", selectedIndex, MechanicTypes);

        if (GUILayout.Button("-", GUILayout.Width(20), GUILayout.Height(20)))
        {
            EditorGUILayout.EndHorizontal();
            EditorGUILayout.EndVertical();
            return null;
        }

        if (newSelectedIndex != selectedIndex)
        {
            mechanic.mechanicType = MechanicTypes[newSelectedIndex];
            mechanic.@params = CreateMechanicParams(mechanic.mechanicType);
        }

        EditorGUILayout.EndHorizontal();

        DisplayMechanicParams(mechanic);

        EditorGUILayout.EndVertical();

        return mechanic;
    }

    static Configurator.AbilityDTO.Params CreateMechanicParams(string mechanicType)
    {
        return mechanicType switch
        {
            "create_projectile" => (Configurator.AbilityDTO.Params)new Configurator.AbilityDTO.CreateProjectileParams(),
            "create_AoE" => (Configurator.AbilityDTO.Params)new Configurator.AbilityDTO.CreateAoEParams(),
            "damage" => (Configurator.AbilityDTO.Params)new Configurator.AbilityDTO.DamageParams(),
            "delay" => (Configurator.AbilityDTO.Params)new Configurator.AbilityDTO.DelayParams(),
            _ => throw new ArgumentException($"Unknown mechanic type: {mechanicType}")
        };
    }

    static void DisplayMechanicParams(Configurator.AbilityDTO.MechanicDTO mechanic)
    {
        switch (mechanic.mechanicType)
        {
            case "create_projectile":
                DisplayProjectileParams(mechanic);
                break;
            case "create_AoE":
                DisplayAoEParams(mechanic);
                break;
            case "damage":
                DisplayDamageParams(mechanic);
                break;
            case "delay":
                DisplayDelayParams(mechanic);
                break;
            default:
                GUILayout.Label($"Error, mechanic {mechanic.mechanicType} not found");
                break;
        }
    }

    static void DisplayProjectileParams(Configurator.AbilityDTO.MechanicDTO mechanic)
    {
        GUILayout.Label("On Hit Mechanics:", EditorStyles.boldLabel);
        var projectileParams = mechanic.@params as Configurator.AbilityDTO.CreateProjectileParams;
        projectileParams.onHitMechanics = DisplayMechanics(projectileParams.onHitMechanics);
        mechanic.@params = projectileParams;
    }

    static void DisplayAoEParams(Configurator.AbilityDTO.MechanicDTO mechanic)
    {
        var aoEParams = mechanic.@params as Configurator.AbilityDTO.CreateAoEParams;
        aoEParams.durationMs = EditorGUILayout.IntField("Duration (ms)", aoEParams.durationMs);
        aoEParams.radius = EditorGUILayout.FloatField("Radius", aoEParams.radius);

        GUILayout.Label("On Hit Mechanics:", EditorStyles.boldLabel);
        aoEParams.onHitMechanics = DisplayMechanics(aoEParams.onHitMechanics);
        mechanic.@params = aoEParams;
    }

    static void DisplayDamageParams(Configurator.AbilityDTO.MechanicDTO mechanic)
    {
        var damageParams = mechanic.@params as Configurator.AbilityDTO.DamageParams;
        damageParams.amount = EditorGUILayout.IntField("Amount", damageParams.amount);

        int targetingStrategySelectedIndex = Mathf.Max(0, Array.IndexOf(TargetingStrategies, damageParams.targetingStrategy));
        int newTargetingStrategySelectedIndex = EditorGUILayout.Popup("Targeting Strategy", targetingStrategySelectedIndex, TargetingStrategies);
        damageParams.targetingStrategy = TargetingStrategies[newTargetingStrategySelectedIndex];

        GUILayout.Label("On Hit Mechanics:", EditorStyles.boldLabel);
        damageParams.onHitMechanics = DisplayMechanics(damageParams.onHitMechanics);
        mechanic.@params = damageParams;
    }

    static void DisplayDelayParams(Configurator.AbilityDTO.MechanicDTO mechanic)
    {
        var delayParams = mechanic.@params as Configurator.AbilityDTO.DelayParams;
        delayParams.delayMs = EditorGUILayout.IntField("Delay", delayParams.delayMs);
        delayParams.executeAfterDelayMechanics = DisplayMechanics(delayParams.executeAfterDelayMechanics);
        mechanic.@params = delayParams;
    }
}

public class AbilityScriptableObjectManager
{
    readonly string abilitiesFolder;
    readonly string abilitiesContainerPrefab;

    public AbilityScriptableObjectManager(string abilitiesFolder, string abilitiesContainerPrefab)
    {
        this.abilitiesFolder = abilitiesFolder;
        this.abilitiesContainerPrefab = abilitiesContainerPrefab;
    }

    public void CreateAbilityScriptableObject(Configurator.AbilityDTO ability, Sprite icon)
    {
        var abilitySO = CreateAbilitySO(ability, icon);
        string name = AssetDatabase.GenerateUniqueAssetPath($"{abilitiesFolder}/{abilitySO.name}.asset");
        AssetDatabase.CreateAsset(abilitySO, name);
        AssetDatabase.SaveAssets();
        AddAbilitiesReferences(abilitySO);
        EditorUtility.FocusProjectWindow();
    }

    public void UpdateAbilityScriptableObject(Configurator.AbilityDTO ability, Sprite icon)
    {
        var oldAbilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(ability.id, abilitiesFolder);
        if (oldAbilitySO == null)
        {
            Debug.LogError($"Ability scriptable object not found for ability id: {ability.id}");
            return;
        }

        var abilitySO = CreateAbilitySO(ability, icon);
        string name = AssetDatabase.GenerateUniqueAssetPath($"{abilitiesFolder}/{abilitySO.name}.asset");
        AssetDatabase.CreateAsset(abilitySO, name);
        AssetDatabase.SaveAssets();

        ReplaceAbilitiesReferences(abilitySO, oldAbilitySO);

        string oldAbilityPath = AssetDatabase.GetAssetPath(oldAbilitySO);
        AssetDatabase.DeleteAsset(oldAbilityPath);

        EditorUtility.FocusProjectWindow();
    }

    public void DeleteAbilityScriptableObject(string abilityId)
    {
        var abilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(abilityId, abilitiesFolder);
        if (abilitySO == null)
        {
            Debug.LogError($"Ability scriptable object not found for ability id: {abilityId}");
            return;
        }

        RemoveAbilitiesReferences(abilitySO);

        string oldAbilityPath = AssetDatabase.GetAssetPath(abilitySO);
        AssetDatabase.DeleteAsset(oldAbilityPath);
    }

    Ability CreateAbilitySO(Configurator.AbilityDTO ability, Sprite icon)
    {
        Ability abilitySO = (ability.targeting switch
        {
            AbilityParameters.TargetId => ScriptableObject.CreateInstance<TargetedAbility>(),
            AbilityParameters.TargetPosition => ScriptableObject.CreateInstance<DirectionalAbility>(),
            _ => throw new ArgumentException($"Unexpected targeting type: {ability.targeting}")
        });

        abilitySO.id = ability.id;
        abilitySO.name = ability.name;
        abilitySO.icon = icon;
        abilitySO.targetLayer = ability.targeting == AbilityParameters.TargetId
            ? LayerMask.GetMask("Players")
            : LayerMask.GetMask("Ground");

        return abilitySO;
    }

    void AddAbilitiesReferences(Ability abilitySO)
    {
        var prefab = AssetDatabase.LoadAssetAtPath<GameObject>(abilitiesContainerPrefab);
        if (prefab == null)
        {
            Debug.LogError($"Prefab not found at path: {abilitiesContainerPrefab}");
            return;
        }

        var abilitiesContainer = prefab.GetComponent<AbilitiesContainer>();
        if (abilitiesContainer != null)
        {
            abilitiesContainer.availableAbilities.Add(abilitySO);
            EditorUtility.SetDirty(prefab);
            AssetDatabase.SaveAssets();
        }
        else
        {
            Debug.LogWarning($"Component 'AbilitiesContainer' not found on prefab at {abilitiesContainerPrefab}");
        }
    }

    void ReplaceAbilitiesReferences(Ability newAbilitySO, Ability oldAbilitySO)
    {
        var prefab = AssetDatabase.LoadAssetAtPath<GameObject>(abilitiesContainerPrefab);
        if (prefab == null)
        {
            Debug.LogError($"Prefab not found at path: {abilitiesContainerPrefab}");
            return;
        }

        var abilitiesContainer = prefab.GetComponent<AbilitiesContainer>();
        if (abilitiesContainer != null)
        {
            abilitiesContainer.availableAbilities.Remove(oldAbilitySO);
            abilitiesContainer.availableAbilities.Add(newAbilitySO);
            EditorUtility.SetDirty(prefab);
            AssetDatabase.SaveAssets();
        }
        else
        {
            Debug.LogWarning($"Component 'AbilitiesContainer' not found on prefab at {abilitiesContainerPrefab}");
        }
    }

    void RemoveAbilitiesReferences(Ability abilitySO)
    {
        var prefab = AssetDatabase.LoadAssetAtPath<GameObject>(abilitiesContainerPrefab);
        if (prefab == null)
        {
            Debug.LogError($"Prefab not found at path: {abilitiesContainerPrefab}");
            return;
        }

        var abilitiesContainer = prefab.GetComponent<AbilitiesContainer>();
        if (abilitiesContainer != null)
        {
            abilitiesContainer.availableAbilities.Remove(abilitySO);
            EditorUtility.SetDirty(prefab);
            AssetDatabase.SaveAssets();
        }
        else
        {
            Debug.LogWarning($"Component 'AbilitiesContainer' not found on prefab at {abilitiesContainerPrefab}");
        }
    }

    public void UpdateAbilityIcon(string abilityId, Sprite icon)
    {
        var abilitySO = ScriptableObjectFinder.FindScriptableObjectByID<Ability>(abilityId, abilitiesFolder);
        if (abilitySO == null)
        {
            Debug.LogError($"Ability scriptable object not found for ability id: {abilityId}");
            return;
        }

        abilitySO.icon = icon;
        EditorUtility.SetDirty(abilitySO);
        AssetDatabase.SaveAssets();
        EditorUtility.FocusProjectWindow();
    }
}

public static class ScriptableObjectFinder
{
    public static T FindScriptableObjectByID<T>(string id, string folder) where T : ScriptableObject
    {
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
