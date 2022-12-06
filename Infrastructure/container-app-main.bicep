param location string = 'eastus2'
param envName string = 'votesy'

param registry string
param registryUsername string

@secure()
param registryPassword string

param storageName string

@secure()
param storageConnectionString string

@secure()
param storageKey string

module law './modules/law.bicep' = {
    name: 'votesy-log-analytics-workspace'
    params: {
      location: location
      name: 'law-${envName}'
    }
}

module containerAppEnvironment './modules/app-environment.bicep' = {
  name: 'container-app-environment'
  params: {
    name: envName
    location: location
    lawClientId:law.outputs.clientId
    lawClientSecret: law.outputs.clientSecret
  }
}

module vosesyapi './modules/containerapp.bicep' = {
  name: 'votesy-api'
  params: {
    name: 'votesy-api'
    location: location
    containerAppEnvironmentId: containerAppEnvironment.outputs.id
    containerImage: '${registry}/votesy-api:latest'
    enableIngress: true
    // We do not need external ingress on this since it doens't need to be accessed outside of the App.
    useExternalIngress: false
    containerPort: 10000
    envVars: [
      {
        name: 'AZURE_CONNECTION_STRING'
        value: storageConnectionString
      }
      {
        name: 'KEY'
        value: storageKey
      }
    ]
    probes: [
      {
        type: 'liveness'
        initialDelaySeconds: 15
        periodSeconds: 10
        failureThreshold: 3
        timeoutSeconds: 1
        httpGet: {
          port: 10000
          path: '/health/liveness'
        }
      }
      {
        type: 'startup'
        initialDelaySeconds: 15
        periodSeconds: 10
        failureThreshold: 3
        timeoutSeconds: 2
        httpGet: {
          port: 10000
          path: '/health/startup'
        }
      }
      {
        type: 'readiness'
        initialDelaySeconds: 15
        periodSeconds: 10
        failureThreshold: 3
        timeoutSeconds: 2
        httpGet: {
          port: 10000
          path: '/health/readiness'
        }
      }
    ]
    registry: registry
    registryUsername: registryUsername
    registryPassword: registryPassword
  }
}

module votseyservice './modules/containerapp.bicep' = {
  name: 'votsey-service'
  params: {
    name: 'votesy-service'
    location: location
    containerAppEnvironmentId: containerAppEnvironment.outputs.id
    containerImage: '${registry}/votesy-service:latest'
    // We do not need external ingress on this since it doens't need to be accessed outside of the App.
    useExternalIngress: false
    enableIngress: false
    containerPort: 7999 // Don't actually need this.
    envVars: [
      {
        name: 'QueueConnectionString'
        value: storageConnectionString
      }
      {
        name: 'Key'
        value: storageKey
      }
      {
        name: 'TableName'
        value: 'votes'
      }
      {
        name: 'AccountName'
        value: storageName
      }
    ]
    registry: registry
    registryUsername: registryUsername
    registryPassword: registryPassword
  }
}

module votseyresults 'modules/containerapp.bicep' = {
  name: 'votesy-results'
  params: {
    name: 'votesy-results'
    location: location
    containerAppEnvironmentId: containerAppEnvironment.outputs.id
    containerImage: '${registry}/votesy-results:latest'
    enableIngress: true
    containerPort: 8080
    useExternalIngress: true
    envVars: [
      {
        name: 'connectionString'
        value: storageConnectionString
      }
      {
        name: 'voteUrl'
        value: 'https://votesy-web.${containerAppEnvironment.outputs.id}.${location}.azurecontainerapps.io'
      }
    ]
    registry: registry
    registryUsername: registryUsername
    registryPassword: registryPassword
  }
}

module votesyweb './modules/containerapp.bicep' = {
  name: 'votesy-web'
  params: {
    name: 'votesy-web'
    location: location
    containerAppEnvironmentId: containerAppEnvironment.outputs.id
    containerImage: '${registry}/votesy-web:latest'
    enableIngress: true
    containerPort: 5002
    useExternalIngress: true
    envVars: [
        {
          name: 'apiUrl'
          value: 'https://${vosesyapi.outputs.fqdn}'
        }
        {
          name: 'resultsURL'
          value: 'https://${votseyresults.outputs.fqdn}'
        }
        {
          name: 'AZURE_STORAGE_CONNECTION_STRING'
          value: storageConnectionString
        }
    ]
    probes: [
      {
        type: 'liveness'
        initialDelaySeconds: 15
        periodSeconds: 10
        failureThreshold: 3
        timeoutSeconds: 1
        httpGet: {
          port: 10000
          path: '/health/liveness'
        }
      }
      {
        type: 'startup'
        initialDelaySeconds: 15
        periodSeconds: 10
        failureThreshold: 3
        timeoutSeconds: 2
        httpGet: {
          port: 10000
          path: '/health/startup'
        }
      }
      {
        type: 'readiness'
        initialDelaySeconds: 15
        periodSeconds: 10
        failureThreshold: 3
        timeoutSeconds: 2
        httpGet: {
          port: 10000
          path: '/health/readiness'
        }
      }
    ]
    registry: registry
    registryUsername: registryUsername
    registryPassword: registryPassword
  }
}

output appEnvId string = containerAppEnvironment.outputs.id
output webfqdn string = votesyweb.outputs.fqdn
output resultsfqdn string = votseyresults.outputs.fqdn
