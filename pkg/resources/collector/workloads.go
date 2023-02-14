package collector

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	opniloggingv1beta1 "github.com/rancher/opni/apis/logging/v1beta1"
	"github.com/rancher/opni/pkg/resources"
	opnimeta "github.com/rancher/opni/pkg/util/meta"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
)

const (
	receiversKey       = "receivers.yaml"
	mainKey            = "config.yaml"
	collectorImageRepo = "docker.io.otel"
	collectorImage     = "opentelemetry-collector-contrib"
	collectorVersion   = "0.68.0"
)

func (r *Reconciler) receiverConfig() (retData []byte, retReceivers []string, retErr error) {
	if r.collector.Spec.LoggingConfig != nil {
		retData = append(retData, []byte(templateLogAgentK8sReceiver)...)
		retReceivers = append(retReceivers, logReceiverK8s)

		receiver, data, err := r.generateDistributionReceiver()
		if err != nil {
			retErr = err
			return
		}
		retData = append(retData, data...)
		if receiver != "" {
			retReceivers = append(retReceivers, receiver)
		}
	}
	return
}

func (r *Reconciler) generateDistributionReceiver() (receiver string, retBytes []byte, retErr error) {
	config := &opniloggingv1beta1.CollectorConfig{}
	retErr = r.client.Get(r.ctx, types.NamespacedName{
		Name:      r.collector.Spec.LoggingConfig.Name,
		Namespace: r.collector.Spec.SystemNamespace,
	}, config)
	if retErr != nil {
		return
	}
	var providerReceiver bytes.Buffer
	switch config.Spec.Provider {
	case opniloggingv1beta1.LogProviderRKE:
		return logReceiverRKE, []byte(templateLogAgentRKE), nil
	case opniloggingv1beta1.LogProviderK3S:
		journaldDir := "/var/log/journald"
		if config.Spec.K3S != nil && config.Spec.K3S.LogPath != "" {
			journaldDir = config.Spec.K3S.LogPath
		}
		retErr = templateLogAgentK3s.Execute(&providerReceiver, journaldDir)
		if retErr != nil {
			return
		}
		return logReceiverK3s, providerReceiver.Bytes(), nil
	case opniloggingv1beta1.LogProviderRKE2:
		journaldDir := "/var/log/journald"
		if config.Spec.RKE2 != nil && config.Spec.RKE2.LogPath != "" {
			journaldDir = config.Spec.RKE2.LogPath
		}
		retErr = templateLogAgentRKE2.Execute(&providerReceiver, journaldDir)
		if retErr != nil {
			return
		}
		return logReceiverRKE2, providerReceiver.Bytes(), nil
	default:
		return
	}
}

func (r *Reconciler) hostLoggingVolumes() (
	retVolumeMounts []corev1.VolumeMount,
	retVolumes []corev1.Volume,
	retErr error,
) {
	config := &opniloggingv1beta1.CollectorConfig{}
	retErr = r.client.Get(r.ctx, types.NamespacedName{
		Name:      r.collector.Spec.LoggingConfig.Name,
		Namespace: r.collector.Spec.SystemNamespace,
	}, config)
	if retErr != nil {
		return
	}

	retVolumeMounts = append(retVolumeMounts, corev1.VolumeMount{
		Name:      "varlogpods",
		MountPath: "/var/log/pods",
	})
	retVolumes = append(retVolumes, corev1.Volume{
		Name: "varlogpods",
		VolumeSource: corev1.VolumeSource{
			HostPath: &corev1.HostPathVolumeSource{
				Path: "/var/log/pods",
			},
		},
	})
	switch config.Spec.Provider {
	case opniloggingv1beta1.LogProviderRKE:
		retVolumeMounts = append(retVolumeMounts, corev1.VolumeMount{
			Name:      "rancher",
			MountPath: "/var/lib/rancher/rke/log",
		})
		retVolumes = append(retVolumes, corev1.Volume{
			Name: "rancher",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: "/var/lib/rancher/rke/log",
				},
			},
		})
	case opniloggingv1beta1.LogProviderK3S:
		journaldDir := "/var/log/journald"
		if config.Spec.K3S != nil && config.Spec.K3S.LogPath != "" {
			journaldDir = config.Spec.K3S.LogPath
		}
		retVolumeMounts = append(retVolumeMounts, corev1.VolumeMount{
			Name:      "journald",
			MountPath: journaldDir,
		})
		retVolumes = append(retVolumes, corev1.Volume{
			Name: "journald",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: journaldDir,
				},
			},
		})
	case opniloggingv1beta1.LogProviderRKE2:
		journaldDir := "/var/log/journald"
		if config.Spec.RKE2 != nil && config.Spec.RKE2.LogPath != "" {
			journaldDir = config.Spec.RKE2.LogPath
		}
		retVolumeMounts = append(retVolumeMounts, corev1.VolumeMount{
			Name:      "journald",
			MountPath: journaldDir,
		})
		retVolumes = append(retVolumes, corev1.Volume{
			Name: "journald",
			VolumeSource: corev1.VolumeSource{
				HostPath: &corev1.HostPathVolumeSource{
					Path: journaldDir,
				},
			},
		})
	}
	return
}

