package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/pmollerus23/agent-cli/auth"
	"github.com/pmollerus23/agent-cli/client"
)

const (
	clientID  = "your-client-id"
	tenantID  = "your-tenant-id"
	serverURL = "http://localhost:5000"
)

func main() {
	ctx := context.Background()

	// Print banner
	printBanner()

	// Initialize auth manager
	authManager, err := auth.NewAuthManager(clientID, tenantID)
	if err != nil {
		fmt.Printf("Failed to initialize auth: %v\n", err)
		os.Exit(1)
	}

	// Get access token (will prompt if needed)
	accessToken, err := authManager.GetAccessToken(ctx)
	if err != nil {
		fmt.Printf("Authentication failed: %v\n", err)
		os.Exit(1)
	}

	// Create API client
	apiClient := client.NewAgentClient(serverURL, accessToken)

	// Run main loop
	if err := runChatLoop(ctx, apiClient, authManager); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

func printBanner() {
	banner := `
╔════════════════════════════════════════════════════════════╗
║          Agent Middleware Console Client v1.0              ║
╚════════════════════════════════════════════════════════════╝
`
	fmt.Println(banner)
}

func runChatLoop(ctx context.Context, apiClient *client.AgentClient, authManager *auth.AuthManager) error {
	scanner := bufio.NewScanner(os.Stdin)
	var currentThreadID *string

	fmt.Println("Chat with your agent (type '/help' for commands)\n")

	for {
		fmt.Print("You: ")
		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())
		if input == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(input, "/") {
			if err := handleCommand(ctx, input, apiClient, authManager, &currentThreadID); err != nil {
				if err.Error() == "exit" {
					return nil
				}
				fmt.Printf("Error: %v\n\n", err)
			}
			continue
		}

		// Send message to agent
		response, err := apiClient.SendMessage(ctx, input, currentThreadID)
		if err != nil {
			fmt.Printf("Error: %v\n\n", err)
			continue
		}

		// Update current thread ID
		currentThreadID = &response.ThreadID

		fmt.Printf("Agent: %s\n\n", response.Message)
	}

	return scanner.Err()
}

func handleCommand(ctx context.Context, cmd string, apiClient *client.AgentClient, authManager *auth.AuthManager, currentThreadID **string) error {
	parts := strings.Fields(cmd)
	if len(parts) == 0 {
		return nil
	}

	switch parts[0] {
	case "/help":
		printHelp()
		return nil

	case "/exit", "/quit":
		fmt.Println("Goodbye!")
		return fmt.Errorf("exit")

	case "/logout":
		if err := authManager.SignOut(ctx); err != nil {
			return err
		}
		return fmt.Errorf("exit")

	case "/new":
		*currentThreadID = nil
		fmt.Println("✓ Started new conversation\n")
		return nil

	case "/history":
		conversations, err := apiClient.ListConversations(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("\nYour conversations:\n")
		for _, conv := range conversations {
			fmt.Printf("  - %s (Thread: %s)\n", conv.Title, conv.ThreadID[:8]+"...")
		}
		fmt.Println()
		return nil

	case "/profile":
		profile, err := apiClient.GetProfile(ctx)
		if err != nil {
			return err
		}
		fmt.Printf("\nYour Profile:\n")
		fmt.Printf("  Instructions: %s\n", profile.PreferredAgentInstructions)
		fmt.Println()
		return nil

	case "/whoami":
		user := authManager.GetCurrentUser(ctx)
		if user == "" {
			fmt.Println("Not authenticated")
		} else {
			fmt.Printf("Logged in as: %s\n\n", user)
		}
		return nil

	default:
		fmt.Printf("Unknown command: %s (type /help for available commands)\n\n", parts[0])
		return nil
	}
}

func printHelp() {
	help := `
Available commands:
  /help       - Show this help message
  /exit       - Exit the application
  /quit       - Exit the application
  /logout     - Sign out and exit
  /new        - Start a new conversation
  /history    - List your past conversations
  /profile    - View your profile settings
  /whoami     - Show current user

Just type your message to chat with the agent.
`
	fmt.Println(help)
}
