#!/bin/bash

#Скрипт обновления и перезапуска telegram-vpn-bot

set -e  # Прерывать выполнение при ошибках

echo "🔄 Pulling latest changes from Git..."
git pull

echo "🔀 Merging origin/main..."
git merge origin/main

echo "🐳 Rebuilding and restarting docker container..."
docker compose up telegram-vpn-bot --build -d

echo "✅ Deployment complete!"

