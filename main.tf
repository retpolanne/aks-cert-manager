resource "azurerm_resource_group" "terratest_certmanager" {
  name     = "terratest-certmanager"
  location = "East US"
}

resource "azurerm_kubernetes_cluster" "terratest_aks_certmanager" {
  location            = azurerm_resource_group.terratest_certmanager.location
  name                = "terratest-aks-certmanager"
  resource_group_name = azurerm_resource_group.terratest_certmanager.name
  dns_prefix          = "terratest-aks-certmanager"

  identity {
    type = "SystemAssigned"
  }

  default_node_pool {
    name       = "agentpool"
    vm_size    = "Standard_D2_v2"
    node_count = 3
  }
  network_profile {
    network_plugin    = "kubenet"
    load_balancer_sku = "standard"
  }
}

resource "local_file" "kubeconfig" {
	content  = azurerm_kubernetes_cluster.terratest_aks_certmanager.kube_config_raw
  filename = "./kubeconfig" 
}
