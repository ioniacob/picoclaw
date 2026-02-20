#!/bin/bash

# Script de despliegue para Vercel con Chat SDK
set -e

echo "🚀 Desplegando PicoClaw con Vercel Chat SDK..."

# Verificar que estamos en el directorio correcto
if [ ! -f "package.json" ]; then
    echo "❌ Error: No se encontró package.json. Ejecuta desde la raíz del proyecto."
    exit 1
fi

# Verificar dependencias
if ! command -v node &> /dev/null; then
    echo "❌ Error: Node.js no está instalado"
    exit 1
fi

if ! command -v npm &> /dev/null; then
    echo "❌ Error: npm no está instalado"
    exit 1
fi

# Instalar dependencias
echo "📦 Instalando dependencias..."
npm install

# Verificar variables de entorno necesarias
if [ -z "$OPENAI_API_KEY" ] && [ -z "$ANTHROPIC_API_KEY" ] && [ -z "$GROQ_API_KEY" ]; then
    echo "⚠️  Advertencia: No se encontraron API keys de AI configuradas"
    echo "Configura al menos una de estas variables de entorno:"
    echo "  - OPENAI_API_KEY"
    echo "  - ANTHROPIC_API_KEY" 
    echo "  - GROQ_API_KEY"
fi

# Crear archivo .env si no existe
if [ ! -f ".env.local" ]; then
    echo "📝 Creando archivo .env.local..."
    cat > .env.local << EOF
# Configuración de administrador
ADMIN_USERNAME=admin
ADMIN_PASSWORD=picoclaw123

# API Keys para AI (configura al menos una)
OPENAI_API_KEY=${OPENAI_API_KEY:-}
ANTHROPIC_API_KEY=${ANTHROPIC_API_KEY:-}
GROQ_API_KEY=${GROQ_API_KEY:-}

# Configuración de WhatsApp
WHATSAPP_BRIDGE_URL=${WHATSAPP_BRIDGE_URL:-wss://api.whatsapp.com/v1}
WHATSAPP_ALLOWED_NUMBERS=${WHATSAPP_ALLOWED_NUMBERS:-}
WHATSAPP_WEBHOOK_TOKEN=${WHATSAPP_WEBHOOK_TOKEN:-tu-token-secreto}

# Configuración de PicoClaw
CONFIG_PATH=/tmp/config.json
NODE_ENV=production

# Configuración de proveedores
OPENROUTER_API_KEY=${OPENROUTER_API_KEY:-}
OPENROUTER_API_BASE=https://openrouter.ai/api/v1
EOF
fi

# Verificar que el código JavaScript es válido
echo "🔍 Verificando código JavaScript..."
node -c api/chat.js
node -c api/whatsapp.js

if [ $? -ne 0 ]; then
    echo "❌ Error: El código JavaScript tiene errores de sintaxis"
    exit 1
fi

# Instalar Vercel CLI si no está instalado
if ! command -v vercel &> /dev/null; then
    echo "📥 Instalando Vercel CLI..."
    npm install -g vercel
fi

# Ejecutar tests si existen
if [ -f "test/api.test.js" ]; then
    echo "🧪 Ejecutando tests..."
    npm test
fi

# Desplegar
echo "🌐 Desplegando en Vercel..."

# Crear archivo de configuración para el despliegue
cat > deploy-config.json << EOF
{
  "name": "picoclaw-vercel-chat",
  "version": 2,
  "buildCommand": null,
  "outputDirectory": null,
  "installCommand": "npm install",
  "devCommand": "npm run dev"
}
EOF

# Desplegar con Vercel
if [ "$1" == "--prod" ]; then
    echo "🚀 Desplegando a producción..."
    vercel deploy --prod --yes --local-config=vercel.json
else
    echo "🔧 Desplegando a preview..."
    vercel deploy --yes --local-config=vercel.json
fi

# Limpiar archivos temporales
rm -f deploy-config.json

echo "✅ Despliegue completado!"
echo ""
echo "📖 Próximos pasos:"
echo "1. Configura las variables de entorno en el panel de Vercel"
echo "2. Visita el panel de administración: https://tu-proyecto.vercel.app/admin"
echo "3. Configura tu bridge WhatsApp para usar el webhook"
echo "4. Verifica la documentación: cat VERCEL_CHAT_SDK.md"
echo ""
echo "🔗 Recursos:"
echo "- Panel de Vercel: https://vercel.com/dashboard"
echo "- Documentación: VERCEL_CHAT_SDK.md"
echo "- Ejemplos de bridge: examples/whatsapp-bridge.js"