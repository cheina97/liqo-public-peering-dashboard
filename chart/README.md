## Values

| Key | Type | Default | Description |
|-----|------|---------|-------------|
| host | string | `""` | Define a host for the ingress. It must match with the TLS certificate's host if present |
| tls.secretName | string | `""` | Name of the secret containing the TLS certificate |
| backend.port | number | `8080` | backend listening port |
| backend.deployment.imageTag | string | `backend-v1` | server docker image tag |