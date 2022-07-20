# Votesy!

Votesy is a simple app designed to:
* Demonstrate Microservices
* Demonstrate Kubernetes

**🚀I used the following app as a reference, and reused a lot from it. However, it's near ancient now, so it's been slightly updated, and I've aded my own flair.**

```
https://github.com/dockersamples/example-voting-app
```

This application is broken out into 4 different projects:
* votesy-api - A **node.js** REST api that will return the current questions and answers. It will get the current questions and answers from Azure Table Storage.
* votesy-web - A **Python Flask** project that will call into **votesy-api** to get the current question and answers and display them, allowing the user to vote. When the user votes, it will send a message into Azure Queue Storage.
* Votesy.Service - A **dotnet core** service that will listen for messages in the queue storage, and then update Azure Table storage with the results.
* Votesy.Results = A **Java Spring Boot** project that will call into Azure Table storage to get the results and display them.

Additionally, there is a **Infrastructure** folder that has all the scripts required to set up the Azure resources.

## votesy-api
Will look for the questions table in Azure table storage, and return the questions and answers that are "isActive". This is used by **votsey-web** to display the options, and also by **votsey.results** to get the current question and vote count.




## FAQ's
