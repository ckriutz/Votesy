using System.Net.Sockets;
using Microsoft.Extensions.Diagnostics.HealthChecks;

public class HealthChecker : IHealthCheck
{
    private readonly ILogger<HealthChecker> _logger;
    private QueueService _queueService;
    private string _tableName;
    public HealthChecker(ILogger<HealthChecker> logger, QueueService queueService)
    {
        _logger = logger;
        _queueService = queueService;
        _tableName = Environment.GetEnvironmentVariable("TableName");
    }

    public async Task<HealthCheckResult> CheckHealthAsync(HealthCheckContext context, CancellationToken cancellationToken = default)
    {

        // check to see if we can query the table storage, and if we can't, return unhealthy.
        var status = await _queueService.GetStatusAsync(_tableName);
        if(status == true)
        {
            // Things look good!
            var healthyResult = new HealthCheckResult( status: HealthStatus.Healthy, description: "Access to table storage is okay.");
            _logger.LogInformation(healthyResult.Status.ToString());
            return healthyResult;
        }
        else
        {
            // Things don't look so good.
            var unhealthyResult = new HealthCheckResult( status: HealthStatus.Unhealthy, description: "Unable to access table storage.");
            _logger.LogError(unhealthyResult.Status.ToString());
            return unhealthyResult;
        }

        
    }


}