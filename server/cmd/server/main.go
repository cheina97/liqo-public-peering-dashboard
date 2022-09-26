package main

import (
	"context"
	"flag"

	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/klog/v2"

	"liqo_dashboard/pkg/dashboard"
)

func main() {
	klog.InitFlags(nil)
	flag.Parse()
	ctx := context.Background()

	cl, err := dashboard.GetKClient(ctx)
	utilruntime.Must(err)

	klog.Fatalf("unable to start the server: %s", dashboard.SetupRouterAndServeHTTP(ctx, cl))
}
