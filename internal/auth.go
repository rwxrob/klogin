package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/coreos/go-oidc/v3/oidc"
)

// ReqOIDCPassAuth grants an authenticatio token suitable for logging
// into Kubernetes by making an OIDC Resource Owner Password Flow
// request to the OIDC issuer URL (iurl). The token is returned along
// with the other response data as a map[string]any.
//
// Issuer URL (iurl) is the OIDC URL at which the discovery URL
// .well-known/openid-configuration endpoint will be
// queried to discover the token URL and other data.
//
// Client ID (cid) is passed das client_id to issuer URL (iurl)
// identifiying an OIDC supporting application.
//
// Client secret (csec) is only required if the OIDC client
// authentication type is set to "confidential". When set to "public"
// the ClientSecret is required even though---depending on the auth
// flow---it is not necessarily a "secret" and can be safely embedded in
// application source code.
func ReqOIDCPassAuth(user, pass, iurl, cid, csec string) (map[string]any, error) {

	params := url.Values{}
	params.Add("grant_type", "password")
	params.Add("client_id", cid)
	params.Add("client_secret", csec)
	params.Add("scope", "openid")
	params.Add("username", user)
	params.Add("password", pass)
	preader := strings.NewReader(params.Encode())

	client := &http.Client{}
	ctx := oidc.ClientContext(context.Background(), client)
	prov, err := oidc.NewProvider(ctx, iurl)
	if err != nil {
		return nil, err
	}
	tokurl := prov.Endpoint().TokenURL

	req, err := http.NewRequest(`POST`, tokurl, preader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200: // do nothing
	case 403, 401:
		return nil, fmt.Errorf("%v authorization denied for server for %q", resp.StatusCode, user)
	default:
		return nil, fmt.Errorf("%v unexpected status code from server for %q", resp.StatusCode, user)
	}

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read OIDC password token response for %q", user)
	}

	rresp := new(map[string]any)
	if err := json.Unmarshal(buf, rresp); err != nil {
		return nil, err
	}

	return *rresp, nil
}
