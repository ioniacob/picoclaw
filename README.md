# 🚀 PicoClaw + Vercel Chat SDK

> **English** | [Español](README_ES.md)

## 📋 Description

This is an enhanced fork of [PicoClaw](https://github.com/sipeed/picoclaw) that adds complete **Vercel Chat SDK** integration to create automated WhatsApp flows with artificial intelligence from multiple providers.

## ✨ New Features

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
- Complete documentation in multiple languages

## 🚀 Quick Deployment

### Option 1: Vercel Chat SDK (Web Panel) - **RECOMMENDED**

```bash
# Quick deployment with web panel
./deploy_vercel_chat.sh --prod

# Admin panel: https://your-project.vercel.app/admin
```

### Option 2: Go + WebSocket (Native)

```bash
# Native Go multi-channel version
./deploy_vercel.sh --prod
```

### Local Development

```bash
# Chat SDK with web panel
./dev.sh

# Native Go
go test ./pkg/channels -v
```

## 📁 Project Structure

```
.
├── api/                    # Vercel handlers
│   ├── chat.js            # Chat with AI
│   ├── whatsapp.js        # WhatsApp with AI
│   └── index.go           # Original Go handler
├── admin/                 # Web panel
│   └── index.html         # Administration
├── pkg/channels/          # Go channels
│   ├── whatsapp.go        # Enhanced WhatsApp
│   ├── whatsapp_test.go   # Unit tests
│   └── whatsapp_validator.go # Validators
├── examples/              # Examples
├── scripts/               # Utility scripts
└── docs/                  # Documentation
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
WHATSAPP_WEBHOOK_TOKEN=your-secret-token
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

## 🧪 Testing

```bash
# Go tests
go test ./pkg/channels -run "TestWhatsApp.*" -v

# Integration tests
./test_whatsapp.sh

# Verify code
npm install
npm run build
```

## 🎨 Admin Panel

The web panel includes:

- ✅ Real-time service status
- ✅ AI provider selector
- ✅ Interactive chat demo
- ✅ WhatsApp configuration
- ✅ Automated flow management

## 🔒 Security

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

- [README.md](README.md) - Complete guide (English) f1ecf1e7
- [README_ES.md](README_ES.md) - Complete guide (Spanish) f1eaf1f8
- [VERCEL_CHAT_SDK.md](VERCEL_CHAT_SDK.md) - Technical documentation
- [VERCEL_DEPLOYMENT.md](VERCEL_DEPLOYMENT.md) - Go deployment
- [VERCEL_README.md](VERCEL_README.md) - Options comparison
- [CHANGES_SUMMARY.md](CHANGES_SUMMARY.md) - Changes summary

## 🤝 Contributing

This development was created by **[@ioniacob](https://github.com/ioniacob)** for the PicoClaw community.

Contributions are welcome! To contribute:

1. Fork this repository
2. Create a feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

## 🌟 Follow Me!

- GitHub: [@ioniacob](https://github.com/ioniacob)
- Twitter: [@ioniacob](https://twitter.com/ioniacob)
- LinkedIn: [linkedin.com/in/ioniacob](https://linkedin.com/in/ioniacob)

Like this project? Give it a ⭐ on GitHub!

## 🙏 Acknowledgments

- [Sipeed](https://github.com/sipeed) for creating PicoClaw
- [Vercel](https://vercel.com) for the Chat SDK
- The PicoClaw community for support

---

## 📄 License

MIT License - maintains original PicoClaw license

---

## 🚀 Ready to deploy your WhatsApp AI assistant!

```bash
git clone https://github.com/ioniacob/picoclaw.git
cd picoclaw
./deploy_vercel_chat.sh --prod
```

**Admin Panel:** `https://your-project.vercel.app/admin`

---

<div align="center">
  <h3>🦐 PicoClaw + Vercel Chat SDK = AI WhatsApp Automation 🚀</h3>
  <p>Developed with ❤️ by <a href="https://github.com/ioniacob">@ioniacob</a></p>
</div>
