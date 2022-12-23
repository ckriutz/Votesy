using Azure.Storage.Queues;
using Azure.Storage.Queues.Models;
using System.Text.Json;
using Azure.Data.Tables;

public class QueueService
{
    private readonly ILogger<QueueService> _logger;

    // I don't like them being here, but, lets just grab these once.
    // The connection string here is the connection string to the Queue Storage.
    string _connectionString;
    
    // The Storage Account Key is the key used to connect to table storage.
    string _storageAccountKey;

    // The table name is the name of the table we will be updating.
    string _tableName;

    // The account name from Azure.
    string _accountName;

    string _storageUri;


    public QueueService(ILogger<QueueService> logger)
    {
        _logger = logger;
        _connectionString = Environment.GetEnvironmentVariable("QueueConnectionString");
        _storageAccountKey = Environment.GetEnvironmentVariable("Key");
        _tableName = Environment.GetEnvironmentVariable("TableName");
        _accountName = Environment.GetEnvironmentVariable("AccountName");
        _storageUri = $"https://{_accountName}.table.core.windows.net/{_tableName}";
    }

    public async Task<Boolean> GetStatusAsync(string queueName)
    {
        QueueClient queueClient = new QueueClient(_connectionString, queueName);
        return await queueClient.ExistsAsync();
    }

    //-------------------------------------------------
    // Create the queue service client
    //-------------------------------------------------
    public async Task GetMessagesFromQueue(string queueName)
    {
        // Instantiate a QueueClient which will be used to create and manipulate the queue
        QueueClient queueClient = new QueueClient(_connectionString, queueName);

        // Get messages from the queue
        try
        {
            QueueMessage[] messages = await queueClient.ReceiveMessagesAsync(maxMessages: 10);
            if(messages.Count() == 0)
            {
                _logger.LogInformation("No votes to record at this time.");
            }
            else
            {
                foreach(QueueMessage m in messages)
                {
                    // "Process" the message
                    Console.WriteLine($"Message: {m.MessageText}");
                    Vote vote = JsonSerializer.Deserialize<Vote>(m.MessageText);

                    await UpdateVoteCount(vote);

                    // Let the service know we're finished with
                    // the message and it can be safely deleted.
                    await queueClient.DeleteMessageAsync(m.MessageId, m.PopReceipt);
                }
            }
        }
        catch(System.Net.Http.HttpRequestException NoHostIsKnownException)
        {
            _logger.LogError(NoHostIsKnownException.ToString());
        }
        catch(Exception ex)
        {
            _logger.LogError(ex.ToString());
        } 
        
        
    }

    public async Task UpdateVoteCount(Vote vote)
    {
        _logger.LogInformation("Updating Vote Count.");
        
        _logger.LogInformation(_storageUri);

        var tableClient = new TableClient(new Uri(_storageUri), _tableName, new TableSharedKeyCredential(_accountName, _storageAccountKey));

        // Create the table in the service.
        await tableClient.CreateIfNotExistsAsync();

        //Azure.Pageable<Vote> queryResultsLINQ = tableClient.Query<Vote>(ent => ent.PartitionKey == "Question1" && ent.RowKey.Equals(questionId.ToString()));

        try
        {
            VoteEntity qEntity = await tableClient.GetEntityAsync<VoteEntity>(vote.questionId, vote.answerId.ToString());
            qEntity.VoteCount++;
            await tableClient.UpdateEntityAsync(qEntity, qEntity.ETag);

        }
        catch (Azure.RequestFailedException rfe)
        {
            // oh shit it's null.
            _logger.LogError(rfe.ToString());
            
            // This isn't maybe the best way to handle this, but I'm single person hackathoning this, so it works for now.
            VoteEntity qEntity = new VoteEntity();
            qEntity.VoteCount = 0;
            qEntity.PartitionKey = vote.questionId;
            qEntity.RowKey = vote.answerId.ToString();
            await tableClient.AddEntityAsync(qEntity);
        }
        catch (Exception ex)
        {
            Console.WriteLine(ex.Message);
        }

    }
}