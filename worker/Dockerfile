FROM python:3.10-slim AS builder
RUN apt update && \
    apt install --no-install-recommends -y build-essential \
    curl gcc libssl-dev python3-opencv ffmpeg libmagickwand-dev \
    libimage-exiftool-perl exiftool libraw-dev
RUN curl --proto '=https' --tlsv1.3 -sSf https://sh.rustup.rs -o rustup-init.sh
RUN sh rustup-init.sh -y
ENV PATH="/root/.cargo/bin:${PATH}"
ADD requirements.txt /requirements.txt
RUN pip3 install --no-cache-dir -r requirements.txt
COPY src /app/src
EXPOSE 15002
WORKDIR /app
CMD ["python3", "-m", "src.main"]