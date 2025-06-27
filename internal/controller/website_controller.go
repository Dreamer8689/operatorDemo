package controller

import (
	"context"

	websitev1 "github.com/Dreamer8689/appoperator/api/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
)

// WebsiteReconciler reconciles a Website object
type WebsiteReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=website.tomcat.io,resources=websites,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=website.tomcat.io,resources=websites/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=website.tomcat.io,resources=websites/finalizers,verbs=update
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete

func (r *WebsiteReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// 获取 Website 实例
	website := &websitev1.Website{}
	err := r.Get(ctx, req.NamespacedName, website)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		return ctrl.Result{}, err
	}

	// 确保 Deployment 存在
	deployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: website.Name, Namespace: website.Namespace}, deployment)
	if err != nil && errors.IsNotFound(err) {
		// 创建新的 Deployment
		dep := r.deploymentForWebsite(website)
		if err = r.Create(ctx, dep); err != nil {
			log.Error(err, "Failed to create Deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	} else if err != nil {
		log.Error(err, "Failed to get Deployment")
		return ctrl.Result{}, err
	}

	// 更新 Deployment
	if *deployment.Spec.Replicas != website.Spec.Replicas {
		deployment.Spec.Replicas = &website.Spec.Replicas
		if err = r.Update(ctx, deployment); err != nil {
			log.Error(err, "Failed to update Deployment")
			return ctrl.Result{}, err
		}
		return ctrl.Result{Requeue: true}, nil
	}

	// 更新状态
	website.Status.AvailableReplicas = deployment.Status.AvailableReplicas
	if err = r.Status().Update(ctx, website); err != nil {
		log.Error(err, "Failed to update Website status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// deploymentForWebsite 返回一个 Website 的 Deployment 对象
func (r *WebsiteReconciler) deploymentForWebsite(website *websitev1.Website) *appsv1.Deployment {
	labels := map[string]string{
		"app": website.Name,
	}
	replicas := website.Spec.Replicas

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      website.Name,
			Namespace: website.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: website.Spec.Image,
						Name:  "website",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 80,
							Name:          "http",
						}},
					}},
				},
			},
		},
	}
	ctrl.SetControllerReference(website, dep, r.Scheme)
	return dep
}

// SetupWithManager sets up the controller with the Manager.
func (r *WebsiteReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&websitev1.Website{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
