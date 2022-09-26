package dashboard

import (
	discoveryv1alpha1 "github.com/liqotech/liqo/apis/discovery/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

// ClusterDto represents the data of a cluster which is returned to the client.
type ClusterDto struct {
	clusterID         string
	Name              string                                       `json:"name"`
	Networking        discoveryv1alpha1.PeeringConditionStatusType `json:"networking"`
	Authentication    discoveryv1alpha1.PeeringConditionStatusType `json:"authentication"`
	OutgoingPeering   discoveryv1alpha1.PeeringConditionStatusType `json:"outgoingPeering"`
	OutgoingResources *ResourceMetrics                             `json:"outgoingResources"`
	IncomingPeering   discoveryv1alpha1.PeeringConditionStatusType `json:"incomingPeering"`
	IncomingResources *ResourceMetrics                             `json:"incomingResources"`
	Age               string                                       `json:"age,omitempty"`
}

// ResourceMetrics represents the metrics of a cluster which is returned to the client.
type ResourceMetrics struct {
	UsedCpus    float64 `json:"usedCpus"`
	UsedMemory  float64 `json:"usedMemory"`
	TotalMemory float64 `json:"totalMemory"`
	TotalCpus   float64 `json:"totalCpus"`
}

// ErrorResponse is returned to the client in case of error.
type ErrorResponse struct {
	Message string `json:"message"`
	Status  int16  `json:"status"`
}

func fromForeignCluster(fc *discoveryv1alpha1.ForeignCluster) *ClusterDto {
	pc := peeringConditionsToMap(fc.Status.PeeringConditions)

	clusterDto := &ClusterDto{
		Name:            fc.Name,
		clusterID:       fc.Spec.ClusterIdentity.ClusterID,
		OutgoingPeering: statusOrDefault(pc, discoveryv1alpha1.OutgoingPeeringCondition),
		IncomingPeering: statusOrDefault(pc, discoveryv1alpha1.IncomingPeeringCondition),
		Networking:      statusOrDefault(pc, discoveryv1alpha1.NetworkStatusCondition),
		Authentication:  statusOrDefault(pc, discoveryv1alpha1.AuthenticationStatusCondition),
	}

	auth, found := pc[discoveryv1alpha1.AuthenticationStatusCondition]
	if found {
		clusterDto.Age = auth.LastTransitionTime.Time.String()
	}

	return clusterDto
}

func newResourceMetrics(cpuUsage, memUsage resource.Quantity, totalResources corev1.ResourceList) *ResourceMetrics {
	return &ResourceMetrics{
		UsedCpus:    cpuUsage.AsApproximateFloat64(),
		TotalCpus:   totalResources.Cpu().AsApproximateFloat64(),
		UsedMemory:  memUsage.AsApproximateFloat64(),
		TotalMemory: totalResources.Memory().AsApproximateFloat64(),
	}
}
