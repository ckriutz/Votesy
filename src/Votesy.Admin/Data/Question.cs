using Microsoft.Azure.Cosmos.Table;
using System;

namespace Votesy.Admin.Data
{
    public class Question : TableEntity
    {
        public string Id { get; set; }
        public string Text { get; set; }
        public string Answer1Id { get; set; }
        public string Answer1Text { get; set; }
        public string Answer2Id { get; set; }
        public string Answer2Text { get; set; }
        public string Answer3Id { get; set; }
        public string Answer3Text { get; set; }
        public string Answer4Id { get; set; }
        public string Answer4Text { get; set; }
        public bool IsCurrent { get; set; }
        public bool IsUsed { get; set; }
        public DateTime CreatedDate { get; set; }
    }
}