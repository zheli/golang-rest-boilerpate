package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"

	"github.com/example/golang-rest-boilerplate/internal/config"
)

// GoogleUser represents user information returned by Google.
type GoogleUser struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

// GoogleOAuthService handles Google OAuth interactions.
type GoogleOAuthService struct {
	config *oauth2.Config
}

// NewGoogleOAuthService constructs a GoogleOAuthService.
func NewGoogleOAuthService(cfg *config.Config) *GoogleOAuthService {
	return &GoogleOAuthService{
		config: &oauth2.Config{
			RedirectURL:  cfg.GoogleRedirectURL,
			ClientID:     cfg.GoogleClientID,
			ClientSecret: cfg.GoogleClientSecret,
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
	}
}

// AuthCodeURL returns the Google OAuth authorization URL for the given state.
func (s *GoogleOAuthService) AuthCodeURL(state string) string {
	return s.config.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

// Exchange converts an authorization code into a token.
func (s *GoogleOAuthService) Exchange(ctx context.Context, code string) (*oauth2.Token, error) {
	return s.config.Exchange(ctx, code)
}

// Client returns an HTTP client authorized with the given token.
func (s *GoogleOAuthService) Client(ctx context.Context, token *oauth2.Token) *http.Client {
	return s.config.Client(ctx, token)
}

// FetchUserInfo retrieves user information from Google using the provided token.
func (s *GoogleOAuthService) FetchUserInfo(ctx context.Context, token *oauth2.Token) (*GoogleUser, error) {
	client := s.Client(ctx, token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch google user info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("google api returned status %s", resp.Status)
	}

	var user GoogleUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, fmt.Errorf("failed to decode google user info: %w", err)
	}

	return &user, nil
}
