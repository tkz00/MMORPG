using TMPro;
using UnityEngine;
using UnityEngine.SceneManagement;

public class AuthenticationManager : MonoBehaviour
{
    [SerializeField] TMP_InputField characterNameInput;

    static string characterName;
    public static string CharacterName { get { return characterName; } }

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
