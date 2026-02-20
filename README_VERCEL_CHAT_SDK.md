# 🚀 NUEVAS CARACTERÍSTICAS: PicoClaw + Vercel Chat SDK

## 📋 Resumen de Desarrollo

Este fork de PicoClaw incluye integración completa con **Vercel Chat SDK** para crear flujos automáticos de WhatsApp con inteligencia artificial. Hemos desarrollado una solución completa que combina la eficiencia de PicoClaw con la potencia de los principales proveedores de AI.

## ✨ Características Principales Añadidas

### 🎯 Vercel Chat SDK Integration

- **Panel de Administración Web** completo con autenticación
- **Múltiples proveedores AI**: OpenAI GPT-4, Anthropic Claude, Groq Mixtral
- **Streaming de respuestas** en tiempo real
- **Gestión de sesiones** persistentes
- **Flujos automáticos** configurables

### 🔒 WhatsApp Channel Mejorado

- **Validación de mensajes** con esquemas estrictos
- **Sanitización de contenido** para prevenir inyecciones
- **Reconexión automática** con backoff exponencial
- **Seguridad TLS/WSS** obligatoria
- **Integridad de mensajes** con HMAC-SHA256
- **Manejo robusto de errores** de red
- **Keepalive con ping/pong** de WebSocket

### 🛠️ Infraestructura y Testing

- **Pruebas unitarias** completas para WhatsAppChannel
- **Simulación de WebSocket** para tests
- **Scripts de despliegue** automatizados
- **Configuración multi-entorno** (dev/prod)
- **Documentación completa** en español

## 🚀 Opciones de Despliegue

### Opción 1: Vercel Chat SDK (Recomendado)

```bash
# Despliegue rápido
./deploy_vercel_chat.sh --prod

# Panel admin: https://tu-proyecto.vercel.app/admin
```

### Opción 2: Go + WebSocket (Original)

```bash
# Versión Go nativa
./deploy_vercel.sh --prod
```

## 📁 Archivos Nuevos Creados

### Core Implementation

- `api/chat.js` - Handler principal de Vercel Chat SDK
- `api/whatsapp.js` - Integración WhatsApp con AI
- `admin/index.html` - Panel de administración web
- `package.json` - Dependencias Node.js
- `vercel.json` - Configuración de rutas

### Testing & Quality

- `pkg/channels/whatsapp_test.go` - Pruebas unitarias completas
- `pkg/channels/whatsapp_secure_test.go` - Pruebas de seguridad
- `pkg/channels/whatsapp_validator.go` - Validadores y estructuras
- `test_whatsapp.sh` - Script de pruebas

### Deployment & Configuration

- `deploy_vercel_chat.sh` - Script de despliegue Chat SDK
- `deploy_vercel.sh` - Script de despliegue Go
- `dev.sh` - Desarrollo local
- `.env.local` - Variables de entorno
- `.env.example` - Plantilla de variables

### Documentation

- `VERCEL_CHAT_SDK.md` - Documentación completa
- `VERCEL_DEPLOYMENT.md` - Guía de despliegue
- `VERCEL_README.md` - Comparación de opciones

## 🔧 Variables de Entorno

```bash
# AI Providers
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GROQ_API_KEY=gsk_...

# Admin
ADMIN_USERNAME=admin
ADMIN_PASSWORD=picoclaw123

# WhatsApp
WHATSAPP_BRIDGE_URL=wss://tu-bridge.com/ws
WHATSAPP_ALLOWED_NUMBERS=+1234567890
```

## 📡 Endpoints

- `GET /admin` - Panel de administración
- `POST /api/chat` - Chat con AI
- `POST /api/whatsapp` - WhatsApp con AI
- `GET /api/whatsapp?action=status` - Estado del servicio

## 🧪 Testing

```bash
# Pruebas Go
go test ./pkg/channels -run "TestWhatsApp.*" -v

# Pruebas de integración
./test_whatsapp.sh

# Desarrollo local
./dev.sh
```

## 🎯 Casos de Uso

1. **Asistente WhatsApp Inteligente** - Responde mensajes automáticamente
2. **Soporte al Cliente 24/7** - Atiende consultas frecuentes
3. **Automatización de Tareas** - Ejecuta acciones por comandos
4. **Integración Multi-Canal** - WhatsApp + otros canales
5. **Panel de Control Web** - Gestión visual de flujos

## 🔒 Seguridad Implementada

- ✅ Validación estricta de mensajes
- ✅ Sanitización de entrada de usuario
- ✅ Autenticación de administrador
- ✅ HTTPS/WSS obligatorio
- ✅ Rate limiting implícito
- ✅ HMAC para integridad
- ✅ Manejo seguro de errores

## 📊 Rendimiento

- **< 10MB** de memoria (Go nativo)
- **< 1s** de tiempo de arranque
- **Streaming** en tiempo real
- **Reconexión automática** robusta
- **Multi-proveedor AI** para redundancia

## 🤝 Contribuir

Este desarrollo fue creado por [@ioniacob](https://github.com/ioniacob) para la comunidad PicoClaw. ¡Contribuciones son bienvenidas!

### Cómo contribuir:

1. Fork este repositorio
2. Crea una rama para tu feature
3. Haz commit de tus cambios
4. Push a la rama
5. Abre un Pull Request

## 📄 Licencia

MIT License - ver archivo LICENSE original

---

## 🌟 ¡Sígueme!

Desarrollado con ❤️ por **@ioniacob**

- GitHub: [github.com/ioniacob](https://github.com/ioniacob)
- Twitter: [@ioniacob](https://twitter.com/ioniacob)
- LinkedIn: [linkedin.com/in/ioniacob](https://linkedin.com/in/ioniacob)

¿Te gusta este proyecto? ¡Dale una ⭐ en GitHub!

---

## 📚 Recursos Adicionales

- [Documentación Vercel Chat SDK](https://sdk.vercel.ai/docs)
- [Guía WhatsApp Business API](https://developers.facebook.com/docs/whatsapp)
- [PicoClaw Original](https://github.com/sipeed/picoclaw)

---

**¡Listo para desplegar!** 🚀

```bash
git clone https://github.com/ioniacob/picoclaw.git
cd picoclaw
./deploy_vercel_chat.sh --prod
```
