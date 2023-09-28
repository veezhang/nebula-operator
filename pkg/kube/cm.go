/*
Copyright 2021 Vesoft Inc.

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

package kube

import (
	"context"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	apiequality "k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/client-go/util/retry"
	"k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type ConfigMap interface {
	CreateOrUpdateConfigMap(cm *corev1.ConfigMap) error
	GetConfigMap(namespace, cmName string) (*corev1.ConfigMap, error)
	DeleteConfigMap(namespace, cmName string) error
}

type cmClient struct {
	kubecli client.Client
}

func NewConfigMap(kubecli client.Client) ConfigMap {
	return &cmClient{kubecli: kubecli}
}

func (c *cmClient) CreateOrUpdateConfigMap(cm *corev1.ConfigMap) error {
	if err := c.kubecli.Create(context.TODO(), cm); err != nil {
		if !apierrors.IsAlreadyExists(err) {
			return err
		}
		merge := func(existing, desired *corev1.ConfigMap) error {
			existing.Data = desired.Data
			existing.Labels = desired.Labels
			for k, v := range desired.Annotations {
				existing.Annotations[k] = v
			}
			return nil
		}
		key := client.ObjectKeyFromObject(cm)
		existing, err := c.getConfigMap(key)
		if err != nil {
			return err
		}
		mutated := existing.DeepCopy()
		if err := merge(mutated, cm); err != nil {
			return err
		}
		if !apiequality.Semantic.DeepEqual(existing, mutated) {
			if err := c.updateConfigMap(mutated); err != nil {
				return err
			}
		}
		return nil
	}
	return nil
}

func (c *cmClient) GetConfigMap(namespace, cmName string) (*corev1.ConfigMap, error) {
	return c.getConfigMap(client.ObjectKey{Namespace: namespace, Name: cmName})
}

func (c *cmClient) updateConfigMap(cm *corev1.ConfigMap) error {
	ns := cm.GetNamespace()
	cmName := cm.GetName()
	cmData := cm.Data

	return retry.RetryOnConflict(retry.DefaultBackoff, func() error {
		if updated, err := c.GetConfigMap(ns, cmName); err == nil {
			cm = updated.DeepCopy()
			cm.Data = cmData
		} else {
			utilruntime.HandleError(fmt.Errorf("get configMap [%s/%s] failed: %v", ns, cmName, err))
			return err
		}

		updateErr := c.kubecli.Update(context.TODO(), cm)
		if updateErr == nil {
			klog.Infof("configMap [%s/%s] updated successfully", ns, cmName)
			return nil
		}
		return updateErr
	})
}

func (c *cmClient) getConfigMap(objKey client.ObjectKey) (*corev1.ConfigMap, error) {
	configMap := &corev1.ConfigMap{}
	err := c.kubecli.Get(context.TODO(), objKey, configMap)
	if err != nil {
		klog.V(4).ErrorS(err, "failed to get configmap", "namespace", objKey.Namespace, "name", objKey.Name)
		return nil, err
	}
	return configMap, err
}

func (c *cmClient) DeleteConfigMap(namespace, cmName string) error {
	cm, err := c.getConfigMap(client.ObjectKey{Namespace: namespace, Name: cmName})
	if err != nil {
		return err
	}
	return c.kubecli.Delete(context.TODO(), cm)
}
