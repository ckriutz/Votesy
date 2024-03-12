from flask import Flask, render_template, request
import os
import time
import socket
import json
import requests
import pika

hostname = socket.gethostname()

connectionString = os.getenv('apiUrl', 'http://localhost:10000')
resultsUrl = os.getenv('resultsURL', "http://localhost:8080")
rabbitmqHost = os.getenv('rabbitmqHost', 'localhost')

app = Flask(__name__)

@app.route("/", methods=['POST','GET'])
def main():
    vote = None
    print(connectionString + '/question/current')
    x = requests.get('http://host.docker.internal:10000/question/current', verify=False)
    print(x.json())
    question = x.json()
    if request.method == 'POST':
        # This is where we send it to queue storage.
        vote = request.form['vote']
        message = {"questionId": question['id'], "answerId" : vote}
        sendMessage(json.dumps(message))

    return render_template("index.html", 
        hostname = hostname, 
        question = question, 
        connectionString = connectionString, 
        vote = vote, 
        resultsUrl = resultsUrl)

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

def sendMessage(message):
    credentials = pika.PlainCredentials('guest', 'guest')
    connection_params = pika.ConnectionParameters(host=rabbitmqHost, port=5672, credentials=credentials)

    connection = pika.BlockingConnection(connection_params)
    channel = connection.channel()

    exchange_name = 'vote-queue'
    routing_key = '/'

    channel.basic_publish(exchange=exchange_name, routing_key=routing_key, body=message)
    print(f"Sent: '{message}'")

    connection.close()


if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5002,debug=True)