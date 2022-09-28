package controllers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes/scheme"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	randomv1 "github.com/superorbital/random-number-controller/api/v1"
	ctrl "sigs.k8s.io/controller-runtime"
)

var err error

func TestMain(m *testing.M) {
	err = randomv1.AddToScheme(scheme.Scheme)
	if err != nil {
		panic(err)
	}
	opts := zap.Options{
		Development: true,
	}
	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))
}
func TestCreateRandomNumber(t *testing.T) {

	fakeClient := fake.NewClientBuilder().WithRuntimeObjects(&randomv1.RandomNumber{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "my-rn",
			Namespace: "default",
		},
		Spec: randomv1.RandomNumberSpec{
			Entropy: randomv1.High,
		},
	}).WithScheme(
		scheme.Scheme,
	).Build()

	reconciler := RandomNumberReconciler{
		Scheme: scheme.Scheme,
		Client: fakeClient,
	}
	_, err = reconciler.Reconcile(context.TODO(), ctrl.Request{
		NamespacedName: types.NamespacedName{
			Name:      "my-rn",
			Namespace: "default",
		}})
	if err != nil {
		t.Fatal(err)
	}

	var configMap corev1.ConfigMap
	err = fakeClient.Get(context.TODO(), client.ObjectKey{
		Namespace: "default",
		Name:      "my-rn",
	}, &configMap)
	if err != nil {
		t.Fatal(err)
	}
	assert.EqualValues(t, configMap.Data["random"], "4", "The random number is 4")
}
