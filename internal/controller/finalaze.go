package controller

// import (
// 	"fmt"

// 	"github.com/andrey4d/kboperator/api/v1alpha1"
// )

// // finalizeMemcached will perform the required operations before delete the CR.
// func (r *KanikoBuildReconciler) doFinalizerOperationsForKanikoBuilder(cr *v1alpha1.KanikoBuild) {

// 	r.Recorder.Event(cr, "Warning", "Deleting", fmt.Sprintf("Custom Resource %s is being deleted ", cr.Name))
// }

// func GetProjectFinalizer() string {
// 	return "kbo.k8s.dav.io/finalizer"
// }

/*
	// Let's add a finalizer. Then, we can define some operations which should occurs before the custom resource to be deleted.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/finalizers
	if !controllerutil.ContainsFinalizer(kaniko, GetProjectFinalizer()) {
		log.Info("Adding Finalizer for KanikoBuilder")
		if ok := controllerutil.AddFinalizer(kaniko, GetProjectFinalizer()); !ok {
			log.Error(err, "Failed to add finalizer into the custom resource")
			return ctrl.Result{Requeue: true}, nil
		}

		if err = r.Update(ctx, kaniko); err != nil {
			log.Error(err, "Failed to update custom resource to add finalizer")
			return ctrl.Result{}, err
		}
	}

	// Check if the Memcached instance is marked to be deleted, which is
	// indicated by the deletion timestamp being set.
	isMemcachedMarkedToBeDeleted := kaniko.GetDeletionTimestamp() != nil
	if isMemcachedMarkedToBeDeleted {
		if controllerutil.ContainsFinalizer(kaniko, GetProjectFinalizer()) {
			log.Info("Performing Finalizer Operations for KanikoBuilder before delete CR")

			// Let's add here an status "Downgrade" to define that this resource begin its process to be terminated.
			meta.SetStatusCondition(&kaniko.Status.Conditions, metav1.Condition{Type: typeDegradedProject,
				Status: metav1.ConditionUnknown, Reason: "Finalizing",
				Message: fmt.Sprintf("Performing finalizer operations for the custom resource: %s ", kaniko.Name)})

			if err := r.Status().Update(ctx, kaniko); err != nil {
				log.Error(err, "Failed to update KanikoBuilder status")
				return ctrl.Result{}, err
			}

			// Perform all operations required before remove the finalizer and allow
			// the Kubernetes API to remove the custom resource.
			r.doFinalizerOperationsForKanikoBuilder(kaniko)

			// TODO(user): If you add operations to the doFinalizerOperationsForMemcached method
			// then you need to ensure that all worked fine before deleting and updating the Downgrade status
			// otherwise, you should requeue here.

			// Re-fetch the memcached Custom Resource before update the status
			// so that we have the latest state of the resource on the cluster and we will avoid
			// raise the issue "the object has been modified, please apply
			// your changes to the latest version and try again" which would re-trigger the reconciliation
			if err := r.Get(ctx, req.NamespacedName, kaniko); err != nil {
				log.Error(err, "Failed to re-fetch KanikoBuilder")
				return ctrl.Result{}, err
			}

			meta.SetStatusCondition(&kaniko.Status.Conditions, metav1.Condition{Type: typeDegradedProject,
				Status: metav1.ConditionTrue, Reason: "Finalizing",
				Message: fmt.Sprintf("Finalizer operations for custom resource %s name were successfully accomplished", kaniko.Name)})

			if err := r.Status().Update(ctx, kaniko); err != nil {
				log.Error(err, "Failed to update KanikoBuilder status")
				return ctrl.Result{}, err
			}

			log.Info("Removing Finalizer for KanikoBuilder after successfully perform the operations")
			if ok := controllerutil.RemoveFinalizer(kaniko, GetProjectFinalizer()); !ok {
				log.Error(err, "Failed to remove finalizer for KanikoBuilder")
				return ctrl.Result{Requeue: true}, nil
			}

			if err := r.Update(ctx, kaniko); err != nil {
				log.Error(err, "Failed to remove finalizer for KanikoBuilder")
				return ctrl.Result{}, err
			}
		}
		return ctrl.Result{}, nil
	}
*/