func (r *Reconciler) configMapName() string {
	return fmt.Sprintf("%s-collector-config", r.collector.Name)
}

func (r *Reconciler) mainConfig(receivers []string) ([]byte, error) {
	config := LoggingConfig{
		Enabled:   r.collector.Spec.LoggingConfig != nil,
		Receivers: receivers,
	}

	var buffer bytes.Buffer
	err := templateMainConfig.Execute(&buffer, config)
	if err != nil {
		return buffer.Bytes(), err
	}

	return buffer.Bytes(), nil
}

func (r *Reconciler) configMap() (resources.Resource, string) {
	cm := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.configMapName(),
			Namespace: r.collector.Spec.SystemNamespace,
			Labels: map[string]string{
				resources.PartOfLabel: "opni",
			},
		},
		Data: map[string]string{},
	}

	receiverData, receivers, err := r.receiverConfig()
	if err != nil {
		return resources.Error(cm, err), ""
	}
	cm.Data[receiversKey] = string(receiverData)

	mainData, err := r.mainConfig(receivers)
	if err != nil {
		return resources.Error(cm, err), ""
	}
	cm.Data[mainKey] = string(mainData)

	combinedData := append(receiverData, mainData...)
	hash := sha256.New()
	hash.Write(combinedData)
	configHash := hex.EncodeToString(hash.Sum(nil))

	ctrl.SetControllerReference(r.collector, cm, r.client.Scheme())

	return resources.Present(cm), configHash
}

func (r *Reconciler) daemonSet(configHash string) resources.Resource {
	imageSpec := r.imageSpec()
	volumeMounts := []corev1.VolumeMount{
		{
			Name:      "collector-config",
			MountPath: "/etc/otel",
		},
	}
	volumes := []corev1.Volume{
		{
			Name: "collector-config",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: r.configMapName(),
					},
				},
			},
		},
	}
	if r.collector.Spec.LoggingConfig != nil {
		hostVolumeMounts, hostVolumes, err := r.hostLoggingVolumes()
		if err != nil {
			return resources.Error(nil, err)
		}
		volumeMounts = append(volumeMounts, hostVolumeMounts...)
		volumes = append(volumes, hostVolumes...)
	}
	ds := &appsv1.DaemonSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-collector-agent", r.collector.Name),
			Namespace: r.collector.Spec.SystemNamespace,
			Labels: map[string]string{
				resources.AppNameLabel: "collector-agent",
				resources.PartOfLabel:  "opni",
			},
		},
		Spec: appsv1.DaemonSetSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					resources.AppNameLabel: "collector-agent",
					resources.PartOfLabel:  "opni",
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						resources.AppNameLabel: "collector-agent",
						resources.PartOfLabel:  "opni",
					},
					Annotations: map[string]string{
						resources.OpniConfigHash: configHash,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name: "otel-collector",
							Command: []string{
								"/otelcol-contrib",
								fmt.Sprintf("--config=/etc/otel/%s", mainKey),
							},
							Image:           *imageSpec.Image,
							ImagePullPolicy: *imageSpec.ImagePullPolicy,
							Env: []corev1.EnvVar{
								{
									Name: "NODE_NAME",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "spec.nodeName",
										},
									},
								},
								{
									Name: "HOST_IP",
									ValueFrom: &corev1.EnvVarSource{
										FieldRef: &corev1.ObjectFieldSelector{
											FieldPath: "status.hostIP",
										},
									},
								},
							},
							VolumeMounts: volumeMounts,
						},
					},
					ImagePullSecrets: imageSpec.ImagePullSecrets,
					Volumes:          volumes,
				},
			},
		},
	}
	ctrl.SetControllerReference(r.collector, ds, r.client.Scheme())

	return resources.Present(ds)
}

func (r *Reconciler) imageSpec() opnimeta.ImageSpec {
	return opnimeta.ImageResolver{
		Version:       collectorVersion,
		ImageName:     collectorImage,
		DefaultRepo:   collectorImageRepo,
		ImageOverride: &r.collector.Spec.ImageSpec,
	}.Resolve()
}
