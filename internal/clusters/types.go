package clusters

// Cluster contains only the high-level meta-data for obtaining an
// OIDC Resource Owner Password Flow (grant_type=password). The
// OIDCIssuerURL is used for OIDC discovery of other URLS, etc. and for
// requesting id_tokens for use in K8S OIDC authentication.
type Cluster struct {
	Name          string
	APIServerURL  string
	OIDCIssuerURL string
	DefNamespace  string
	ClientID      string
	ClientSecret  string
	CA            []byte
}
