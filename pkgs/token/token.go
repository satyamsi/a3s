package token

import (
	"crypto"
	"fmt"
	"strings"
	"time"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v4"
	"go.aporeto.io/a3s/pkgs/permissions"
)

// An Issuer is an object that can issue IdentityToken.
type Issuer interface {
	Issue() *IdentityToken
}

// A Source represents the authentication source info
// used to derive an IdentityToken.
type Source struct {
	Type      string `json:"type"`
	Namespace string `json:"namespace,omitempty"`
	Name      string `json:"name,omitempty"`
}

// An IdentityToken represents a normalized identity token.
type IdentityToken struct {
	Identity     []string                  `json:"identity"`
	Opaque       map[string]string         `json:"opaque,omitempty"`
	Restrictions *permissions.Restrictions `json:"restrictions,omitempty"`
	Source       Source                    `json:"-"`

	jwt.RegisteredClaims
}

// NewIdentityToken returns a new IdentityToken with the
// given Source.
func NewIdentityToken(source Source) *IdentityToken {
	return &IdentityToken{
		Source: source,
	}
}

// Parse returns a validated IdentityToken from the given token string.
func Parse(tokenString string, keychain *JWKS, trustedIssuer string, requiredAudience string) (*IdentityToken, error) {

	t := &IdentityToken{}
	token, err := jwt.ParseWithClaims(tokenString, t, makeKeyFunc(keychain))
	if err != nil {
		return nil, fmt.Errorf("unable to parse jwt: %w", err)
	}

	claims := token.Claims.(*IdentityToken)

	for _, c := range claims.Identity {
		switch {
		case strings.HasPrefix(c, "@sourcename="):
			claims.Source.Name = strings.TrimPrefix(c, "@sourcename=")
		case strings.HasPrefix(c, "@sourcenamespace="):
			claims.Source.Namespace = strings.TrimPrefix(c, "@sourcenamespace=")
		case strings.HasPrefix(c, "@sourcetype="):
			claims.Source.Type = strings.TrimPrefix(c, "@sourcetype=")
		}
	}

	if claims.Source.Type == "" {
		return nil, fmt.Errorf("invalid token: missing @sourcetype in identity claims")
	}

	if claims.Issuer != trustedIssuer {
		return nil, fmt.Errorf("issuer '%s' is not acceptable. want '%s'", claims.Issuer, trustedIssuer)
	}

	if !claims.VerifyAudience(requiredAudience, false) {
		return nil, fmt.Errorf("audience '%s' is not acceptable. want '%s'", claims.Audience, requiredAudience)
	}

	return t, nil
}

// JWT returns the signed JWT string.
func (t *IdentityToken) JWT(key crypto.PrivateKey, kid string, issuer string, audience jwt.ClaimStrings, exp time.Time, cloak []string) (string, error) {

	t.ID = uuid.Must(uuid.NewV4()).String()
	t.IssuedAt = jwt.NewNumericDate(time.Now())
	t.Issuer = issuer
	t.Audience = audience

	if !exp.IsZero() {
		t.ExpiresAt = jwt.NewNumericDate(exp)
	}

	if len(cloak) > 0 {
		var cloakedIdentity []string
		for _, claim := range t.Identity {
			for _, c := range cloak {
				if strings.HasPrefix(claim, c) {
					cloakedIdentity = append(cloakedIdentity, claim)
				}
			}
		}
		t.Identity = cloakedIdentity
	}

	if t.Source.Type == "" {
		return "", fmt.Errorf("invalid identity token: missing source type")
	}

	t.Identity = append(t.Identity, fmt.Sprintf("@sourcetype=%s", t.Source.Type))

	if t.Source.Namespace != "" {
		t.Identity = append(t.Identity, fmt.Sprintf("@sourcenamespace=%s", t.Source.Namespace))
	}

	if t.Source.Name != "" {
		t.Identity = append(t.Identity, fmt.Sprintf("@sourcename=%s", t.Source.Name))
	}

	j := jwt.NewWithClaims(jwt.SigningMethodES256, t)

	if kid != "" {
		j.Header["kid"] = kid
	}

	return j.SignedString(key)
}

// Restrict applies the given permissions to the token. If the token is not already restricted
// the restrictions will be applied as is. If it is already restricted, the new restrictions will
// be applied over the existing ones, and the function will return an error if the requested
// restrictions break the limits of the current ones.
func (t *IdentityToken) Restrict(restrictions permissions.Restrictions) (err error) {

	if t.Restrictions == nil {
		t.Restrictions = &permissions.Restrictions{}
	}

	if t.Restrictions.Namespace, err = t.Restrictions.RestrictNamespace(restrictions.Namespace); err != nil {
		return err
	}

	if t.Restrictions.Networks, err = t.Restrictions.RestrictNetworks(restrictions.Networks); err != nil {
		return err
	}

	if t.Restrictions.Permissions, err = t.Restrictions.RestrictPermissions(restrictions.Permissions); err != nil {
		return err
	}

	if t.Restrictions.Zero() {
		t.Restrictions = nil
	}

	return nil
}
