package controllers

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

func testSampleController() {
	It("should make configmap", func() {
		ns := testNamespacePrefix + "1"
		Expect(k8sClient.Create(context.Background(), &corev1.Namespace{
			ObjectMeta: ctrl.ObjectMeta{Name: ns},
		})).ShouldNot(HaveOccurred())
		defer k8sClient.Delete(context.Background(), &corev1.Namespace{
			ObjectMeta: ctrl.ObjectMeta{Name: ns},
		})
	})
}
