# Votesy API

The Votesy API is designed to provide CRUD operations on the list of questions and answers, and the number of Votes that are in the database, including which is the current question.

This API is written in GO. Why? I dont know, I thought writing it in Go would be fun. It was, but there were times when I thought that writing it in dotnet would have been easier. Time will tellI suppose.

The data **was** using Redis to store all the data, but that has since been removed, and replaced with Azure Table Storage. Ideally, in the future, the storage would be a database, and data cached in Redis, but that time is not now. Honestly, I think the performance will be fine. If it proves me wrong, I'll try another solution.

To get this to work, you will need to add a ENV variable for table storage:

```
STORAGE_CONNECTION_STRING
```