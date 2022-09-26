package dashboard

import (
	"context"
	"fmt"

	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	liqoconsts "github.com/liqotech/liqo/pkg/consts"
	liqogetters "github.com/liqotech/liqo/pkg/utils/getters"
	liqolabels "github.com/liqotech/liqo/pkg/utils/labels"
	liqorestcfg "github.com/liqotech/liqo/pkg/utils/restcfg"
	virtualkubeletconsts "github.com/liqotech/liqo/pkg/virtualKubelet/forge"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog/v2"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var (
	scheme = runtime.NewScheme()
)

func init() {
	_ = discoveryv1alpha1.AddToScheme(scheme)
	_ = sharingv1alpha1.AddToScheme(scheme)
	_ = metricsv1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
}

// GetKClient creates a kubernetes API client and returns it.
func GetKClient(ctx context.Context) (client.Client, error) {
	config := liqorestcfg.SetRateLimiter(ctrl.GetConfigOrDie())

	cl, err := client.New(config, client.Options{
		Scheme: scheme,
	})
	if err != nil {
		klog.Fatalf("error creating manager: %", err)
	}

	return cl, nil
}

// It gets all the available foreign clusters connected to the local cluster and calculate
// incoming and outgoing resources.
func getForeignClusters(ctx context.Context, cl client.Client) (*[]ClusterDto, error) {
	foreignClusterList := &discoveryv1alpha1.ForeignClusterList{}
	err := cl.List(ctx, foreignClusterList)
	if err != nil {
		klog.Errorf("error retrieving foreign clusters: %s", err)
		return nil, err
	}

	podMetricsList := &metricsv1beta1.PodMetricsList{}
	err = cl.List(ctx, podMetricsList, client.MatchingLabels{
		liqoconsts.LocalPodLabelKey: "true",
	})
	if err != nil {
		klog.Warningf("error retrieving pod metrics: %s", err)
		return nil, err
	}
	podMetricsMap := podMetricListToMap(podMetricsList.Items)

	var clusters []ClusterDto
	for i := range foreignClusterList.Items {
		clusterDto := fromForeignCluster(&foreignClusterList.Items[i])

		if isPeeringEstablished(clusterDto.OutgoingPeering) {
			klog.V(5).Infof("Calculating outgoing resources for cluster %s", clusterDto.clusterID)
			outgoingResources, err := calculateOutgoingResources(ctx, cl, clusterDto.clusterID, podMetricsMap)
			if err == nil {
				clusterDto.OutgoingResources = outgoingResources
			} // otherwise, the outgoing resources are not calculated so they are nil
		}

		if isPeeringEstablished(clusterDto.IncomingPeering) {
			incomingResources, err := calculateIncomingResources(ctx, cl, clusterDto.clusterID)
			if err == nil {
				clusterDto.IncomingResources = incomingResources
			} // otherwise, the incoming resources are not calculated so they are nil
		}

		// in this moment clusters without any resource are also added to the list but we can decide to filter them
		clusters = append(clusters, *clusterDto)
	}

	return &clusters, nil
}

// CalculateOutgoingResources The outgoing˙ resources aren't calculated, but they are simply retrieved from the metrics of the virtual node. The
// clusterID identifies the virtual node by the label liqo.io/remote-cluster-id=clusterID.
func calculateOutgoingResources(ctx context.Context, cl client.Client, clusterID string,
	shadowPodsMetrics map[string]*metricsv1beta1.PodMetrics) (*ResourceMetrics, error) {
	resourceOffer, err := liqogetters.GetResourceOfferByLabel(ctx, cl, metav1.NamespaceAll, liqolabels.RemoteLabelSelector(clusterID))
	if err != nil {
		klog.V(5).Infof("error retrieving resourceOffers: %s", err)
		return nil, err
	}

	podList, err := getOutgoingPods(ctx, cl, clusterID)
	if err != nil {
		klog.Errorf("error retrieving outgoing pods: %s", err)
		return nil, err
	}

	currentPodMetrics := []metricsv1beta1.PodMetrics{}
	for i := range podList {
		singlePodMetrics, found := shadowPodsMetrics[podList[i].Name]
		if found {
			currentPodMetrics = append(currentPodMetrics, *singlePodMetrics)
		}
	}

	cpuUsage, memUsage := aggregatePodsMetrics(currentPodMetrics)
	totalResources := resourceOffer.Spec.ResourceQuota.Hard
	return newResourceMetrics(cpuUsage, memUsage, totalResources), nil
}

// Calculates the resources that the local cluster is giving to a remote cluster identified by a given clusterID.
// In order to calculate these resources the function sums the resources consumed by all pods having the label
// virtualkubelet.liqo.io/origin=clusterID which is present only on pods that have been scheduled from the
// remote cluster.
func calculateIncomingResources(ctx context.Context, cl client.Client, clusterID string) (*ResourceMetrics, error) {
	resourceOffer, err := liqogetters.GetResourceOfferByLabel(ctx, cl, metav1.NamespaceAll, liqolabels.LocalLabelSelector(clusterID))
	if err != nil {
		klog.Warningf("error retrieving resourceOffers: %s", err)
		return nil, err
	}

	podMetricsList := &metricsv1beta1.PodMetricsList{}
	if err := cl.List(ctx, podMetricsList, client.MatchingLabels{
		virtualkubeletconsts.LiqoOriginClusterIDKey: clusterID,
	}); err != nil {
		return nil, err
	}

	cpuUsage, memUsage := aggregatePodsMetrics(podMetricsList.Items)

	totalResources := resourceOffer.Spec.ResourceQuota.Hard
	return newResourceMetrics(cpuUsage, memUsage, totalResources), nil
}

func getOutgoingPods(ctx context.Context, cl client.Client, clusterID string) ([]corev1.Pod, error) {
	nodeList := &corev1.NodeList{}
	if err := cl.List(ctx, nodeList, client.MatchingLabels{
		liqoconsts.RemoteClusterID: clusterID,
	}); err != nil {
		klog.V(5).Infof("error retrieving nodes: %s", err)
		return nil, err
	}

	if len(nodeList.Items) != 1 {
		return nil, fmt.Errorf("expected exactly one element in the list of Nodes but got %d", len(nodeList.Items))
	}

	node := nodeList.Items[0].Name
	podList := &corev1.PodList{}
	err := cl.List(ctx, podList, client.MatchingFields{
		"spec.nodeName": node,
	})
	if err != nil {
		klog.V(5).Infof("error retrieving pods: %w", err)
		return nil, fmt.Errorf("error retrieving pods: %w", err)
	}

	return podList.Items, nil
}
