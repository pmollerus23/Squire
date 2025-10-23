namespace AgentMiddleware.Server.Models;

public class ConversationMetadata
{
    public int Id { get; set; }
    public int UserId { get; set; }
    public required string AzureThreadId { get; set; }
    public string? Title { get; set; }
    public DateTime CreatedAt { get; set; }
    public DateTime LastMessageAt { get; set; }
    
    public User User { get; set; } = null!;
}