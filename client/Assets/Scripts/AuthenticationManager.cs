using System;
using System.Threading.Tasks;
using Newtonsoft.Json;
using TMPro;
using UnityEngine;
using UnityEngine.Networking;
using UnityEngine.SceneManagement;
using UnityEngine.UI;

public class AuthenticationManager : MonoBehaviour
{
    [SerializeField] TMP_InputField characterNameInput;
    [SerializeField] Transform quickAccessCharactersContainer;
    [SerializeField] GameObject characterQuickAccessPrefab;

    static string characterName;
    public static string CharacterName { get { return characterName; } }

    void Start()
    {
        PopulateCharacterSelect();
    }

    async Task PopulateCharacterSelect()
    {
        CharactersHttpClient httpClient = new();
        CharacterDTO[] characters = await httpClient.GetCharacters();

        foreach (CharacterDTO character in characters)
        {
            GameObject button = Instantiate(characterQuickAccessPrefab, quickAccessCharactersContainer);
            button.GetComponentInChildren<TMP_Text>().text = character.name;
            button.GetComponent<Button>().onClick.AddListener(() => SelectCharacter(character.name));
        }
    }

    public void SelectCharacter(string characterName)
    {
        AuthenticationManager.characterName = characterName;
        SceneManager.LoadScene("SampleScene", LoadSceneMode.Single);
    }

    public void Submit()
    {
        if (characterNameInput.text.Length == 0)
        {
            Debug.LogError("Character name cannot be empty");
            return;
        }

        characterName = characterNameInput.text;
        SceneManager.LoadScene("SampleScene", LoadSceneMode.Single);
    }
}

public class CharactersHttpClient
{
    const string BACKEND_URL = "http://0.0.0.0:8081";

    public async Task<CharacterDTO[]> GetCharacters()
    {
        using var request = UnityWebRequest.Get($"{BACKEND_URL}/characters");
        await SendWebRequestAsync(request);

        if (request.result != UnityWebRequest.Result.Success)
            throw new Exception(request.error);

        var jsonResponse = request.downloadHandler.text;
        var response = JsonConvert.DeserializeObject<CharacterDTO[]>(jsonResponse);
        return response;

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
