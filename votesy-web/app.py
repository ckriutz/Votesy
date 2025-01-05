import os
import time
import socket
import json
import requests
from azure.storage.queue import QueueServiceClient
from flask import Flask, render_template, request

# This application will get the current question from the API and display it to the user.
# It will also allow the user to vote on the question and send the vote to the queue.

hostname = socket.gethostname()

# Get the environment variables
# Connection string to the API
connectionString = os.getenv('API_URL', 'http://localhost:10000')
storageConnectionString = os.getenv('STORAGE_CONNECTION_STRING')

# Create a mock question object to use.
question = {
    "PartitionKey":"questions",
    "RowKey":"wzso1",
    "Timestamp":"2024-12-24T03:23:08.3550273Z",
    "id":"wzso1",
    "text":"Bear or Owl?",
    "answer1Id":"3ief7",
    "answer1Text":"Bear",
    "answer2Id":"4cdhw",
    "answer2Text":"Owl",
    "answer3Id":"",
    "answer3Text":"",
    "answer4Id":"",
    "answer4Text":"",
    "isCurrent":"true",
    "isUsed":"false",
    "CreatedDate":"2024-12-24T03:23:08.381516842Z"
}

# Mok return for votes:
votes = [{"PartitionKey":"votes","RowKey":"3ief7","Timestamp":"2024-12-24T17:53:08.4070264Z","id":"3ief7","voteCount":11},{"PartitionKey":"votes","RowKey":"4cdhw","Timestamp":"2024-12-24T17:53:14.0582694Z","id":"4cdhw","voteCount":2}]

app = Flask(__name__)

def send_message_to_queue(message):
    queue_service_client = QueueServiceClient.from_connection_string(storageConnectionString)
    queue_client = queue_service_client.get_queue_client("votes")
    queue_client.send_message(message)

@app.route("/", methods=['POST','GET'])
def main():
    vote = None
    print(connectionString + '/questions/current')
    #x = requests.get(connectionString + '/questions/current', verify=False)
    #print(x.json())
    #question = x.json()
    if request.method == 'POST':
        # This is where we send it to queue storage.
        vote = request.form['vote']
        message = {"answerId" : vote}
        send_message_to_queue(json.dumps(message))

    # Get the vote count for answer1Id
    

    return render_template("index.html", 
        hostname = hostname, 
        question = question, 
        connectionString = storageConnectionString, 
        vote = vote)

@app.route("/results", methods=['GET'])
def results():
    # I have to pass in both the question (with the answers),
    # and the number of votes for each answer.
    print(connectionString + '/votes/' + question["PartitionKey"] + question['RowKey'])
    answer1_vote_count = next((item['voteCount'] for item in votes if item['RowKey'] == question['answer1Id']), 0)
    #x = requests.get(connectionString + '/votes/' + question['id'], verify=False)
    #print(x.json())
    return render_template("results.html", hostname = hostname, question = question, votes = votes,
        answer1_vote_count = answer1_vote_count)

@app.route("/health/readiness", methods=['GET'])
def readiness():
    time.sleep(1)
    resp = {"status_code":200}
    return json.dumps(resp)

@app.route("/health/liveness", methods=['GET'])
def liveness():
    time.sleep(.25)
    resp = {"status_code":200}
    return json.dumps(resp)

@app.route("/health/startup", methods=['GET'])
def startup():
    time.sleep(.5)
    resp = {"status_code":200}
    return json.dumps(resp)

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5002,debug=True)