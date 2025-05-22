# Jambon

An AI that try to auto resolve Kubernetes warning, using Qwen!

## Web S10 - Groupe 3 :
- Membres :  Rémi ESPIE, Axelandre SOLLIER, Mathias BOULAY
- Excalicraw : https://excalidraw.com/#room=8e326d92283961a77ae6,q3tEB-fYYD-r5WcvvG09EQ

## Jambon's inner workings

Jambon is developed in Go and is composed of 2 main components: the Watcher and the Caller.

### Watcher

The Watcher is a Kubernetes Deployment that watches for events in the cluster. It is responsible for listening to the Kubernetes API and reacting to events that occur in the cluster. It uses the [Kubernetes Go client](https://github.com/kubernetes/client-go).

When a warning is detected, the Watcher will start a Job that will run the Caller. The event will be put in a "to-be resolved" queue, and the Job will be responsible for resolving it. 

### Caller

The Caller is a Kubernetes Job that is responsible for resolving the warning. It will use the Ollama API to interact with the Qwen model and ask it to resolve the warning.

The Caller will first try to resolve the issue by itself. It will find the manifest from a hard coded path, and will use the Qwen model to generate a fix. The fix will be applied to the cluster using the [Kubernetes Go client](https://github.com/kubernetes/client-go).

If the Caller is unable to fix the issue, it will use the Speaches API to transcribe the warning and generate a response, with the idea to make a call to a user to inform him of the problem and ask him to resolve it.

## AI

### Ollama 

[Ollama](https://ollama.com/) is a command line tool that allows you to run large language models (LLMs) locally on a machine. It provides a simple interface for downloading, running, and managing models, as well as an API that we are leveraging to interact with the models.

We are using the [Qwen3](https://ollama.com/library/qwen3) model, the 1.7b version for test purpose and the 8b version for production. Even the small models are really fast to run and perform pretty well.

We are also thinking of using Qwen2.5-coder, a model that is specifically designed for code generation and understanding. It is a smaller model but it is trained on a large amount of code data, which may makes it more suitable for our use case.

We are running Ollama as a service on port 11434. To make it listen on all interfaces, we need to set the `Ollama_HOST` environment variable to `0.0.0.0:11434` before starting the service:
```ini
[Unit]
Description=Ollama Service
After=network-online.target

[Service]
ExecStart=/usr/local/bin/ollama serve
User=ollama
Group=ollama
Restart=always
RestartSec=3
Environment="OLLAMA_HOST=0.0.0.0:11434"
Environment="PATH=/home/debian/.local/bin:/home/debian/.local/bin:/usr/local/bin:/usr/bin:/bin:/usr/local/games:/usr/games"

[Install]
WantedBy=default.target
```

```bash
sudo systemctl daemon-reload
sudo systemctl enable ollama
sudo systemctl restart ollama
```

During usage, the model will be downloaded by the Caller Job if it is not already present.

To interact with it, we are using an unofficial [Ollama Go client](github.com/xyproto/ollamaclient).

### Speaches

[Speaches](https://speaches.ai/) is an OpenAI API-compatible server supporting streaming transcription, translation, and speech generation. Speach-to-Text is powered by [faster-whisper](https://github.com/SYSTRAN/faster-whisper) and for Text-to-Speech [piper](https://github.com/rhasspy/piper) and [Kokoro](https://huggingface.co/hexgrad/Kokoro-82M) are used. This project aims to be Ollama, but for TTS/STT models. It is available on port 8000.

Speaches is still in beta, and because the documentation is up-to-date with the latest code, we are using the release candidate 0.8.0 build. We are also running Speaches as 2 docker containers:
```yaml
services:
  speaches:
    container_name: speaches
    build:
      dockerfile: Dockerfile
      context: .
      platforms:
        - linux/amd64
        - linux/arm64
    restart: unless-stopped
    ports:
      - 8000:8000
    develop:
      watch:
        - action: rebuild
          path: ./uv.lock
        - action: sync+restart
          path: ./src
          target: /home/ubuntu/speaches/src
    healthcheck:
      test: ["CMD", "curl", "--fail", "http://0.0.0.0:8000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 5s
```

```yaml
services:
  speaches:
    extends:
      file: compose.yaml
      service: speaches
    image: ghcr.io/speaches-ai/speaches:0.8.0-rc.2-cpu
    build:
      args:
        BASE_IMAGE: ubuntu:24.04
    volumes:
      - hf-hub-cache:/home/ubuntu/.cache/huggingface/hub
volumes:
  hf-hub-cache:
```

```bash
docker compose up --detach
```

Furthermore, Speaches does not come with prebuilt models, so we need to download them ourselves. We had to install and use their CLI to do so:
```bash
curl -LsSf https://astral.sh/uv/install.sh | sh # Install uv
# uv will install the speaches-cli as well
uvx speaches-cli model download speaches-ai/Kokoro-82M-v1.0-ONNX # TTS
uvx speaches-cli model download Systran/faster-whisper-small # STT
```

Because Speaches is compatible with OpenAI's API, we use the official [OpenAI Go client](https://github.com/openai/openai-go) to interact with it.