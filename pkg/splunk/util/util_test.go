package util

import (
	"context"
	"reflect"
	"testing"

	spltest "github.com/splunk/splunk-operator/pkg/splunk/test"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCreateResource(t *testing.T) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "test",
		},
		Data: map[string][]byte{"one": []byte("value1")},
	}

	c := spltest.NewMockClient()
	err := CreateResource(c, &secret)
	if err != nil {
		t.Errorf("CreateResource() returned %v; want nil", err)
	}
	c.CheckCalls(t, "TestCreateResource", map[string][]spltest.MockFuncCall{
		"Create": {
			{CTX: context.TODO(), Obj: &secret},
		},
	})
}

func TestUpdateResource(t *testing.T) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "test",
		},
		Data: map[string][]byte{"one": []byte("value1")},
	}

	c := spltest.NewMockClient()
	err := UpdateResource(c, &secret)
	if err != nil {
		t.Errorf("UpdateResource() returned %v; want nil", err)
	}
	c.CheckCalls(t, "TestUpdateResource", map[string][]spltest.MockFuncCall{
		"Update": {
			{CTX: context.TODO(), Obj: &secret},
		},
	})
}

func TestDeleteResource(t *testing.T) {
	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "secret",
			Namespace: "test",
		},
		Data: map[string][]byte{"one": []byte("value1")},
	}

	c := spltest.NewMockClient()
	err := DeleteResource(c, &secret)
	if err != nil {
		t.Errorf("DeleteResource() returned %v; want nil", err)
	}
	c.CheckCalls(t, "TestUpdateResource", map[string][]spltest.MockFuncCall{
		"Delete": {
			{CTX: context.TODO(), Obj: &secret},
		},
	})
}

func TestDeepCopyInto(t *testing.T) {
	cr := TestResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stack1",
			Namespace: "test",
		},
	}

	copy := cr.DeepCopy()

	if copy.Name != cr.Name {
		t.Errorf("TestResource copy.Name = %s; want %s", copy.Name, cr.Name)
	}

	if copy.Namespace != cr.Namespace {
		t.Errorf("TestResource copy.Namespace = %s; want %s", copy.Namespace, cr.Namespace)
	}
}

func TestDeepCopyObject(t *testing.T) {
	cr := TestResource{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "stack1",
			Namespace: "test",
		},
	}

	copy := cr.DeepCopyObject()

	if !reflect.DeepEqual(copy, &cr) {
		t.Errorf("TestResource \n got = %+v; \n want %+v \n", copy, cr)
	}
}

func TestPodExecCommand(t *testing.T) {
	// Create pod
	pod := &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "splunk-stack1-0",
			Namespace: "test",
			Labels: map[string]string{
				"controller-revision-hash": "v0",
			},
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					VolumeMounts: []corev1.VolumeMount{
						{
							MountPath: "/mnt/splunk-secrets",
							Name:      "mnt-splunk-secrets",
						},
					},
				},
			},
			Volumes: []corev1.Volume{
				{
					Name: "mnt-splunk-secrets",
					VolumeSource: corev1.VolumeSource{
						Secret: &corev1.SecretVolumeSource{
							SecretName: "test-secret",
						},
					},
				},
			},
		},
	}

	// Create client and add object
	c := spltest.NewMockClient()
	_, _, _ = PodExecCommand(c, "splunk-stack1-0", "test", []string{"/bin/sh"}, "ls -ltr", false, true)

	// Add object
	c.AddObject(pod)
	_, _, _ = PodExecCommand(c, "splunk-stack1-0", "test", []string{"/bin/sh"}, "ls -ltr", false, true)

	// Hit some error legs
	_, _, _ = PodExecCommand(c, "splunk-stack1-0", "test", []string{"/bin/sh"}, "ls -ltr", false, false)
}