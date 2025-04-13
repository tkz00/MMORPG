using System.Collections.Generic;
using Newtonsoft.Json;

public class GameStateDTO : DTO
{
    public List<CharacterDTO> players;
    public List<ProjectileDTO> projectiles;
    [JsonProperty("area_effects")]
    public List<AreaEffectDTO> aoEs;
    public List<CharacterDTO> npcs;
    [JsonProperty("entities_to_destroy")]
    public List<string> entitiesToDestroy;
}
