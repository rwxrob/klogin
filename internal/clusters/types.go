package clusters

// Cluster contains only the high-level meta-data for obtaining an
// OIDC Resource Owner Password Flow (grant_type=password). The
// OIDCIssuerURL is used for OIDC discovery of other URLS, etc. and for
// requesting id_tokens for use in K8S OIDC authentication.
type Cluster struct {
	Name          string // prod
	APIServerURL  string // https://192.168.50.133:8443
	OIDCIssuerURL string // https://home.rwx.gg:8443/realms/prod
	DefNamespace  string // <user> by default
	ClientID      string // minikube-prod
	ClientSecret  string // 6yw8EH1dVAx9XYKTM2mK2FZZmFhD19Hz
	CA            []byte // certificate authority data (NOT server cert)
}
