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
* votesy-api - A **Go** REST api that will return the current questions and answers. It will get the current questions and answers from Azure Table Storage, and also allows for most CRUD operations on the table storage.
* votesy-web - A **Python Flask** project that will call into **votesy-api** to get the current question and answers and display them, allowing the user to vote. When the user votes, it will send a message into Azure Queue Storage.
* Votesy.Service - A **dotnet core** service that will listen for messages in the queue storage, and then update Azure Table storage with the results.
* Votesy.Results = A **Java Spring Boot** project that will call into Azure Table storage to get the results and display them.

Additionally, there is a **Infrastructure** folder that has all the scripts required to set up the Azure resources, as well as the scripts to deploy.

## votesy-api
A REST API that will look for the questions table in Azure table storage, and return the questions and votes that are "isActive". This is used by **votsey-web** to display the questions and answers, and also by **Votsey.Results** to get the current question and vote count. Originally this was a node.js application, but after some thinking, I decided to re-write it in node.js. If you go back to [this commit](https://github.com/ckriutz/Votesy/tree/22dbf2bbd6457ad0d2f21b0d2dfe506ef428e101) you can see the node.js version in all it's glory. The Go version works well though!

The endpoint is set to port 10000, and that's hardcoded (for now), so the orchistration needs to match that as well.

## votesy-web
A Python Flask web application. It will call out to **votesy-api** and get the current question, which includes the answers for the options. When the user votes, it will indicate it, and then send out a message into Azure Queue Storage.

The website will be exposed on port 5002, whoch is hardcoded (for now), so the orchistration needs to match that as well.

## Votesy.Service
A dotnet 6 Service. It's designed to check the queue storage every 5 seconds, process the message, and update the number of votes in the process. None of the other services connect to this service.

## Votesy.Results
A Java Spring Boot Application. It connects to the **votesy-api** to get the current question and votes in order to display that information. I originally started building this out as a Blazor application, but wanted yet another framework to build out. I'd never done Java Spring Boot before, so I wanted an oppertunity. It was a pain, but works well.

The endpoint is set to port 8080, and that's hardcoded (for now), so the orchistration needs to match that as well.

## FAQ's
