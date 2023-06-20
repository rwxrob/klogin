To get these certificate authority PEM certificates the most reliable method is to grab them from the ConfigMap inside of a configured cluster. (Note that `openssl s_client -connect` methods will return the *server* cert, NOT the actual CA, which will work but is incorrect.)

```sh
kubectl get cm -n kube-system kube-root-ca.crt -o json | jq -r '.data."ca.crt"'
```

Note that you may have to temporarily `--insecure-skip-tls-verify` to make the request for the ConfigMap.

