/*
The clusters package is designed to be customized with hard-coded values
for OIDC cluster authentication and embedded CA certificates.

To obtain reliable copies of certificate authority PEM files fetch them
from the ConfigMap inside every configured cluster using the following
command:

		kubectl get cm -n kube-system kube-root-ca.crt -o json | jq -r '.data."ca.crt"'

This command works for self-signed Kubernetes clusters as well as those
that have been bootstrapped with custom certificates. Note that `openssl
s_client -connect` methods (as mistakenly documented in StackExchange
and elsewhere) will wrongly return the *server* cert, NOT the actual CA,
which will work but is incorrect and can be verified by examining the
content of the cert with something like the following:

		openssl x509 -in kube-root-ca.crt -noout -text

Note that you may have to temporarily enable
`--insecure-skip-tls-verify` to make the request for the ConfigMap but
this will be removed once the `certificate-authority-data` has been
properly embedded into the cluster configuration.

*/
package clusters
