FROM python:3.8-slim
ADD requirements.txt /requirements.txt
RUN pip install -r requirements.txt
COPY . .
EXPOSE 5002
EXPOSE 15002
CMD ["python", "main.py"]