using UnityEngine;

[CreateAssetMenu(fileName = "Ability", menuName = "ScriptableObjects/CreateAbility")]
public class Ability : ScriptableObject
{
    public new string name;
    public Sprite icon;
}
