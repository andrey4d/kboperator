/*
Copyright 2024 andrey4d.dev@gmail.com.

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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	kbov1alpha1 "github.com/andrey4d/kboperator/api/v1alpha1"
	"github.com/andrey4d/kboperator/internal/k8s/configmaps"
	"github.com/andrey4d/kboperator/internal/k8s/jobs"
	"github.com/andrey4d/kboperator/internal/k8s/persistence"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	kbatch "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

const (
	typeAvailableProject = "Available"
	typeDegradedProject  = "Degraded"
)

type contextKey string

const objectLogKey contextKey = "object"

// KanikoBuildReconciler reconciles a KanikoBuild object
type KanikoBuildReconciler struct {
	client.Client
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=kbo.k8s.dav.io,resources=,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=kbo.k8s.dav.io,resources=kanikobuilds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=kbo.k8s.dav.io,resources=kanikobuilds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the KanikoBuild object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *KanikoBuildReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	kaniko := &kbov1alpha1.KanikoBuild{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, kaniko)
	if err != nil {
		if apierrors.IsNotFound(err) {
			log.Info("Application resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Info("Failed to get Application resource. Re-running reconcile.")
		return ctrl.Result{}, err
	}

	log.Info("Application resource found", "name", kaniko.Name)

	if len(kaniko.Status.Conditions) == 0 {
		meta.SetStatusCondition(&kaniko.Status.Conditions, metav1.Condition{Type: typeAvailableProject, Status: metav1.ConditionUnknown, Reason: "Reconciling", Message: "Starting reconciliation"})
		if err = r.Status().Update(ctx, kaniko); err != nil {
			log.Error(err, "Failed to update project status")
			return ctrl.Result{}, err
		}

		if err := r.Get(ctx, req.NamespacedName, kaniko); err != nil {
			log.Error(err, "Failed to re-fetch project")
			return ctrl.Result{}, err
		}
	}

	// ConfigMaps
	result, err := r.ConfigMap(ctx, req, kaniko)
	if err != nil {
		return result, err
	}

	if kaniko.Spec.Persistence.Enabled {
		result, err = r.PersistenceVolume(ctx, req, kaniko)
		if err != nil {
			return result, err
		}
	}
	// Jobs
	result, err = r.Job(ctx, req, kaniko)
	if err != nil {
		return result, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *KanikoBuildReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&kbov1alpha1.KanikoBuild{}).
		Named("kanikobuild").
		// Owns(&kbatch.Job{}).
		Owns(&corev1.ConfigMap{}).
		Complete(r)
}

func (r *KanikoBuildReconciler) ConfigMap(ctx context.Context, req ctrl.Request, kaniko *kbov1alpha1.KanikoBuild) (ctrl.Result, error) {
	log := log.FromContext(ctx)
	found := &corev1.ConfigMap{}
	err := r.Get(ctx, r.objectKey(kaniko), found)

	if err != nil && apierrors.IsNotFound(err) {
		cm, err := configmaps.NewConfigMap(kaniko, r.Scheme).BuilderConfigMap()
		if err := r.SetErrorStatus(context.WithValue(ctx, objectLogKey, "ConfigMap"), kaniko, err); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Creating a new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)

		if err = r.Create(ctx, cm); err != nil {
			log.Error(err, "Failed to create new ConfigMap", "ConfigMap.Namespace", cm.Namespace, "ConfigMap.Name", cm.Name)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get ConfigMap")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *KanikoBuildReconciler) Job(ctx context.Context, req ctrl.Request, kaniko *kbov1alpha1.KanikoBuild) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithName("Job")
	found := &kbatch.Job{}
	err := r.Get(ctx, r.objectKey(kaniko), found)
	if err != nil && apierrors.IsNotFound(err) {
		job, err := jobs.NewJob(kaniko, r.Scheme).BuilderJob()
		if err := r.SetErrorStatus(context.WithValue(ctx, objectLogKey, "Job"), kaniko, err); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Creating a new Job", "Namespace", job.Namespace, "Name", job.Name)

		if err = r.Create(ctx, job); err != nil {
			log.Error(err, "Failed to create new Job", "Namespace", job.Namespace, "Name", job.Name)
			return ctrl.Result{}, err
		}

	} else if err != nil {
		log.Error(err, "Failed to get Job")
		// Let's return the error for the reconciliation be re-trigged again
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *KanikoBuildReconciler) SetErrorStatus(ctx context.Context, kaniko *kbov1alpha1.KanikoBuild, err error) error {
	log := log.FromContext(ctx).WithValues("object", ctx.Value("object"))
	if err != nil {
		log.Error(err, "Failed to define new resource")
		// The following implementation will update the status
		meta.SetStatusCondition(&kaniko.Status.Conditions, metav1.Condition{
			Type:   typeAvailableProject,
			Status: metav1.ConditionFalse, Reason: "Reconciling",
			Message: fmt.Sprintf("Failed to create the custom resource (%s): (%s)", kaniko.Name, err),
		})

		if err := r.Status().Update(ctx, kaniko); err != nil {
			log.Error(err, "Failed to update KanikoBuilder status")
			return err
		}
		return err
	}
	return err
}

func (r *KanikoBuildReconciler) PersistenceVolume(ctx context.Context, req ctrl.Request, kaniko *kbov1alpha1.KanikoBuild) (ctrl.Result, error) {
	log := log.FromContext(ctx).WithName("PersistenceVolume")
	found := &corev1.PersistentVolumeClaim{}
	err := r.Get(ctx, r.objectKey(kaniko), found)
	if err != nil && apierrors.IsNotFound(err) {
		pvc, err := persistence.NewPersistence(kaniko, r.Scheme).BuilderPvc()
		if err := r.SetErrorStatus(context.WithValue(ctx, objectLogKey, "PVC"), kaniko, err); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("Creating a new ConfigMap", "Namespace", pvc.Namespace, "Name", pvc.Name)

		if err = r.Create(ctx, pvc); err != nil {
			log.Error(err, "Failed to create new PVC", "Namespace", pvc.Namespace, "Name", pvc.Name)
			return ctrl.Result{}, err
		}
	} else if err != nil {
		log.Error(err, "Failed to get PVC")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *KanikoBuildReconciler) objectKey(k *kbov1alpha1.KanikoBuild) client.ObjectKey {
	return client.ObjectKey{
		Name:      k.Spec.Name,
		Namespace: k.Namespace,
	}
}
