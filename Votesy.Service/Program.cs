using System;
using Microsoft.Extensions.Hosting;
using System.Threading;
using System.Threading.Tasks;
using Microsoft.Extensions.Logging;
using Microsoft.Extensions.DependencyInjection;
using MassTransit;


var builder = WebApplication.CreateBuilder(args);

// Need to get some environment variables
var rabbitMqHost = Environment.GetEnvironmentVariable("RABBITMQ_HOST") ?? "192.168.0.239";

// Add services to the container.
builder.Services.AddControllersWithViews();

builder.Services.AddMassTransit(x =>
{
    x.AddConsumer<Consumers.CountVoteConsumer>();

    x.UsingRabbitMq((context,cfg) =>
    {
        cfg.Host(rabbitMqHost, "/", h => {
            h.Username("guest");
            h.Password("guest");
        });

        cfg.ReceiveEndpoint("vote-queue", endpointConfigurator =>
        {
            endpointConfigurator.ConfigureConsumeTopology = false;
            endpointConfigurator.UseRawJsonDeserializer();
            endpointConfigurator.DefaultContentType = new System.Net.Mime.ContentType("application/json");
            endpointConfigurator.ConfigureConsumer<Consumers.CountVoteConsumer>(context);
        });

        cfg.ConfigureEndpoints(context);
    });
});

builder.Services.AddHealthChecks().AddCheck<HealthChecker>(name: "HealthChecker", failureStatus: Microsoft.Extensions.Diagnostics.HealthChecks.HealthStatus.Degraded);

var app = builder.Build();

app.UseHttpsRedirection();
app.UseStaticFiles();

app.UseRouting();

app.UseAuthorization();

app.UseHealthChecks("/health");

app.Run();