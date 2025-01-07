const { app } = require('@azure/functions');
const { TableClient, AzureNamedKeyCredential } = require('@azure/data-tables');

// Define your Table Storage connection details
const accountName = process.env.StorageAccountName;
const accountKey = process.env.StorageAccountKey;
const tableName = 'votes';
const credential = new AzureNamedKeyCredential(accountName, accountKey);
const tableClient = new TableClient(`https://${accountName}.table.core.windows.net`, tableName, credential);


app.storageQueue('storageQueueTrigger', {
    queueName: 'votes',
    connection: 'AzureStorageConnectionString',
    handler: async (queueItem, context) => {
        context.log('Queue item received:', queueItem);
        // The JSON message we get should look like this:
        // {"answerId": "3ief7"}

        const answerId = queueItem.answerId;
        context.log('AnswerId Voted For: ', answerId);

        try {
            // Retrieve the entity from Table Storage
            const entity = await tableClient.getEntity('votes', answerId);

            // Increment the voteCount
            entity.voteCount = (entity.voteCount || 0) + 1;
            //context.log('Updated number of votes: ', entity.voteCount);

            // Upsert the updated entity back into Table Storage
            await tableClient.upsertEntity(entity);
            context.log('Entity updated successfully');
        } catch (error) {
            if (error.statusCode === 404) {
                // Entity not found, create a new one
                context.log('Entity not found.');
            } else {
                context.log('Error updating entity:', error.message);
            }
        }
    }
});
