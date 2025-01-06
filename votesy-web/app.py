import os
import time
import socket
import json
import requests
from azure.storage.queue import QueueServiceClient
from flask import Flask, render_template, request, make_response

# This application will get the current question from the API and display it to the user.
# It will also allow the user to vote on the question and send the vote to the queue.

hostname = socket.gethostname()

# Get the environment variables
# Connection string to the API
connectionString = os.getenv('API_URL', 'http://localhost:10000')
storageConnectionString = os.getenv('STORAGE_CONNECTION_STRING')

app = Flask(__name__)

def send_message_to_queue(message):
    queue_service_client = QueueServiceClient.from_connection_string(storageConnectionString)
    queue_client = queue_service_client.get_queue_client("votes")
    queue_client.send_message(message)

def get_cookie_for_vote():
    vote_history = request.cookies.get('vote_history')
    if vote_history is None or vote_history == "":
        return []
    return json.loads(vote_history)

@app.route("/", methods=['POST','GET'])
def main():
    history = get_cookie_for_vote()
    print("history: "+ json.dumps(history))
    vote = None
    #print(connectionString + '/questions/current')
    x = requests.get(connectionString + '/questions/current', verify=False)
    #print(x.json())
    question = x.json()
    if request.method == 'POST':
        # This is where we send it to queue storage.
        vote = request.form['vote']
        message = {"answerId" : vote}
        send_message_to_queue(json.dumps(message))
        
        # build out the cookie object to include the questionId, and the answerid that was slected, and add it to the existing cookie if there is one.
        # This will allow us to keep track of the question and the answer that was selected.
        selection = {"question": question['id'], "answerId" : vote}
        history.append(selection)

        resp = make_response(render_template("index.html", 
            hostname = hostname, 
            question = question, 
            voting_history = json.dumps(history), 
            vote = vote))
        resp.set_cookie('vote_history', json.dumps(history))
        return resp

    resp = make_response(render_template("index.html", 
        hostname = hostname, 
        question = question, 
        voting_history = json.dumps(history),
        vote = vote))
    return resp

@app.route("/results", methods=['GET'])
def results():
    # I have to pass in both the question (with the answers),
    # and the number of votes for each answer.
    x = requests.get(connectionString + '/questions/current', verify=False)
    question = x.json()

    # Now get the vote count for each answer.
    x = requests.get(connectionString + '/votes/' + question['PartitionKey'] +'/' + question['RowKey'], verify=False)
    votes = x.json()
    
    combined = {
        "question": question['text'], 
        "answer1": question['answer1Text'], 
        "answer1Votes": 0, 
        "answer2": question['answer2Text'], 
        "answer2Votes": 0, 
        "answer3": question['answer3Text'], 
        "answer3Votes": 0, 
        "answer4": question['answer4Text'], 
        "answer4Votes": 0
    }
    # I need to populate the votes into the combined object from the votes object.
    for vote in votes:
        if vote['id'] == question['answer1Id']:
            combined['answer1Votes'] = vote['voteCount']
        elif vote['id'] == question['answer2Id']:
            combined['answer2Votes'] = vote['voteCount']
        elif vote['id'] == question['answer3Id']:
            combined['answer3Votes'] = vote['voteCount']
        elif vote['id'] == question['answer4Id']:
            combined['answer4Votes'] = vote['voteCount']
    

    #x = requests.get(connectionString + '/votes/' + question['id'], verify=False)
    print(combined)
    return render_template("results.html", hostname = hostname, combined = combined)

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