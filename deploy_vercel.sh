#!/bin/bash

# Script de despliegue para Vercel con WhatsApp
set -e

echo "🚀 Desplegando PicoClaw con WhatsApp en Vercel..."

# Verificar que estamos en el directorio correcto
if [ ! -f "go.mod" ]; then
    echo "❌ Error: No se encontró go.mod. Ejecuta desde la raíz del proyecto."
    exit 1
fi

# Verificar variables de entorno necesarias
if [ -z "$OPENROUTER_API_KEY" ]; then
    echo "❌ Error: OPENROUTER_API_KEY no está configurada"
    echo "Por favor configura tu API key: export OPENROUTER_API_KEY=sk-or-v1-..."
    exit 1
fi

if [ -z "$WHATSAPP_BRIDGE_URL" ]; then
    echo "❌ Error: WHATSAPP_BRIDGE_URL no está configurada"
    echo "Por favor configura la URL de tu bridge: export WHATSAPP_BRIDGE_URL=wss://..."
    exit 1
fi

# Construir el proyecto
echo "📦 Construyendo el proyecto..."
go build -o /tmp/picoclaw ./cmd/picoclaw/main.go

# Verificar que el binario se creó
if [ ! -f "/tmp/picoclaw" ]; then
    echo "❌ Error: Falló la construcción del binario"
    exit 1
fi

echo "✅ Binario construido exitosamente"

# Ejecutar tests
echo "🧪 Ejecutando tests..."
go test ./pkg/channels -run "TestWhatsApp.*" -v

echo "✅ Tests pasados"

# Instalar Vercel CLI si no está instalado
if ! command -v vercel &> /dev/null; then
    echo "📥 Instalando Vercel CLI..."
    npm install -g vercel
fi

# Desplegar
echo "🌐 Desplegando en Vercel..."

# Crear archivo de variables de entorno para Vercel
cat > .env.production << EOF
OPENROUTER_API_KEY=$OPENROUTER_API_KEY
WHATSAPP_BRIDGE_URL=$WHATSAPP_BRIDGE_URL
WHATSAPP_ALLOWED_NUMBERS=${WHATSAPP_ALLOWED_NUMBERS:-}
WHATSAPP_WEBHOOK_TOKEN=${WHATSAPP_WEBHOOK_TOKEN:-}
ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY:-}
OPENAI_API_KEY=${OPENAI_API_KEY:-}
EOF

echo "✅ Archivo de entorno creado"
echo "📋 Variables configuradas:"
echo "  - WHATSAPP_BRIDGE_URL: $WHATSAPP_BRIDGE_URL"
echo "  - WHATSAPP_ALLOWED_NUMBERS: ${WHATSAPP_ALLOWED_NUMBERS:-no configurado}"
echo "  - WHATSAPP_WEBHOOK_TOKEN: ${WHATSAPP_WEBHOOK_TOKEN:-no configurado}"

# Desplegar con Vercel
if [ "$1" == "--prod" ]; then
    echo "🚀 Desplegando a producción..."
    vercel --prod --yes
else
    echo "🔧 Desplegando a preview..."
    vercel --yes
fi

echo "✅ Despliegue completado!"
echo "📖 Verifica el estado en: https://vercel.com/dashboard"
echo "🔗 Configura tu bridge WhatsApp para usar el webhook generado"
echo "📚 Documentación: cat VERCEL_DEPLOYMENT.md"