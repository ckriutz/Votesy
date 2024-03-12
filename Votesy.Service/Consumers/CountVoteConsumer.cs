using MassTransit;
using NRedisStack;
using NRedisStack.RedisStackCommands;
using StackExchange.Redis;

namespace Consumers
{
    public class CountVoteConsumer : IConsumer<Contracts.Vote>
    {
        // get the environment variable for the Redis host
        static string redisConnection = Environment.GetEnvironmentVariable("REDIS_HOST") ?? "192.168.0.239";

        static ConnectionMultiplexer redis = ConnectionMultiplexer.Connect(redisConnection);
        static IDatabase db = redis.GetDatabase();

        public async Task Consume(ConsumeContext<Contracts.Vote> context)
        {
            // This is where we count the vote.
            Console.WriteLine($"Going to increment the vote count for answer votes:{context.Message.AnswerId}.");
            // We need to increment the vote count for the answer in Redis.
            var result = db.StringIncrement($"votes:{context.Message.AnswerId}");
            Console.WriteLine($"Vote count for answer {context.Message.AnswerId} is now {result}.");
        }
    }
}