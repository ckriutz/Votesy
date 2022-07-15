var azure = require('azure-storage');
const express = require("express");

const app = express();

// To solve the cors issue
const cors=require('cors');

const tabeleName = process.env.TABLE_NAME;
const key = process.env.KEY

var tableSvc = azure.createTableService(tabeleName, key);
var TableQuery = azure.TableQuery;
var TableUtilities = azure.TableUtilities;

const port = 5100;

app.listen(port, () => console.log("Server Started on port " + port));
	
app.use(express.static('public'));
app.use(cors());

app.get('/questions/:question/', searchElement);
app.get('/questions/', getAllQuestions);
app.get('/votes/:question', getVotesForQuestion);

function searchElement(request, response) {
    var tableQuery;
    if (request.params.question === 'current')
    {
        tableQuery = new TableQuery().where(TableQuery.booleanFilter('isCurrent', TableUtilities.QueryComparisons.EQUAL, true));
    }
    else
    {
        tableQuery = new TableQuery().where(TableQuery.stringFilter('RowKey', TableUtilities.QueryComparisons.EQUAL, request.params.question));
    }
    
    tableSvc.queryEntities('questions', tableQuery, null, function (error, result) {
        if(!error) {
            var ents = [];
            result.entries.forEach(element => {
                console.log(element);
                var entity = {"PartitionKey": element.PartitionKey._, "RowKey": element.RowKey._, "Text": element.text._, "Answer1Id": parseInt(element.answer1Id._), "Answer1Text": element.answer1Text._,"Answer2Id": parseInt(element.answer2Id._), "Answer2Text": element.answer2Text._, "IsCurrent": element.isCurrent._};
                ents.push(entity);
            });         
            
            response.send(ents);
        }
        else {
            console.log(error);
        }
    }); 
}

function getAllQuestions(request, response) {
    var tableQuery = new TableQuery().where(TableQuery.stringFilter('PartitionKey', TableUtilities.QueryComparisons.EQUAL, 'Questions'));

    tableSvc.queryEntities('questions', tableQuery, null, function (error, result) {
        if(!error) {
            var ents = [];
            result.entries.forEach(element => {
                console.log(element);
                var entity = {"PartitionKey": element.PartitionKey._, "RowKey": element.RowKey._, "Text": element.text._, "Answer1Id": parseInt(element.answer1Id._), "Answer1Text": element.answer1Text._,"Answer2Id": parseInt(element.answer2Id._), "Answer2Text": element.answer2Text._, "IsCurrent": element.isCurrent._};
                ents.push(entity);
            });           
            
            response.send(ents);
        }
        else {
            console.log(error);
        }
    }); 

}

function getVotesForQuestion(request, response) {
    var tableQuery = new TableQuery().where(TableQuery.stringFilter('PartitionKey', TableUtilities.QueryComparisons.EQUAL, request.params.question));
    tableSvc.queryEntities('votes', tableQuery, null, function (error, result) {
        if(!error) {
            var ents = [];
            result.entries.forEach(element => {
                console.log(element);
                var entity = {"PartitionKey": element.PartitionKey._, "RowKey": element.RowKey._, "VoteCount": element.VoteCount._};
                ents.push(entity);
            });
            //var entities = result.entries;
            
            response.send(ents);
            //console.log(entity);
        }
        else {
            console.log(error);
        }
    });
}
