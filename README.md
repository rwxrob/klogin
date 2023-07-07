# Kubernetes `klogin` command

* Authenticate using OIDC Resource Owner Password Flow (usually backed by LDAP).
* Hard-coded cluster authentication info (`internal/clusters`).
* Assumes host already trusts OIDC issuer URL (TLS certificate).
* Always sets `certificate-authority-data` per cluster.

## Differences from previous versions

* Uses TLS certificate validation for all OIDC queries preventing potential man
  in the middle attacks (instead of `insecure-skip-tls-verify` which so many do
  by default).

* Corrects a problem preventing Kubernetes dashboard from working with `kubectl
  proxy` in previous versions.

* Allows easily changing between different user names for a given cluster.

* Respects existing configurations (other than those reserved in
  `internal/clusters`) including the Namespace of any context.

* Uses simplified 24-hour OIDC JWT (`id_token`) authentication and cleans up
  old authentication data from configuration file(s) for reserved
  cluster/contexts.

* Adds tab completion support for more shells than just bash (zsh, fish, and
  PowerShell).

* Uses (only) official and industry-standard Kubernetes packages:
    * `k8s.io/cli-runtime`
    * `k8s.io/client-go`
    * `github.com/coreos/go-oidc/v3/oidc`
    * `github.com/spf13/cobra`

* Observes and respects use of multi-file list path values in KUBECONFIG
  environment variable (same as kubectl).

* Contains zero `execs` of the `kubectl config` program other than to display contexts.

## Legal

Released under Apache 2.0.

## Related

* Resource Owner Password Flow with OIDC  
  <https://auth0.com/docs/authenticate/login/oidc-conformant-authentication/oidc-adoption-rop-flow>
* OAuth 2.0 and OpenID Connect Overview \| Okta Developer  
  <https://developer.okta.com/docs/concepts/oauth-openid/>
