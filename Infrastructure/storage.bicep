targetScope = 'resourceGroup'

param location string

// Here, we create a storage account, which is key to our application!
resource createStorage 'Microsoft.Storage/storageAccounts@2021-06-01' = {
  name: 'votesysa'
  location: location
  sku: {
    name: 'Standard_LRS'
  }
  kind: 'StorageV2'

  // We need two tables.
  resource tableService 'tableServices' = {
    name: 'default'

    resource tableQ 'tables' = {
      name: 'questions'
    }

    resource tableV 'tables' = {
      name: 'votes'
    }
  }

  // ...and we need a queue!
  resource queueService 'queueServices' = {
    name: 'default'
    
    resource queue 'queues' = {
      name: 'votes'
    }
  }
}
