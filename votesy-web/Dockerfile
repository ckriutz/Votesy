# For more information, please refer to https://aka.ms/vscode-docker-python
FROM python:3.13-slim-bullseye

WORKDIR /app

# Install pip requirements
COPY requirements.txt ./requirements.txt
RUN python -m pip install -r /app/requirements.txt

COPY . .

EXPOSE 5002

# Keeps Python from generating .pyc files in the container
ENV PYTHONDONTWRITEBYTECODE=1

# Turns off buffering for easier container logging
ENV PYTHONUNBUFFERED=1

# During debugging, this entry point will be overridden. For more information, please refer to https://aka.ms/vscode-docker-python-debug
CMD ["gunicorn", "--bind", "0.0.0.0:5002", "app:app"]
