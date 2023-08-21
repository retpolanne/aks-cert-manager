package test

import (
	"flag"
	"os"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/azure"
	"github.com/gruntwork-io/terratest/modules/k8s"
	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var clusterName string = "terratest-aks-certmanager"
var rgName string = "terratest-certmanager"
var subscriptionId string

var destroy = flag.Bool("destroy", false, "In case you need to destroy the resources")

func getKubectlOptions(namespace string) *k8s.KubectlOptions {
	return &k8s.KubectlOptions{
		ContextName: clusterName,
		ConfigPath:  "../kubeconfig",
		Namespace:   namespace,
		Logger:      logger.Discard,
	}
}

func TestMain(t *testing.T) {
	subscriptionId = os.Getenv("ARM_SUBSCRIPTION_ID")
	if subscriptionId == "" {
		t.Fatal("Please provide the ARM_SUBSCRIPTION_ID variable")
	}
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../",
	})
	if *destroy {
		defer terraform.Destroy(t, terraformOptions)
	}
	terraform.InitAndApply(t, terraformOptions)
	TestAKSCreated(t)
	TestAKSHasCertManagerInstalled(t)
}

func TestAKSCreated(t *testing.T) {
	cluster, err := azure.GetManagedClusterE(t, rgName, clusterName, subscriptionId)
	if err != nil {
		t.Fatalf("Error getting cluster %s, err %s\n", clusterName, err)
	}
	require.NotNil(t, cluster)
	assert.Equal(t, clusterName, *cluster.Name)
}

func TestAKSHasCertManagerInstalled(t *testing.T) {
	k8sOptions := getKubectlOptions("cert-manager")
	deployment := k8s.GetDeployment(t, k8sOptions, "cert-manager")
	assert.True(t, k8s.IsDeploymentAvailable(deployment))
}

func TestAKSCertManagerHasCRDs(t *testing.T) {
	k8sOptions := getKubectlOptions("cert-manager")
	out, err := k8s.RunKubectlAndGetOutputE(t, k8sOptions, "get", "crd")
	if err != nil {
		t.Fatalf("Error running kubectl get crd %s\n", err)
	}
	crdExists := false
	crd := "cert-manager.io"
	for _, line := range strings.Split(out, "\n") {
		if strings.Contains(line, crd) {
			crdExists = true
		}
	}
	assert.Truef(t, crdExists, "expected %s to exist\n", crd)
}
