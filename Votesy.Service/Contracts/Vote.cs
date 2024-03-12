namespace Contracts;
public record Vote
{
    public string? QuestionId { get; set; }
    public string AnswerId { get; set; }
}