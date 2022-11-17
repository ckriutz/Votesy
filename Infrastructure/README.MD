# Deploying Votesy!
The fun part of this app is deploying it! We have options to deploy the application into the following infrastructure:

* Locally using docker compose.
* Azure Kubernetes Service (AKS)
* Azure Container Apps

### Step 1, lets connect to our Azure Subscription.
az login
az account set --subscription "subsciption name"

### Step 2, regardless of which method we choose, we need to create a storage account in the resource group.
Run the following command:
```
az deployment group create --resource-group votesy-rg --template-file storage.bicep
```

When the command prompts for a location, type in the preferred location.

```
Please provide string value for 'location' (? for help):
```

## Instructions for Docker Compose

## Instructions for Azure Container Apps

## Instructions for Azure Kubernetes Service