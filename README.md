# Template for `klogin` command to authenticate using OIDC Resource Owner Password Flow

Clone from template and make custom changes to `internal/clusters` package.

* Assumes OIDC provider TLS is signed by trusted CA on host.
* Assumes Kubernetes API server CA is not in host CA trust chain.

## Legal

Released under Apache 2.0.

## Related

* Resource Owner Password Flow with OIDC  
  <https://auth0.com/docs/authenticate/login/oidc-conformant-authentication/oidc-adoption-rop-flow>
* OAuth 2.0 and OpenID Connect Overview \| Okta Developer  
  <https://developer.okta.com/docs/concepts/oauth-openid/>
