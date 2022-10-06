# Liqo Resources Dashboard Server

The server is written in golang and to run it locally you need to execute the run command in a console where you previously exported a KUBECONFIG variable.

```bash
export KUBECONFIG="path/to/kubeconfig"
go run cmd/server/main.go
```

The exported kubeconfig must be able to list, watch and get all the following resources:

- Pods
- Nodes
- PodMetrics (metrics.k8s.io)
- ForeignClusters (discovery.liqo.io)
- ResourceOffers (sharing.liqo.io)

## Code Linting

To keep the code cleaned we use a linter called golangci-lint, you can learn more about it [here](https://github.com/golangci/golangci-lint). If you don't want to install golangci-lint you can simply run it in a docker container:

```bash
docker run --rm -v $(pwd):/app -w /app golangci/golangci-lint:v1.50.0 golangci-lint run -v
```

## Tests

[Ginkgo](https://github.com/onsi/ginkgo) has been used to write tests.
