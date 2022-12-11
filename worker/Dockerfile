FROM python:3.10-slim
RUN apt update && \
    apt install --no-install-recommends -y build-essential gcc && \
    apt clean && rm -rf /var/lib/apt/lists/*
ADD requirements.txt /requirements.txt
RUN pip3 install --no-cache-dir -r requirements.txt
COPY src /
EXPOSE 15002
CMD ["python3", "main.py"]