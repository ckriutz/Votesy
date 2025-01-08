using System;
using System.Collections.Generic;
using System.Net.Http;
using System.Net.Http.Json;
using System.Threading.Tasks;

namespace Votesy.Admin.Data
{
    public class QuestionService
    {
        private readonly HttpClient _httpClient;

        public QuestionService(HttpClient httpClient)
        {
            _httpClient = httpClient;
        }

        public async Task<List<Question>> GetQuestionsAsync()
        {
            try
            {
                var response = await _httpClient.GetFromJsonAsync<List<Question>>("https://votesy-api.nicecliff-3e424aeb.eastus2.azurecontainerapps.io/questions/");
                if (response == null)
                {
                    return new List<Question>();
                }
                // Order the response by the most recent questions
                response.Sort((q1, q2) => q2.CreatedDate.CompareTo(q1.CreatedDate));
                return response;
            }
            catch (Exception ex)
            {
                // Handle exceptions as needed
                Console.WriteLine($"Error fetching questions: {ex.Message}");
                return new List<Question>();
            }
        }

        public async Task AddQuestionAsync(Question question)
        {
            var response = await _httpClient.PostAsJsonAsync<Question>("https://votesy-api.nicecliff-3e424aeb.eastus2.azurecontainerapps.io/question", question);
            response.EnsureSuccessStatusCode();
        }

        public async Task DeleteQuestionAsync(string id)
        {
            var response = await _httpClient.DeleteAsync($"https://votesy-api.nicecliff-3e424aeb.eastus2.azurecontainerapps.io/question/questions/{id}");
            response.EnsureSuccessStatusCode();
        }
    }
}