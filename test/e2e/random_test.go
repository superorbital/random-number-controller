//go:build e2e
// +build e2e

package e2e

import (
	"context"
	"testing"
	"time"

	. "github.com/onsi/gomega"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"k8s.io/apimachinery/pkg/api/errors"

	randomv1 "github.com/superorbital/random-number-controller/api/v1"
)

var (
	timeout  = time.Second * 10
	duration = time.Second * 10
	interval = time.Millisecond * 250
)

func TestThatThingsHappen(t *testing.T) {
	g := NewWithT(t)

	env := envtest.Environment{
		UseExistingCluster: boolPointer(true),
	}
	config, err := env.Start()
	g.Expect(err).ToNot(HaveOccurred())

	randomv1.AddToScheme(scheme.Scheme)
	k8sClient, err := client.New(config, client.Options{Scheme: scheme.Scheme})
	g.Expect(err).ToNot(HaveOccurred())

	rn := &randomv1.RandomNumber{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "rn-e2e",
			Namespace: "default",
		},
		Spec: randomv1.RandomNumberSpec{
			Entropy: randomv1.High,
		},
	}
	g.Expect(k8sClient.Create(context.TODO(), rn)).Should(Succeed())

	cmObjectKey := client.ObjectKey{
		Name:      "rn-e2e",
		Namespace: "default",
	}
	var configMap corev1.ConfigMap

	g.Eventually(func() bool {
		err = k8sClient.Get(context.TODO(), cmObjectKey, &configMap)

		if err != nil {
			return false
		}
		return true
	}, timeout, interval).Should(BeTrue())

	g.Expect(configMap.Data["random"]).Should(Equal("4"))

	err = k8sClient.Delete(context.TODO(), rn)
	g.Expect(err).ToNot(HaveOccurred())

	// cleanup
	g.Eventually(func() bool {
		var configMapDeleted corev1.ConfigMap
		err = k8sClient.Get(context.TODO(), cmObjectKey, &configMapDeleted)
		if err == nil {
			return false
		}
		return errors.IsNotFound(err)
	}, timeout, interval).Should(BeTrue())

}

func boolPointer(b bool) *bool {
	return &b
}
