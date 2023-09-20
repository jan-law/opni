package monitoring

import (
	"github.com/rancher/opni/pkg/logger"
	"github.com/rancher/opni/pkg/resources"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/util/retry"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (r *Reconciler) updateImageStatus() (bool, error) {
	lg := r.logger
	var image string
	var pullPolicy corev1.PullPolicy
	if imgOverride := r.mc.Spec.Cortex.Image.GetImageWithDefault(""); imgOverride != "" {
		image = imgOverride
	} else {
		var err error
		image, pullPolicy, err = resources.FindManagerImage(r.ctx, r.client)
		if err != nil {
			return false, err
		}
	}

	if r.mc.Status.Image != image {
		err := retry.RetryOnConflict(retry.DefaultRetry, func() error {
			if err := r.client.Get(r.ctx, client.ObjectKeyFromObject(r.mc), r.mc); err != nil {
				return err
			}
			r.mc.Status.Image = image
			r.mc.Status.ImagePullPolicy = pullPolicy
			return r.client.Status().Update(r.ctx, r.mc)
		})
		if err != nil {
			lg.Error("failed to update monitoring cluster status", logger.Err(err))
			return false, err
		}
		return true, nil
	}
	return false, nil
}
