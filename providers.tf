provider "helm" {
  kubernetes {
    config_path = "./kubeconfig"
  }

  registry {
    url = "https://charts.jetstack.io"
  }
}

provider "azurerm" {
  features {}
  skip_provider_registration = true
}
