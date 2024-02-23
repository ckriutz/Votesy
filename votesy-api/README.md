# Votesy API

The Votesy API is designed to provide CRUD operations on the list of questions and answers that are in the database, including which is the current question.

This API is written in GO.

Right now, I've made the decision to store everything in Redis. This should be more than fine for several hundred questions, but after that, the payload is going to be too large to send all the questions back to the user, so we're going to have to update this.

The golang Redis tool doesn't allow for the RedisStack commands which would make this so much easier. That would allow me to query Redis in ways that would make it work beyond several hundred questions. I can also investigate writing those commands myself, and that might work.

If those don't work, I'll need to strongly consider writing questions and answers to a standard database. I'm not thrilled with the Redis database tools either, so it's a mixed bag. This should work for now.