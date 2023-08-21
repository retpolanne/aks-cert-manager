provider "helm" {
  kubernetes {
    config_path = "./kubeconfig"
  }
}

provider "azurerm" {
  features {}
  skip_provider_registration = true
}
