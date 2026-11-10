package integrations

import (
	"context"
	"net/http/httptest"
	"testing"
)

func TestOAuthPassThroughAuthenticatorRequiresRequestToken(t *testing.T) {
	authenticator := NewOAuthPassThroughAuthenticator()

	withToken := httptest.NewRequest("GET", "https://example.com", nil)
	withToken = withToken.WithContext(context.WithValue(withToken.Context(), "pass_through_oauth_token", "Bearer user-token"))
	if err := authenticator.Authenticate(withToken); err != nil {
		t.Fatalf("Authenticate with request token: %v", err)
	}
	if got := withToken.Header.Get("Authorization"); got != "Bearer user-token" {
		t.Fatalf("Authorization header = %q, want %q", got, "Bearer user-token")
	}

	withoutToken := httptest.NewRequest("GET", "https://example.com", nil)
	if err := authenticator.Authenticate(withoutToken); err == nil {
		t.Fatal("Authenticate without request token succeeded")
	}
}
