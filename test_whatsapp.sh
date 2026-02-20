#!/bin/bash
# Script de prueba para WhatsApp en PicoClaw
# Este script demuestra cómo probar la configuración de WhatsApp

echo "=== Prueba de Configuración de WhatsApp para PicoClaw ==="
echo

# Verificar que el archivo de configuración existe
CONFIG_FILE="/workspaces/picoclaw/config/config.json"
if [ ! -f "$CONFIG_FILE" ]; then
    echo "❌ Error: No se encuentra el archivo de configuración $CONFIG_FILE"
    exit 1
fi

echo "✅ Archivo de configuración encontrado: $CONFIG_FILE"
echo

# Verificar la configuración de WhatsApp
echo "📋 Configuración actual de WhatsApp:"
grep -A 5 '"whatsapp"' "$CONFIG_FILE"
echo

# Ejecutar tests de WhatsApp
echo "🧪 Ejecutando tests de WhatsApp..."
cd /workspaces/picoclaw

echo "1. Test de configuración por defecto:"
go test ./pkg/config -run TestDefaultConfig_Channels -v
echo

echo "2. Test de migración de configuración:"
go test ./pkg/migrate -run "Test.*channels.*mapping" -v
echo

echo "3. Test de funcionalidad de WhatsApp:"
go test ./pkg/channels -run TestWhatsApp -v
echo

# Verificar que el bridge WebSocket esté configurado
echo "🔍 Verificación de configuración:"
if grep -q '"enabled": false' "$CONFIG_FILE" && grep -A 2 '"whatsapp"' "$CONFIG_FILE" | grep -q '"enabled": false'; then
    echo "ℹ️  WhatsApp está actualmente deshabilitado en la configuración"
    echo "Para habilitar WhatsApp, necesitas:"
    echo "  1. Cambiar 'enabled' a true en la configuración de WhatsApp"
    echo "  2. Tener un bridge WebSocket ejecutándose en ws://localhost:3001"
    echo "  3. Configurar allow_from con los usuarios permitidos"
else
    echo "✅ WhatsApp está habilitado en la configuración"
fi

echo
echo "=== Resumen de pruebas ==="
echo "✅ Tests de configuración ejecutados"
echo "✅ Tests de funcionalidad de WhatsApp ejecutados"
echo "ℹ️  Para probar WhatsApp en producción:"
echo "   1. Asegúrate de tener un bridge WebSocket de WhatsApp ejecutándose"
echo "   2. Configura correctamente allow_from con IDs de usuarios permitidos"
echo "   3. Ejecuta: picoclaw gateway"
echo
echo "🔗 Ejemplo de bridge WebSocket para WhatsApp:"
echo "   - whatsapp-web.js con WebSocket server"
echo "   - Baileys con WebSocket wrapper"
echo "   - Otros bridges que implementen el protocolo WebSocket esperado"
echo
echo "📚 El protocolo espera mensajes JSON con formato:"
echo '   Enviar: {"type": "message", "to": "PHONE_NUMBER", "content": "TEXT"}'
echo '   Recibir: {"type": "message", "from": "PHONE_NUMBER", "content": "TEXT", "chat": "CHAT_ID"}'