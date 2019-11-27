/*

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
	"strconv"

	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/record"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	samplecontrollerv1alpha1 "github.com/sample-controller/api/v1alpha1"
)

// SampleReconciler reconciles a Sample object
type SampleReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=samplecontroller.k8s.io,resources=samples,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=samplecontroller.k8s.io,resources=samples/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;path

func (r *SampleReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("sample", req.NamespacedName)

	var sample samplecontrollerv1alpha1.Sample
	log.Info("fetching Sample Resource")
	if err := r.Get(ctx, req.NamespacedName, &sample); err != nil {
		log.Error(err, "unable to fetch Sample")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// cleanup
	// T.B.D

	for i := 0; i < sample.Spec.Maps; i++ {

		deploy := &corev1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      strconv.Itoa(i),
				Namespace: req.Namespace,
			},
		}

		if _, err := ctrl.CreateOrUpdate(ctx, r.Client, deploy, func() error {

			// set a label
			// T.B.D

			// set the owner so that garbage collection can kicks in
			if err := ctrl.SetControllerReference(&sample, deploy, r.Scheme); err != nil {
				log.Error(err, "unable to set ownerReference from Sample to ConfigMap")
				return err
			}
			// end of ctrl.CreateOrUpdate
			return nil
		}); err != nil {

			// error handling of ctrl.CreateOrUpdate
			log.Error(err, "unable to ensure maneki is correct")
			return ctrl.Result{}, err

		}
	}

	return ctrl.Result{}, nil
}

func (r *SampleReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&samplecontrollerv1alpha1.Sample{}).
		Complete(r)
}
