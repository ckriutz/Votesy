
// This is the main service that does all the work to manage Votes.
// At the moment, it checks Azure Queue Storage, but that can always change.
public class VoteService : BackgroundService
{
    private readonly ILogger<VoteService> _logger;
    private readonly QueueService _queueService;
    private readonly string _tableName;

    public VoteService(ILogger<VoteService> logger, QueueService queueService)
    {
        _logger = logger;
        _queueService = queueService;
        _tableName = Environment.GetEnvironmentVariable("TableName");
    }
    protected override async Task ExecuteAsync(CancellationToken stoppingToken)
    {
        while (!stoppingToken.IsCancellationRequested)
        {
            _logger.LogInformation("Checking for Messages...");
            await _queueService.GetMessagesFromQueue(_tableName);

            await Task.Delay(TimeSpan.FromSeconds(5), stoppingToken);
        }
    }
}