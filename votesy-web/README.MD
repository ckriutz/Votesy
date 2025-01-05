# Votesy Web

This is the main web page that shows the questions, and allows the user to vote.

It is a Python Flask app (becasue why not?) and is fairly simple for now.
The goal is to have a user see the question, and the voting options, and then click one.
The one that is clicked sends a message to Azure Queue storage, where it is processed elsewhere.
When the user votes, we will redirect to the page that shows the results!

The results page used to be a Java app, but that added more complexity than we liked.

In order for this to work we need the following env variables:

```
API_URL=
STORAGE_CONNECTION_STRING=
```