using System.Collections;
using System.Collections.Generic;
using UnityEngine;

public class PlayerJoinedDTO : DTO
{
    public string id;
    public AbilityDTO[] abilities;
}

public class AbilityDTO
{
   public string id;
   public string name;
   public float range;
}

