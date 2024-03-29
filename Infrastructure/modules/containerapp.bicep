// general Azure Container App settings
param location string
param name string
param containerAppEnvironmentId string

// Container Image ref
param containerImage string

// Networking
param enableIngress bool
param useExternalIngress bool
param containerPort int

param registry string
param registryUsername string

@secure()
param registryPassword string

param envVars array = []

param probes array = []

resource containerApp 'Microsoft.App/containerApps@2022-03-01' = {
  name: name
  location: location
  properties: {
    managedEnvironmentId: containerAppEnvironmentId
    configuration: {
      secrets: [
        {
          name: 'registrypassword'
          value: registryPassword
        }
      ]      
      registries: [
        {
          server: registry
          username: registryUsername
          passwordSecretRef: 'registrypassword'
        }
      ]
      ingress: enableIngress ? {
        external: useExternalIngress
        targetPort: containerPort
        transport: 'auto'
        traffic: [
          {
            latestRevision: true
            weight: 100
          }
        ]
      } : null
    }
    template: {
      containers: [
        {
          image: containerImage
          name: name
          env: envVars
          probes: probes
        }
      ]
      scale: {
        minReplicas: 0
      }
    }
  }
}

output fqdn string = enableIngress ? containerApp.properties.configuration.ingress.fqdn : 'Ingress not enabled'
