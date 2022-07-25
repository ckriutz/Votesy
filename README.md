# Votesy!

Votesy is a simple app designed to:
* Demonstrate Microservices
* Demonstrate Docker Compose
* Demonstrate Kubernetes


🚀 This was inspired by [example-voting-app](https://github.com/dockersamples/example-voting-app). I used it as a reference and borrowed a lot of code from it, however it's near ancient now, and I wanted to add my own flair.

## Background

If you wanted to make this app, you could certainly do it many, many different ways. I wanted to show how all the different services work together to form a single solution, and then how they can be *orchestrated* using docker-compose and Kubernetes.

I also wanted to use as many different languages and frameworks I could without losing my mind, so I chose C#, Java, Python, and Node. Java was the worst to work with by far. Still, they all work great!

## The Different Services

This application is broken out into 4 different projects:
* votesy-api - A **node.js** REST api that will return the current questions and answers. It will get the current questions and answers from Azure Table Storage.
* votesy-web - A **Python Flask** project that will call into **votesy-api** to get the current question and answers and display them, allowing the user to vote. When the user votes, it will send a message into Azure Queue Storage.
* Votesy.Service - A **dotnet core** service that will listen for messages in the queue storage, and then update Azure Table storage with the results.
* Votesy.Results = A **Java Spring Boot** project that will call into Azure Table storage to get the results and display them.

Additionally, there is a **Infrastructure** folder that has all the scripts required to set up the Azure resources, as well as the scripts to deploy.

## votesy-api
Will look for the questions table in Azure table storage, and return the questions and answers that are "isActive". This is used by **votsey-web** to display the options, and also by **votsey.results** to get the current question and vote count.




## FAQ's
