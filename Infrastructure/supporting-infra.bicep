// Some general variables!
@description('Provide a location for the registry.')
param location string = 'eastus2'

// Setting target scope
targetScope = 'subscription'

// First, create the resource group!
resource rg 'Microsoft.Resources/resourceGroups@2021-01-01' = {
  name: 'rg-votesy'
  location: location
}

// Then, make sure the storage account is there.
module storage './modules/storage.bicep' = {
  name: 'storageDeployment'
  params: {
    location: location
  }
  scope: rg
}

// Finally, create the container registry for our services.
module acr './modules/container-registry.bicep' = {
  name: 'acrDeployment'
  scope: rg
  params: {
    location: location
  }
}
