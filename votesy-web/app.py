from flask import Flask, render_template, request
from azure.storage.queue import (QueueClient, BinaryBase64EncodePolicy, BinaryBase64DecodePolicy)
import os
import uuid
import socket
import json
import requests

hostname = socket.gethostname()

connectionString = os.getenv('apiUrl', 'http://localhost:10000')
resultsUrl = os.getenv('resultsURL', "http://localhost:8080")

app = Flask(__name__)

@app.route("/", methods=['POST','GET'])
def main():
    vote = None
    print(connectionString + '/question/current')
    x = requests.get(connectionString + '/question/current')
    question = x.json()
    if request.method == 'POST':
        # This is where we send it to queue storage.
        vote = request.form['vote']
        message = {"questionId": question['RowKey'], "answerId" : vote}
        sendMessage(json.dumps(message))

    return render_template("index.html", 
        hostname = hostname, 
        question = question, 
        connectionString = connectionString, 
        vote = vote, 
        resultsUrl = resultsUrl)

def sendMessage(message):
    # Get the AZURE_STORAGE_CONNECTION_STRING
    azure_connection_string = os.getenv('AZURE_STORAGE_CONNECTION_STRING', None)

    # Create a QueueClient
    queue_client = QueueClient.from_connection_string(azure_connection_string, 'votes')

    #base64_queue_client = QueueClient.from_connection_string(
    #    conn_str=azure_connection_string, queue_name='votes',
    #    message_encode_policy = BinaryBase64EncodePolicy(),
    #    message_decode_policy = BinaryBase64DecodePolicy()
    #)
    #print("adding message " + message)
    queue_client.send_message(message)

if __name__ == "__main__":
    app.run(host='0.0.0.0', port=5002,debug=True)