using System.Collections.Generic;

public class GameStateDTO : DTO {
	public List<CharacterDTO> players;
	public List<ProjectileDTO> projectiles;
	public List<CharacterDTO> npcs;
}
