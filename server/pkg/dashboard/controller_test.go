package dashboard

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ = Describe("Tests to check if the calculation functions work properly", Ordered, func() {
	var cl client.Client
	BeforeAll(func() {
		cl = clientBuilder.Build()
	})

	When("outgoing resources calculation occur", func() {
		var (
			resource *ResourceMetrics
			err      error
		)
		BeforeAll(func() {
			resource, err = calculateOutgoingResources(ctx, cl, clusterID, podMetricListToMap(shadowPodsMetrics))
		})

		It("should succeed", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("total CPU should be 10", func() {
			Expect(resource.TotalCpus).To(BeNumerically("==", 10))
		})

		It("used cpus should be 1", func() {
			Expect(resource.UsedCpus).To(BeNumerically("==", 1))
		})

		It("total memory should be 10GB", func() {
			Expect(resource.TotalMemory).To(BeNumerically("==", 1e+10))
		})

		It("used memory should be 1GB", func() {
			Expect(resource.UsedMemory).To(BeNumerically("==", 1e+09))
		})
	})

	When("incoming resources calculation occur", func() {
		var (
			resource *ResourceMetrics
			err      error
		)

		BeforeAll(func() {
			resource, err = calculateIncomingResources(ctx, cl, clusterID)
		})

		It("should succeed", func() {
			Expect(err).ToNot(HaveOccurred())
		})

		It("total CPU should be 1", func() {
			Expect(resource.TotalCpus).To(BeNumerically("==", 1))
		})

		It("used cpus should be 0.2", func() {
			Expect(resource.UsedCpus).To(BeNumerically("==", 0.2))
		})

		It("total memory should be 20GB", func() {
			Expect(resource.TotalMemory).To(BeNumerically("==", 2e+10))
		})

		It("used memory should be 2GB", func() {
			Expect(resource.UsedMemory).To(BeNumerically("==", 2e+9))
		})
	})
})
