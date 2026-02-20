# Despliegue de PicoClaw con WhatsApp en Vercel

## 📋 Requisitos previos

1. Cuenta en [Vercel](https://vercel.com)
2. Bridge WebSocket de WhatsApp (ver ejemplos abajo)
3. API key de OpenRouter o proveedor AI

## 🚀 Despliegue rápido

### 1. Configurar variables de entorno en Vercel

```bash
# WhatsApp
WHATSAPP_BRIDGE_URL=wss://tu-bridge-whatsapp.com/ws
WHATSAPP_ALLOWED_NUMBERS=+1234567890,+0987654321
WHATSAPP_WEBHOOK_TOKEN=tu-token-secreto

# AI Provider
OPENROUTER_API_KEY=sk-or-v1-tu-api-key-aqui
```

### 2. Desplegar desde GitHub

1. Conecta tu repositorio en [Vercel Dashboard](https://vercel.com/dashboard)
2. Importa el proyecto desde GitHub
3. Configura las variables de entorno
4. Despliega

### 3. Bridge WebSocket para WhatsApp

Necesitas un bridge que convierta mensajes de WhatsApp a WebSocket. Opciones:

#### Opción A: whatsapp-web.js
```javascript
// server.js
const { Client, LocalAuth } = require('whatsapp-web.js');
const WebSocket = require('ws');

const client = new Client({
    authStrategy: new LocalAuth(),
    puppeteer: { headless: true }
});

const wss = new WebSocket.Server({ port: 3001 });

client.on('message', msg => {
    wss.clients.forEach(ws => {
        ws.send(JSON.stringify({
            type: 'message',
            from: msg.from,
            content: msg.body,
            chat: msg.from,
            id: msg.id.id,
            timestamp: msg.timestamp
        }));
    });
});

wss.on('connection', ws => {
    ws.on('message', data => {
        const msg = JSON.parse(data);
        if (msg.type === 'message' && msg.to && msg.content) {
            client.sendMessage(msg.to, msg.content);
        }
    });
});

client.initialize();
```

#### Opción B: Servicio cloud (recomendado para Vercel)
- [Ultramsg](https://ultramsg.com/)
- [WATI](https://www.wati.io/)
- [Twilio WhatsApp API](https://www.twilio.com/docs/whatsapp/api)

## 🔧 Configuración del webhook

Una vez desplegado, configura el webhook en tu bridge:

```bash
# URL del webhook (reemplaza con tu dominio de Vercel)
WEBHOOK_URL=https://tu-proyecto.vercel.app/webhook/whatsapp

# Token de autenticación
WEBHOOK_TOKEN=tu-token-secreto
```

## 📡 Endpoints disponibles

| Endpoint | Método | Descripción |
|----------|--------|-------------|
| `GET /` | GET | Información del servicio |
| `GET /health` | GET | Health check |
| `GET /ready` | GET | Estado del servicio |
| `POST /webhook/whatsapp` | POST | Webhook para mensajes de WhatsApp |
| `POST /api/chat` | POST | API de chat directo |

## 🔒 Seguridad

### 1. Autenticación webhook
```bash
# Configurar token en Vercel
WHATSAPP_WEBHOOK_TOKEN=mi-super-token-secreto-12345
```

### 2. Validación de números
```bash
# Solo permitir números específicos
WHATSAPP_ALLOWED_NUMBERS=+1234567890,+0987654321
```

### 3. HTTPS obligatorio
- Vercel proporciona HTTPS automáticamente
- El bridge debe usar WSS (WebSocket Secure)

## 🧪 Pruebas

### Verificar despliegue:
```bash
curl https://tu-proyecto.vercel.app/health
```

### Probar webhook:
```bash
curl -X POST https://tu-proyecto.vercel.app/webhook/whatsapp \
  -H "Authorization: Bearer tu-token-secreto" \
  -H "Content-Type: application/json" \
  -d '{
    "type": "message",
    "from": "+1234567890",
    "content": "Hola, ¿cómo estás?",
    "chat": "+1234567890",
    "id": "msg123",
    "timestamp": 1234567890
  }'
```

## 🔍 Monitoreo

- Logs en [Vercel Dashboard](https://vercel.com/dashboard)
- Métricas de uso y errores
- Alertas configurables

## ⚠️ Limitaciones en Vercel

1. **Tiempo de ejecución**: 10s (hobby) / 60s (pro)
2. **Memoria**: 1GB (hobby) / 3GB (pro)
3. **Almacenamiento**: Solo temporal (/tmp)
4. **WebSockets**: Solo inbound (no persistentes)

## 🚀 Optimizaciones recomendadas

1. **Bridge cloud**: Usa un servicio de bridge externo
2. **Redis**: Para estado persistente
3. **CDN**: Para archivos multimedia
4. **Rate limiting**: Implementa límites de uso

## 📞 Soporte

- [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- [Vercel Documentation](https://vercel.com/docs)
- [WhatsApp Business API](https://developers.facebook.com/docs/whatsapp/business-management-api)