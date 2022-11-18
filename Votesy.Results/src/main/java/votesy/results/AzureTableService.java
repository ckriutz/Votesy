package votesy.results;

import java.util.ArrayList;

// Include the following imports to use table APIs
import com.azure.data.tables.TableClient;
import com.azure.data.tables.TableClientBuilder;
import com.azure.data.tables.TableServiceClient;
import com.azure.data.tables.TableServiceClientBuilder;
import com.azure.data.tables.models.ListEntitiesOptions;

public class AzureTableService {
    // Retrieve storage account from connection-string.
    String connectionStringEnv = System.getenv("connectionString");

    public void listTables() {
        try {
            // Create a TableServiceClient with a connection string.
            TableServiceClient tableServiceClient = new TableServiceClientBuilder()
                .connectionString(connectionStringEnv)
                .buildClient();

            // Loop through a collection of table names.
            tableServiceClient.listTables().forEach(tableItem -> 
                System.out.println(tableItem.getName())
            );
        }
        catch (Exception e) {
            // Output the stack trace.
            e.printStackTrace();
        }
    }

    public Question getCurrentQuestion() {
        final String tableName = "questions";

        // Seriously there should just be one of them.
        ArrayList<Question> questions = new ArrayList<Question>();

        try {
            // Create a TableClient with a connection string and a table name.
            TableClient tableClient = new TableClientBuilder().connectionString(connectionStringEnv).tableName(tableName).buildClient();

            // Create a filter condition where the partition key is "Sales".
            ListEntitiesOptions options = new ListEntitiesOptions().setFilter("PartitionKey eq 'Questions' and isCurrent eq true");
            tableClient.listEntities(options, null, null).forEach(tableEntity -> {
                Question question = new Question();
                question.PartitionKey = tableEntity.getPartitionKey();
                question.RowKey = tableEntity.getRowKey();
                question.text = tableEntity.getProperty("text").toString();
                question.answer1Id = tableEntity.getProperty("answer1Id").toString();
                question.answer1Text = tableEntity.getProperty("answer1Text").toString();
                question.answer2Id = tableEntity.getProperty("answer2Id").toString();
                question.answer2Text = tableEntity.getProperty("answer2Text").toString();

                questions.add(question);
            });
        }
        catch (Exception e) {
            // Output the stack trace.
            e.printStackTrace();
        }

        return questions.get(0);
    }

    public ArrayList<Vote> getVotesForQuestion(Question question) {
        final String tableName = "votes";

        // Seriously there should just be one of them.
        ArrayList<Vote> votes = new ArrayList<Vote>();

        try {
            // Create a TableClient with a connection string and a table name.
            TableClient tableClient = new TableClientBuilder().connectionString(connectionStringEnv).tableName(tableName).buildClient();
            ListEntitiesOptions options = new ListEntitiesOptions().setFilter("PartitionKey eq '" + question.RowKey +"'");
            tableClient.listEntities(options, null, null).forEach(tableEntity -> {
                System.out.println(tableEntity.getPartitionKey() +
                    " " + tableEntity.getRowKey() +
                    "\t" + tableEntity.getProperty("VoteCount"));
                Vote v = new Vote();
                v.PartitionKey = tableEntity.getPartitionKey();
                v.RowKey = tableEntity.getRowKey();
                v.VoteCount = Integer.parseInt(tableEntity.getProperty("VoteCount").toString());
                votes.add(v);
            });
        }
        catch (Exception e) {
            // Output the stack trace.
            e.printStackTrace();
        }

        return votes;
    }
}
