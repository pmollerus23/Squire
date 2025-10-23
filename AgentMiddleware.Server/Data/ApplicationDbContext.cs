using Microsoft.EntityFrameworkCore;
using AgentMiddleware.Server.Models;

namespace AgentMiddleware.Server.Data;

public class ApplicationDbContext : DbContext
{
    public ApplicationDbContext(DbContextOptions<ApplicationDbContext> options) : base(options) { }

    public DbSet<User> Users => Set<User>();
    public DbSet<UserProfile> UserProfiles => Set<UserProfile>();
    public DbSet<ConversationMetadata> ConversationMetadata => Set<ConversationMetadata>();

    protected override void OnModelCreating(ModelBuilder modelBuilder)
    {
        modelBuilder.Entity<User>(entity =>
        {
            entity.HasIndex(e => e.EntraObjectId).IsUnique();
            entity.HasOne(e => e.Profile)
                  .WithOne(e => e.User)
                  .HasForeignKey<UserProfile>(e => e.UserId)
                  .OnDelete(DeleteBehavior.Cascade);
        });

        modelBuilder.Entity<ConversationMetadata>(entity =>
        {
            entity.HasIndex(e => e.AzureThreadId);
            entity.HasOne(e => e.User)
                  .WithMany(e => e.Conversations)
                  .HasForeignKey(e => e.UserId)
                  .OnDelete(DeleteBehavior.Cascade);
        });
    }
}