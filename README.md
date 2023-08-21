# AKS Cert Manager 

Automated via Terraform and tested via Terratest.

## Setting up

For Mac, just run: 

```sh
brew bundle
az login
export ARM_SUBSCRIPTION_ID=`az account subscription list --query "[].subscriptionId" -o tsv`
```

You may need to change your DNS server to 8.8.8.8 due to [30549](https://github.com/hashicorp/terraform/issues/30549)

## Testing

```sh
go test -v ./...
```
