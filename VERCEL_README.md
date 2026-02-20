# PicoClaw - WhatsApp AI con Vercel

## 🚀 Opciones de Despliegue

### Opción 1: Vercel Chat SDK (Recomendado)

Implementación moderna con panel web de administración y Vercel Chat SDK.

**Características:**
- ✅ Panel de administración web
- ✅ Vercel Chat SDK integrado
- ✅ Múltiples proveedores AI (OpenAI, Anthropic, Groq)
- ✅ Flujos automáticos configurables
- ✅ Streaming de respuestas en tiempo real

**Desplegar:**
```bash
# Opción rápida
./deploy_vercel_chat.sh

# O con Vercel CLI
npm install
vercel deploy --prod
```

**Panel Admin:** `https://tu-proyecto.vercel.app/admin`

📖 [Documentación completa](VERCEL_CHAT_SDK.md)

---

### Opción 2: Go + WebSocket (Original)

Implementación original en Go con WebSocket y canales múltiples.

**Características:**
- ✅ Código Go original
- ✅ Múltiples canales (WhatsApp, Telegram, Discord, etc.)
- ✅ WebSocket seguro con reconexión automática
- ✅ Validación y sanitización de mensajes
- ✅ Integridad de mensajes con HMAC

**Desplegar:**
```bash
# Configurar variables de entorno
export OPENROUTER_API_KEY=sk-or-v1-...
export WHATSAPP_BRIDGE_URL=wss://...

# Desplegar
./deploy_vercel.sh --prod
```

📖 [Documentación Go](VERCEL_DEPLOYMENT.md)

---

## 📋 Comparación

| Característica | Chat SDK | Go Original |
|----------------|----------|-------------|
| Panel Web | ✅ | ❌ |
| Vercel Chat SDK | ✅ | ❌ |
| Múltiples AI | ✅ | ✅ |
| WhatsApp | ✅ | ✅ |
| Otros canales | ❌ | ✅ |
| Streaming | ✅ | ❌ |
| WebSocket | ❌ | ✅ |
| Go nativo | ❌ | ✅ |

---

## 🛠️ Desarrollo Local

### Chat SDK (Recomendado)
```bash
# Instalar dependencias
npm install

# Configurar entorno
cp .env.example .env.local
# Edita .env.local con tus keys

# Iniciar desarrollo
./dev.sh
# O: vercel dev
```

### Go Original
```bash
# Ejecutar tests
go test ./pkg/channels -run "TestWhatsApp.*" -v

# Iniciar servidor
go run cmd/picoclaw/main.go
```

---

## 🔧 Configuración de WhatsApp

### Bridge WhatsApp
Necesitas un bridge WebSocket que conecte WhatsApp con PicoClaw:

**Opción 1: Local**
```bash
cd examples
npm install
node whatsapp-bridge.js
```

**Opción 2: Cloud**
- [Ultramsg](https://ultramsg.com/)
- [WATI](https://www.wati.io/)
- [Twilio WhatsApp](https://www.twilio.com/docs/whatsapp/api)

### Variables de Entorno
```bash
# WhatsApp
WHATSAPP_BRIDGE_URL=wss://tu-bridge.com/ws
WHATSAPP_ALLOWED_NUMBERS=+1234567890,+0987654321
WHATSAPP_WEBHOOK_TOKEN=tu-token-secreto

# AI Providers
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GROQ_API_KEY=gsk_...
OPENROUTER_API_KEY=sk-or-v1-...
```

---

## 📡 Endpoints

### Chat SDK
- `GET /admin` - Panel de administración
- `POST /api/chat` - API de chat con AI
- `POST /api/whatsapp` - API de WhatsApp
- `GET /api/whatsapp?action=status` - Estado del servicio

### Go Original
- `GET /health` - Health check
- `GET /ready` - Ready check
- `POST /webhook/whatsapp` - Webhook WhatsApp
- `POST /api/chat` - API de chat

---

## 🔒 Seguridad

### Chat SDK
- Autenticación de administrador
- Validación de números permitidos
- HTTPS/WSS obligatorio
- Rate limiting implícito

### Go Original
- Validación de mensajes con HMAC
- Sanitización de contenido
- Reconexión automática segura
- WebSocket con TLS

---

## 🎨 Personalización

### Chat SDK
```javascript
// api/chat.js - Agregar nuevo provider
import { nuevoProvider } from '@ai-sdk/nuevo';
const providers = {
  ...existentes,
  nuevo: nuevoProvider('model')
};
```

### Go Original
```go
// pkg/channels/whatsapp.go - Modificar validación
func (c *WhatsAppChannel) validateMessage(msg interface{}) error {
  // Tu validación personalizada
}
```

---

## 📊 Monitoreo

### Chat SDK
- Panel web con estado en tiempo real
- Logs en Vercel Dashboard
- Métricas de uso por provider

### Go Original
- Logs estructurados
- Health checks
- Métricas de conexión

---

## 🆘 Soporte

**Problemas comunes:**
1. **WhatsApp no conecta** → Verifica QR code y bridge
2. **AI no responde** → Revisa API keys y límites
3. **Error 500** → Verifica variables de entorno
4. **WebSocket falla** → Comprueba certificados TLS

**Recursos:**
- 📖 [Docs Chat SDK](VERCEL_CHAT_SDK.md)
- 📖 [Docs Go](VERCEL_DEPLOYMENT.md)
- 🐛 [GitHub Issues](https://github.com/sipeed/picoclaw/issues)
- 💬 [Vercel Discord](https://discord.gg/vercel)

---

## 🎯 Recomendación

**Para nuevos proyectos:** Usa **Vercel Chat SDK** por su facilidad de uso y panel web.

**Para integraciones complejas:** Usa **Go Original** por su soporte multi-canal y WebSocket robusto.

---

¡Listo para desplegar! 🚀

```bash
# Opción rápida con panel web
./deploy_vercel_chat.sh --prod

# O la versión Go multi-canal
./deploy_vercel.sh --prod
```