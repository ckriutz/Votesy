# Votesy!

Votesy is a simple app designed to:
* Demonstrate Microservices
* Demonstrate Docker Compose
* Demonstrate Kubernetes


🚀 This was inspired by [example-voting-app](https://github.com/dockersamples/example-voting-app). I used it as a reference and borrowed a lot of code from it, however it's near ancient now, and I wanted to add my own flair.

## Background

If you wanted to make this app, you could certainly do it many, many different ways. I wanted to show how all the different services work together to form a single solution, and then how they can be *orchestrated* using docker-compose and Kubernetes.

I also wanted to use as many different languages and frameworks I could without losing my mind, so I chose C#, Java, Python, Go, and Node. Java was the worst to work with by far. Still, they all work great!

## The Different Services

This application is broken out into 4 different projects:
* votesy-api - A **Go** REST api that will return the current questions and answers. It will get the current questions and answers from Redis, and also allows for most CRUD operations as well.
* votesy-web - A **Python Flask** project that will call into **votesy-api** to get the current question and answers and display them, allowing the user to vote. When the user votes, it will send a message into RabbitMQ.
* Votesy.Service - A **dotnet core** service that will listen for messages in the queue storage, and then update Azure Table storage with the results.
* Votesy.Results = A **Java Spring Boot** project that will call into Azure Table storage to get the results and display them.

Additionally, there is a **Infrastructure** folder that has all the scripts required to set up the Azure resources, as well as the scripts to deploy.

## votesy-api
Will look for the questions table in Azure table storage, and return the questions and answers that are "isActive". This is used by **votsey-web** to display the options, and also by **votsey.results** to get the current question and vote count. On initial run, it will check to see if there are any questions in the database, and if not, it will add one. Makes it easy for initial deployments!

The website will be exposed on port 5002, whoch is hardcoded (for now), so the orchistration needs to match that as well.

## Votesy.Service
A dotnet 6 Service. It's designed to check the queue storage every 5 seconds, process the message, and update the number of votes in the process. None of the other services connect to this service. It's a asp.net app, sort of, because we enables health probes, which for now, require us to be able to resolve an endpoint. That Endpoint is on 7999.

## Votesy.Results
A Java Spring Boot Application. It connects to the **votesy-api** to get the current question and votes in order to display that information. I originally started building this out as a Blazor application, but wanted yet another framework to build out. I'd never done Java Spring Boot before, so I wanted an oppertunity. It was a pain, but works well.

The endpoint is set to port 8080, and that's hardcoded (for now), so the orchistration needs to match that as well.

## Efforts to add health checks to the services:
✅ votesy-api
🚫 votesy-results
✅ votesy-service
✅ votesy-web

## FAQ's
