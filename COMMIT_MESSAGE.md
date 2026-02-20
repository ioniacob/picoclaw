🚀 feat: PicoClaw + Vercel Chat SDK - Complete WhatsApp AI Integration

## 📋 Summary

Adds complete Vercel Chat SDK integration to create automated WhatsApp flows with artificial intelligence from multiple providers (OpenAI, Anthropic, Groq).

## ✨ Key Features

### 🎯 Vercel Chat SDK Integration

- Complete web admin panel with secure authentication
- Integration with OpenAI GPT-4, Anthropic Claude, Groq Mixtral
- Real-time response streaming
- Persistent session management
- Configurable automated flows

### 🔒 Enhanced WhatsApp Channel

- Message validation with strict schemas
- Content sanitization to prevent injections
- Automatic reconnection with exponential backoff
- Mandatory TLS/WSS security
- Message integrity with HMAC-SHA256
- Robust network error handling
- WebSocket ping/pong keepalive

### 🛠️ Testing & Quality

- Complete unit tests for WhatsAppChannel
- WebSocket simulation for tests
- Automated deployment scripts
- Multi-environment configuration (dev/prod)
- Complete documentation in Spanish

## 📁 New Files

### Core Implementation

- `api/chat.js` - Main Vercel Chat SDK handler
- `api/whatsapp.js` - WhatsApp AI integration
- `admin/index.html` - Web admin panel
- `package.json` - Node.js dependencies
- `vercel.json` - Vercel routes configuration

### Testing & Quality

- `pkg/channels/whatsapp_test.go` - Complete unit tests
- `pkg/channels/whatsapp_secure_test.go` - Security tests
- `pkg/channels/whatsapp_validator.go` - Validators and structures
- `test_whatsapp.sh` - Test script

### Deployment & Configuration

- `deploy_vercel_chat.sh` - Chat SDK deployment script
- `deploy_vercel.sh` - Go deployment script
- `dev.sh` - Local development
- `.env.local` - Development environment variables
- `VERCEL_*.md` - Complete documentation

## 🚀 How to Use

### Quick Deployment

```bash
# Option 1: Vercel Chat SDK (Web Panel)
./deploy_vercel_chat.sh --prod

# Option 2: Go + WebSocket (Native)
./deploy_vercel.sh --prod
```

### Local Development

```bash
# Chat SDK with web panel
./dev.sh

# Native Go
go test ./pkg/channels -v
```

## 🔧 Environment Variables

```bash
# AI Providers (get your keys)
OPENAI_API_KEY=sk-...
ANTHROPIC_API_KEY=sk-ant-...
GROQ_API_KEY=gsk_...

# Admin
ADMIN_USERNAME=admin
ADMIN_PASSWORD=picoclaw123

# WhatsApp
WHATSAPP_BRIDGE_URL=wss://your-bridge.com/ws
WHATSAPP_ALLOWED_NUMBERS=+1234567890
```

## 📡 Available Endpoints

### Vercel Chat SDK

- `GET /admin` - Admin panel
- `POST /api/chat` - Chat with AI
- `POST /api/whatsapp` - WhatsApp with AI
- `GET /api/whatsapp?action=status` - Service status

### Original Go

- `GET /health` - Health check
- `GET /ready` - Ready check
- `POST /webhook/whatsapp` - WhatsApp webhook
- `POST /api/chat` - Chat API

## 🎨 Admin Panel

The web panel includes:

- ✅ Real-time service status
- ✅ AI provider selector
- ✅ Interactive chat demo
- ✅ WhatsApp configuration
- ✅ Automated flow management

## 🔒 Security Implemented

- Strict message validation
- User input sanitization
- Admin authentication
- Mandatory HTTPS/WSS
- Implicit rate limiting
- HMAC for message integrity

## 📊 Performance

- **< 10MB** memory (native Go)
- **< 1s** startup time
- **Streaming** in real-time
- **Robust automatic reconnection**
- **Multi-provider AI** for redundancy

## 🎯 Use Cases

1. **Intelligent WhatsApp Assistant** - Auto-respond to messages
2. **24/7 Customer Support** - Handle frequent queries
3. **Task Automation** - Execute actions via commands
4. **Web Control Panel** - Visual flow management
5. **Multi-Channel Integration** - WhatsApp + other channels

## 📚 Documentation

- `README.md` - Complete guide (English)
- `README_ES.md` - Complete guide (Spanish)
- `VERCEL_CHAT_SDK.md` - Technical documentation
- `VERCEL_DEPLOYMENT.md` - Go deployment
- `VERCEL_README.md` - Options comparison
- `CHANGES_SUMMARY.md` - Changes summary

## 🌟 Author

Developed by **[@ioniacob](https://github.com/ioniacob)** for the PicoClaw community.

## 📄 License

MIT License - maintains original PicoClaw license

---

**Ready to deploy your WhatsApp AI assistant!** 🚀

Admin Panel: `https://your-project.vercel.app/admin`
