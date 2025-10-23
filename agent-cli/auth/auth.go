package auth

import (
    "context"
    "fmt"

    "github.com/AzureAD/microsoft-authentication-library-for-go/apps/public"
)

const (
    authority = "https://login.microsoftonline.com/"
)

type AuthManager struct {
    client public.Client
    scopes []string
}

func NewAuthManager(clientID, tenantID string) (*AuthManager, error) {
    // Build the authority URL
    authorityURL := authority + tenantID

    // Create public client
    client, err := public.New(clientID, public.WithAuthority(authorityURL))
    if err != nil {
        return nil, fmt.Errorf("failed to create client: %w", err)
    }

    am := &AuthManager{
        client: client,
        scopes: []string{fmt.Sprintf("api://%s/access_as_user", clientID)},
    }

    return am, nil
}


func (am *AuthManager) GetAccessToken(ctx context.Context) (string, error) {
    // Try to get accounts from cache
    accounts, err := am.client.Accounts(ctx)
    if err != nil {
        return "", fmt.Errorf("failed to get accounts: %w", err)
    }

    var result public.AuthResult

    // Try silent authentication first
    if len(accounts) > 0 {
        result, err = am.client.AcquireTokenSilent(ctx, am.scopes, public.WithSilentAccount(accounts[0]))
        if err == nil {
            fmt.Printf("✓ Authenticated as: %s\n\n", accounts[0].PreferredUsername)
            return result.AccessToken, nil
        }
    }

    // Silent auth failed, use device code flow
    fmt.Println("Authenticating with Microsoft...")

    deviceCode, err := am.client.AcquireTokenByDeviceCode(ctx, am.scopes)
    if err != nil {
        return "", fmt.Errorf("authentication failed: %w", err)
    }

    // Display the device code message to the user
    fmt.Println(deviceCode.Result.Message)

    // Wait for the user to authenticate
    result, err = deviceCode.AuthenticationResult(ctx)
    if err != nil {
        return "", fmt.Errorf("authentication failed: %w", err)
    }

    fmt.Printf("\n✓ Authentication successful!\n")
    fmt.Printf("Logged in as: %s\n\n", result.Account.PreferredUsername)

    return result.AccessToken, nil
}

func (am *AuthManager) SignOut(ctx context.Context) error {
    accounts, err := am.client.Accounts(ctx)
    if err != nil {
        return fmt.Errorf("failed to get accounts: %w", err)
    }

    for _, account := range accounts {
        if err := am.client.RemoveAccount(ctx, account); err != nil {
            return fmt.Errorf("failed to remove account: %w", err)
        }
    }

    fmt.Println("✓ Signed out successfully")
    return nil
}

func (am *AuthManager) GetCurrentUser(ctx context.Context) string {
    accounts, err := am.client.Accounts(ctx)
    if err != nil || len(accounts) == 0 {
        return ""
    }
    return accounts[0].PreferredUsername
}
