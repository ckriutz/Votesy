using Azure;
using Azure.Data.Tables;

public class VoteEntity : ITableEntity
{
    public string PartitionKey { get; set; }
    public string RowKey { get; set; }
    public DateTimeOffset? Timestamp { get; set; }
    public ETag ETag { get; set; }
    public int VoteCount { get; set; }
}
