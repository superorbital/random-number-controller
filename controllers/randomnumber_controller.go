/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	randomv1 "github.com/superorbital/random-number-controller/api/v1"
)

var RandomNumber = "4" // chosen by fair dice roll, guaranteed to be random

// RandomNumberReconciler reconciles a RandomNumber object
type RandomNumberReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=random.superorbital.io,resources=randomnumbers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=random.superorbital.io,resources=randomnumbers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=random.superorbital.io,resources=randomnumbers/finalizers,verbs=update
//+kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

// Reconcile creates a random configmap from a RandomNumber
func (r *RandomNumberReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	var rn randomv1.RandomNumber
	if err := r.Get(ctx, req.NamespacedName, &rn); err != nil {
		logger.Error(err, "unable to fetch RandomNumber")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	key := types.NamespacedName{
		Name:      rn.GetName(),
		Namespace: rn.GetNamespace(),
	}
	// By default, don't create a new configmap
	create := false
	var configMap corev1.ConfigMap
	if err := r.Get(ctx, key, &configMap); err != nil {
		if !errors.IsNotFound(err) {
			return ctrl.Result{}, err
		}
		create = true
		configMap = corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      rn.GetName(),
				Namespace: rn.GetNamespace(),
			},
		}
	}
	// set our random number
	configMap.Data = map[string]string{"random": RandomNumber}
	if err := controllerutil.SetControllerReference(&rn, &configMap, r.Scheme); err != nil {
		logger.Error(err, "unable to set controller reference on configmap")
		return ctrl.Result{}, err
	}

	var err error
	if create {
		logger.Info("Creating configmap")
		err = r.Create(ctx, &configMap)
	} else {
		logger.Info("Updating configmap")
		err = r.Update(ctx, &configMap)
	}
	if err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *RandomNumberReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&randomv1.RandomNumber{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}
