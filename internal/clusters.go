package internal

import _ "embed"

//go:embed certs/prod.crt
var ProdCA []byte

//go:embed certs/dev.crt
var DevCA []byte

//go:embed certs/inf.crt
var InfCA []byte

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

// Clusters contains a hard-coded list of cluster information with URLs
// for authentication, the API, and all OIDC information.
var Clusters = map[string]Cluster{

	`prod`: Cluster{
		Name:          `prod`,
		APIServerURL:  `https://prod.home.rwx.gg:8443`,
		OIDCIssuerURL: `https://home.rwx.gg:8443/realms/prod`,
		ClientID:      `minikube-prod`,
		ClientSecret:  `x46aVnFlygaEdHBom1200AZm37ZTWPhe`,
		CA:            ProdCA,
	},

	`dev`: Cluster{
		Name:          `dev`,
		APIServerURL:  `https://192.168.50.117:8443`,
		OIDCIssuerURL: `https://home.rwx.gg:8443/realms/dev`,
		ClientID:      `minikube-dev`,
		ClientSecret:  `12l3IMNHqDvqpp5gm0292g1Hs7SzWgyX`,
		CA:            DevCA,
	},

	`inf`: Cluster{
		Name:          `inf`,
		APIServerURL:  `https://192.168.61.20:8443`,
		OIDCIssuerURL: `https://home.rwx.gg:8443/realms/inf`,
		ClientID:      `minikube-inf`,
		ClientSecret:  `6yw8EH1dVAx9XYKTM2mK2FZZmFhD19Hz`,
		CA:            InfCA,
	},
}
