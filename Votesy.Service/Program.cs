using System;
using Microsoft.Extensions.Hosting;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.DependencyInjection;

public class Program
{
    public static void Main(string[] args)
    {
        CreateHostBuilder(args).Build().Run();
    }

    public static IHostBuilder CreateHostBuilder(string[] args) => Host.CreateDefaultBuilder(args)
    .ConfigureServices((hostContext, services) =>
    {
        // Now we can add this service.
        services.AddHostedService<QueueService>();
    })
    .ConfigureLogging((hostContext, configLogging) =>
    {
        // Add general console logging.
        configLogging.AddConfiguration(hostContext.Configuration.GetSection("Logging"));
        configLogging.AddConsole();
    });
    // See https://aka.ms/new-console-template for more information
}

