package dashboard

import (
	// generic import.
	"context"
	"testing"

	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	sharingv1alpha1 "github.com/liqotech/liqo/apis/sharing/v1alpha1"
	liqoconsts "github.com/liqotech/liqo/pkg/consts"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	metricsv1beta1 "k8s.io/metrics/pkg/apis/metrics/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

var (
	ctx               context.Context
	clientBuilder     fake.ClientBuilder
	clusterID         string
	shadowPodsMetrics []metricsv1beta1.PodMetrics
)

var _ = BeforeSuite(func() {
	scheme := runtime.NewScheme()
	ctx = context.Background()
	shadowPodsMetrics = []metricsv1beta1.PodMetrics{}
	foreignCluster := discoveryv1alpha1.ForeignCluster{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				"discovery.liqo.io/cluster-id": "cluster1",
			},
		},
	}
	nodeMetrics := metricsv1beta1.NodeMetrics{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				"liqo.io/remote-cluster-id": "cluster1",
			},
		},
		Usage: map[corev1.ResourceName]resource.Quantity{
			corev1.ResourceCPU:    resource.MustParse("1000m"),
			corev1.ResourceMemory: resource.MustParse("2000M"),
		},
	}
	resourceOfferOut := sharingv1alpha1.ResourceOffer{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				"liqo.io/originID": "cluster1",
			},
		},
		Spec: sharingv1alpha1.ResourceOfferSpec{
			ResourceQuota: corev1.ResourceQuotaSpec{
				Hard: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse("10000m"),
					corev1.ResourceMemory: resource.MustParse("10000M"),
				},
			},
		},
	}
	resourceOfferIn := sharingv1alpha1.ResourceOffer{
		ObjectMeta: v1.ObjectMeta{
			Name: "rome",
			Labels: map[string]string{
				"discovery.liqo.io/cluster-id": "cluster1",
			},
		},
		Spec: sharingv1alpha1.ResourceOfferSpec{
			ResourceQuota: corev1.ResourceQuotaSpec{
				Hard: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse("1000m"),
					corev1.ResourceMemory: resource.MustParse("20000M"),
				},
			},
		},
	}
	podMetrics1 := metricsv1beta1.PodMetrics{
		ObjectMeta: v1.ObjectMeta{
			Name: "pod1",
			Labels: map[string]string{
				"virtualkubelet.liqo.io/origin": "cluster1",
			},
		},
		Containers: []metricsv1beta1.ContainerMetrics{
			{
				Usage: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("1000M"),
				},
			},
		},
	}
	podMetrics2 := metricsv1beta1.PodMetrics{
		ObjectMeta: v1.ObjectMeta{
			Name: "pod2",
			Labels: map[string]string{
				"virtualkubelet.liqo.io/origin": "cluster1",
			},
		},
		Containers: []metricsv1beta1.ContainerMetrics{
			{
				Usage: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse("100m"),
					corev1.ResourceMemory: resource.MustParse("1000M"),
				},
			},
		},
	}
	node := corev1.Node{
		ObjectMeta: v1.ObjectMeta{
			Labels: map[string]string{
				"liqo.io/remote-cluster-id": "cluster1",
			},
			Name: "node1",
		},
	}
	shadowPod1 := corev1.Pod{
		ObjectMeta: v1.ObjectMeta{
			Name: "shadowpod1",
			Labels: map[string]string{
				liqoconsts.LocalPodLabelKey: "true",
			},
		},
		Spec: corev1.PodSpec{
			NodeName: "node1",
		},
	}
	shadowPodMetrics1 := metricsv1beta1.PodMetrics{
		ObjectMeta: v1.ObjectMeta{
			Name: "shadowpod1",
			Labels: map[string]string{
				liqoconsts.LocalPodLabelKey: "true",
			},
		},
		Containers: []metricsv1beta1.ContainerMetrics{
			{
				Usage: map[corev1.ResourceName]resource.Quantity{
					corev1.ResourceCPU:    resource.MustParse("1000m"),
					corev1.ResourceMemory: resource.MustParse("1000M"),
				},
			},
		},
	}
	clusterID = "cluster1"
	_ = discoveryv1alpha1.AddToScheme(scheme)
	_ = sharingv1alpha1.AddToScheme(scheme)
	_ = metricsv1beta1.AddToScheme(scheme)
	_ = corev1.AddToScheme(scheme)
	clientBuilder = *fake.
		NewClientBuilder().
		WithScheme(scheme).
		WithObjects(
			&foreignCluster,
			&nodeMetrics,
			&resourceOfferOut,
			&resourceOfferIn,
			&podMetrics1,
			&podMetrics2,
			&node,
			&shadowPod1,
			&shadowPodMetrics1)
	shadowPodsMetrics = append(shadowPodsMetrics, shadowPodMetrics1)
})

func TestTest(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Test Suite")
}
