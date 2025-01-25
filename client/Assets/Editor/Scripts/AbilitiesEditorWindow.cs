using System;
using System.Collections.Generic;
using Newtonsoft.Json;
using UnityEditor;
using UnityEngine;

public class AbilitiesEditorWindow : EditorWindow
{
    private string responseText = string.Empty;
    private List<Configurator.AbilityDTO> abilitiesList = new List<Configurator.AbilityDTO>();
    private Vector2 scrollPosition;
    private Configurator.AbilityDTO selectedAbility = null; // Tracks the selected ability

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

    private void DrawDetailView()
    {
        // Back button to return to the list view
        if (GUILayout.Button("Back to List", GUILayout.Height(30)))
        {
            selectedAbility = null; // Clear selected ability to return to the list view
            return;
        }

        GUILayout.Space(10);
        GUILayout.Label($"Details for: {selectedAbility.name}", EditorStyles.boldLabel);

        GUILayout.Label($"ID: {selectedAbility.id}");
        GUILayout.Label($"Name: {selectedAbility.name}");
        GUILayout.Label($"Range: {selectedAbility.range}");
        GUILayout.Label($"Cooldown: {selectedAbility.cooldown}");
        GUILayout.Label($"Targeting: {selectedAbility.targeting}");
        GUILayout.Label($"Character State: {selectedAbility.characterAction}");

        GUILayout.Space(10);
        GUILayout.Label("Mechanics:", EditorStyles.boldLabel);

        // Display mechanics in a scrollable area
        scrollPosition = GUILayout.BeginScrollView(scrollPosition, GUILayout.Height(200));
        foreach (var mechanic in selectedAbility.mechanics)
        {
            GUILayout.Label($"Mechanic Type: {mechanic.mechanicType}");
            // GUILayout.Label($"Params: {JsonConvert.SerializeObject(mechanic.@params, Formatting.Indented)}");

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
    }

    private void FetchAbilities()
    {
        string url = "http://0.0.0.0:8080/abilities";

        var request = UnityEngine.Networking.UnityWebRequest.Get(url);
        request.SendWebRequest().completed += (asyncOperation) =>
        {
            if (request.result == UnityEngine.Networking.UnityWebRequest.Result.Success)
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
}
