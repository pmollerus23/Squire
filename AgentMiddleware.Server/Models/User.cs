namespace AgentMiddleware.Server.Models;

public class User
{
    public int Id { get; set; }
    public required string EntraObjectId { get; set; } // Azure Entra user Object ID
    public string? Email { get; set; } // Cached from Azure token
    public string? Username { get; set; } // Cached from Azure token
    public DateTime CreatedAt { get; set; }

    public UserProfile? Profile { get; set; }
    public ICollection<ConversationMetadata> Conversations { get; set; } = new List<ConversationMetadata>();
}
