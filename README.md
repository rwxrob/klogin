# Customizable `kubectl-login` plugin

Clone from template and make custom changes to the following to match your multi-cluster environment:

* `intern/clusters`
* `intern/certs`

## Considerations

* Only supports (deprecated) `grant_type=password` flow (at the moment).
* Assumes OIDC provider TLS is signed by trusted CA on host.
* Assumes Kubernetes API server CA is not in host CA trust chain.

## Legal

Released under Apache 2.0.
