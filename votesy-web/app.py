import os
import time
import socket
import json
import requests
from flask import Flask, render_template

# This application will get the current question from the API and display it to the user.
# It will also allow the user to vote on the question and send the vote to the queue.

hostname = socket.gethostname()

# Get the environment variables
# Connection string to the API
connectionString = os.getenv('apiUrl', 'http://localhost:10000')
storageConnectionString = os.getenv('storageConnectionString')

# Create a mock question object to use.
question = {
    "id": 1,
    "text": "What is the best programming language?",
    "answer1Text": "Python",
    "answer1Id": 1,
    "answer2Text": "Java",
    "answer2Id": 2,
    "answer3Text": "C#",
    "answer3Id": 3,
    "answer4Text": "JavaScript",
    "answer4Id": 4,
}

app = Flask(__name__)

@app.route("/", methods=['POST','GET'])
def main():
    vote = None
    print(connectionString + '/questions/current')
    #x = requests.get(connectionString + '/questions/current', verify=False)
    #print(x.json())
    #question = x.json()
    #if request.method == 'POST':
    #    # This is where we send it to queue storage.
    #    vote = request.form['vote']
    #    message = {"questionId": question['id'], "answerId" : vote}
    #    sendMessage(json.dumps(message))

    return render_template("index.html", 
        hostname = hostname, 
        question = question, 
        connectionString = connectionString, 
        vote = vote)

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

#def sendMessage(message):
#    credentials = pika.PlainCredentials('guest', 'guest')
#    connection_params = pika.ConnectionParameters(host=rabbitmqHost, port=5672, credentials=credentials)

#    connection = pika.BlockingConnection(connection_params)
#    channel = connection.channel()

#    exchange_name = 'vote-queue'
#    routing_key = '/'

#    channel.basic_publish(exchange=exchange_name, routing_key=routing_key, body=message)
#    print(f"Sent: '{message}'")

#    connection.close()


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5002,debug=True)