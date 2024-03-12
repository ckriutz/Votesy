using System.Net.Sockets;
using Microsoft.Extensions.Diagnostics.HealthChecks;

public class HealthChecker : IHealthCheck
{
    private readonly ILogger<HealthChecker> _logger;
    public HealthChecker(ILogger<HealthChecker> logger)
    {
        _logger = logger;
    }

    public async Task<HealthCheckResult> CheckHealthAsync(HealthCheckContext context, CancellationToken cancellationToken = default)
    {
        // Check to see if we can connect to Redis, and RabbitMQ, and if we can't, return unhealthy.

        // check to see if we can query the table storage, and if we can't, return unhealthy.
        //var status = await _queueService.GetStatusAsync(_tableName);
        //if(status == true)
        //{
            // Things look good!
            //var healthyResult = new HealthCheckResult( status: HealthStatus.Healthy, description: "Access to table storage is okay.");
            //_logger.LogInformation(healthyResult.Status.ToString());
            //return healthyResult;
        //}
        //else
        //{
            // Things don't look so good.
            //var unhealthyResult = new HealthCheckResult( status: HealthStatus.Unhealthy, description: "Unable to access table storage.");
            //_logger.LogError(unhealthyResult.Status.ToString());
            //return unhealthyResult;
        //}
        return HealthCheckResult.Healthy("This is a test");
        
    }


}