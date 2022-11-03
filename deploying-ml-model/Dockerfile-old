# Use the official lightweight Python image.
# https://hub.docker.com/_/python
FROM python:3.7-slim
 
# Allow statements and log messages to immediately appear in the Knative logs
ENV PYTHONUNBUFFERED True
 
# Copy local code to the container image.
ENV APP_HOME /app
WORKDIR $APP_HOME
 
# Install production dependencies.
RUN pip install flask gunicorn transformers
RUN pip install torch==1.13.0+cpu -f https://download.pytorch.org/whl/torch_stable.html
 
# The bootstrap script will download the model file into the container image
COPY ./bootstrap.py .
RUN python bootstrap.py
 
# Now copy the app
COPY ./app.py .
 
# Run the web service on container startup. Here we use gunicorn
CMD exec gunicorn --bind :$PORT --workers 1 --threads 1 --timeout 0 app:app
