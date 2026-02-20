#!/bin/bash

# Script de desarrollo local para Vercel Chat SDK
set -e

echo "🚀 Iniciando PicoClaw con Vercel Chat SDK en modo desarrollo..."

# Verificar que estamos en el directorio correcto
if [ ! -f "package.json" ]; then
    echo "❌ Error: No se encontró package.json. Ejecuta desde la raíz del proyecto."
    exit 1
fi

# Verificar variables de entorno
if [ ! -f ".env.local" ]; then
    echo "📝 Creando archivo .env.local con configuración por defecto..."
    cp .env.example .env.local
fi

# Cargar variables de entorno
export $(cat .env.local | xargs)

# Verificar dependencias
echo "📦 Verificando dependencias..."
if [ ! -d "node_modules" ]; then
    echo "📥 Instalando dependencias de Node.js..."
    npm install
fi

# Verificar que Vercel CLI esté instalado
if ! command -v vercel &> /dev/null; then
    echo "📥 Instalando Vercel CLI..."
    npm install -g vercel
fi

# Verificar código JavaScript
echo "🔍 Verificando código JavaScript..."
node -c api/chat.js
node -c api/whatsapp.js

if [ $? -ne 0 ]; then
    echo "❌ Error: El código JavaScript tiene errores de sintaxis"
    exit 1
fi

echo "✅ Código JavaScript válido"

# Iniciar servidor de desarrollo
echo "🌐 Iniciando servidor de desarrollo..."
echo ""
echo "📋 Endpoints disponibles:"
echo "  - Panel de admin: http://localhost:3000/admin"
echo "  - API de chat: http://localhost:3000/api/chat"
echo "  - API de WhatsApp: http://localhost:3000/api/whatsapp"
echo "  - Estado WhatsApp: http://localhost:3000/api/whatsapp?action=status"
echo ""
echo "🔧 Variables de entorno:"
echo "  - Admin user: $ADMIN_USERNAME"
echo "  - Admin pass: $ADMIN_PASSWORD"
echo "  - AI Provider: ${OPENAI_API_KEY:+OpenAI} ${ANTHROPIC_API_KEY:+Anthropic} ${GROQ_API_KEY:+Groq}"
echo ""
echo "⚠️  Asegúrate de tener configuradas tus API keys en .env.local"
echo ""

# Iniciar Vercel dev
vercel dev --listen 3000