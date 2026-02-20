# 🚀 PicoClaw + Vercel Chat SDK

> [English](README.md) | **Español**

## 📋 Descripción

Este es un fork mejorado de [PicoClaw](https://github.com/sipeed/picoclaw) que añade integración completa con **Vercel Chat SDK** para crear flujos automáticos de WhatsApp con inteligencia artificial de múltiples proveedores.

## ✨ Características Nuevas

### 🎯 Integración Vercel Chat SDK

- Panel de administración web completo con autenticación segura
- Integración con OpenAI GPT-4, Anthropic Claude, Groq Mixtral
- Streaming de respuestas en tiempo real
- Gestión de sesiones persistentes
- Flujos automáticos configurables

### 🔒 Canal WhatsApp Mejorado

- Validación de mensajes con esquemas estrictos
- Sanitización de contenido para prevenir inyecciones
- Reconexión automática con backoff exponencial
- Seguridad TLS/WSS obligatoria
- Integridad de mensajes con HMAC-SHA256
- Manejo robusto de errores de red
- Keepalive con ping/pong de WebSocket

### 🛠️ Testing y Calidad

- Pruebas unitarias completas para WhatsAppChannel
- Simulación de WebSocket para tests
- Scripts de despliegue automatizados
- Configuración multi-entorno (dev/prod)
- Documentación completa en múltiples idiomas

## 🚀 Despliegue Rápido

### Opción 1: Vercel Chat SDK (Panel Web) - **RECOMENDADO**

```bash
# Despliegue rápido con panel web
./deploy_vercel_chat.sh --prod

# Panel admin: https://tu-proyecto.vercel.app/admin
```

### Opción 2: Go + WebSocket (Nativo)

```bash
# Versión Go nativa multi-canal
./deploy_vercel.sh --prod
```

### Desarrollo Local

```bash
# Chat SDK con panel web
./dev.sh

# Go nativo
go test ./pkg/channels -v
```

## 📁 Estructura del Proyecto

```
.
├── api/                    # Handlers Vercel
│   ├── chat.js            # Chat con AI
│   ├── whatsapp.js        # WhatsApp con AI
│   └── index.go           # Handler Go original
├── admin/                 # Panel web
│   └── index.html         # Administración
├── pkg/channels/          # Canales Go
│   ├── whatsapp.go        # WhatsApp mejorado
│   ├── whatsapp_test.go   # Pruebas unitarias
│   └── whatsapp_validator.go # Validadores
├── examples/              # Ejemplos
├── scripts/               # Scripts de utilidad
└── docs/                  # Documentación
```

## 🔧 Variables de Entorno

```bash
# Proveedores AI (obtén tus claves)
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GROQ_API_KEY=gsk_...

# Admin
ADMIN_USERNAME=admin
ADMIN_PASSWORD=picoclaw123

# WhatsApp
WHATSAPP_BRIDGE_URL=wss://tu-bridge.com/ws
WHATSAPP_ALLOWED_NUMBERS=+1234567890
WHATSAPP_WEBHOOK_TOKEN=tu-token-secreto
```

## 📡 Endpoints Disponibles

### Vercel Chat SDK

- `GET /admin` - Panel de administración
- `POST /api/chat` - Chat con AI
- `POST /api/whatsapp` - WhatsApp con AI
- `GET /api/whatsapp?action=status` - Estado del servicio

### Go Original

- `GET /health` - Health check
- `GET /ready` - Ready check
- `POST /webhook/whatsapp` - Webhook WhatsApp
- `POST /api/chat` - API de chat

## 🧪 Testing

```bash
# Pruebas Go
go test ./pkg/channels -run "TestWhatsApp.*" -v

# Pruebas de integración
./test_whatsapp.sh

# Verificar código
npm install
npm run build
```

## 🎨 Panel de Administración

El panel web incluye:

- ✅ Estado del servicio en tiempo real
- ✅ Selector de proveedor AI
- ✅ Demo de chat interactivo
- ✅ Configuración de WhatsApp
- ✅ Gestión de flujos automáticos

## 🔒 Seguridad

- Validación estricta de mensajes
- Sanitización de entrada de usuario
- Autenticación de administrador
- HTTPS/WSS obligatorio
- Rate limiting implícito
- HMAC para integridad de mensajes

## 📊 Rendimiento

- **< 10MB** de memoria (Go nativo)
- **< 1s** de tiempo de arranque
- **Streaming** en tiempo real
- **Reconexión automática** robusta
- **Multi-proveedor AI** para redundancia

## 🎯 Casos de Uso

1. **Asistente WhatsApp Inteligente** - Responde mensajes automáticamente
2. **Soporte al Cliente 24/7** - Atiende consultas frecuentes
3. **Automatización de Tareas** - Ejecuta acciones por comandos
4. **Panel de Control Web** - Gestión visual de flujos
5. **Integración Multi-Canal** - WhatsApp + otros canales

## 📚 Documentación

- [README.md](README.md) - Guía completa (Inglés) f1ecf1e7
- [README_ES.md](README_ES.md) - Guía completa (Español) f1eaf1f8
- [VERCEL_CHAT_SDK.md](VERCEL_CHAT_SDK.md) - Documentación técnica
- [VERCEL_DEPLOYMENT.md](VERCEL_DEPLOYMENT.md) - Despliegue Go
- [VERCEL_README.md](VERCEL_README.md) - Comparación de opciones
- [CHANGES_SUMMARY.md](CHANGES_SUMMARY.md) - Resumen de cambios

## 🤝 Contribuir

Este desarrollo fue creado por **[@ioniacob](https://github.com/ioniacob)** para la comunidad PicoClaw.

¡Las contribuciones son bienvenidas! Para contribuir:

1. Fork este repositorio
2. Crea una rama para tu feature (`git checkout -b feature/AmazingFeature`)
3. Haz commit de tus cambios (`git commit -m 'Add some AmazingFeature'`)
4. Push a la rama (`git push origin feature/AmazingFeature`)
5. Abre un Pull Request

## 🌟 ¡Sígueme!

- GitHub: [@ioniacob](https://github.com/ioniacob)
- Twitter: [@ioniacob](https://twitter.com/ioniacob)
- LinkedIn: [linkedin.com/in/ioniacob](https://linkedin.com/in/ioniacob)

¿Te gusta este proyecto? ¡Dale una ⭐ en GitHub!

## 🙏 Agradecimientos

- [Sipeed](https://github.com/sipeed) por crear PicoClaw
- [Vercel](https://vercel.com) por el Chat SDK
- La comunidad de PicoClaw por el apoyo

---

## 📄 Licencia

MIT License - mantiene la licencia original de PicoClaw

---

## 🚀 ¡Listo para desplegar tu asistente WhatsApp con AI!

```bash
git clone https://github.com/ioniacob/picoclaw.git
cd picoclaw
./deploy_vercel_chat.sh --prod
```

**Panel Admin:** `https://tu-proyecto.vercel.app/admin`

---

<div align="center">
  <h3>🦐 PicoClaw + Vercel Chat SDK = Automatización WhatsApp con AI 🚀</h3>
  <p>Desarrollado con ❤️ por <a href="https://github.com/ioniacob">@ioniacob</a></p>
</div>
