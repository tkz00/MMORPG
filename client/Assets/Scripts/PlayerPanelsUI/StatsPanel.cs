using System.Collections.Generic;
using System.Linq;
using TMPro;
using UnityEngine;

public class StatsPanel : MonoBehaviour
{
    Dictionary<string, int> currentStats = new();
    [SerializeField] TextMeshProUGUI statsText;

    public void UpdatePanel(Dictionary<string, int> stats)
    {
        foreach (KeyValuePair<string, int> stat in stats)
        {
            currentStats[stat.Key] = stat.Value;
            statsText.text = string.Join("\n", currentStats.Select(kvp => $"{kvp.Key}: {kvp.Value}"));
        }
    }
}
