using Newtonsoft.Json;

public class EquipItemDTO : DTO
{
    [JsonProperty("item_id")]
    public string itemId;
}

public class UnequipItemDTO : DTO
{
    [JsonProperty("item_id")]
    public string itemId;
}