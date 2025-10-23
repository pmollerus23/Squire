namespace AgentMiddleware.Server.Models;

public class UserProfile
{
    public int Id { get; set; }
    public int UserId { get; set; }
    public string? PreferredAgentInstructions { get; set; }
    public string? CustomWorkflowsJson { get; set; }
    public DateTime UpdatedAt { get; set; }
    
    public User User { get; set; } = null!;
}