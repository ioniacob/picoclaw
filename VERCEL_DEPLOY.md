# Vercel Deployment Configuration

## 🚀 Quick Deploy from GitHub

1. **Connect Repository**: https://github.com/ioniacob/picoclaw
2. **Framework Preset**: Vercel Functions (Mixed Go + Node.js)
3. **Root Directory**: `/`
4. **Build Command**: `npm install && npx vercel build`
5. **Output Directory**: `.vercel/output`

## 🔧 Environment Variables (Required)

Add these in Vercel dashboard:

```bash
# AI Providers
OPENAI_API_KEY=sk-your-openai-key
ANTHROPIC_API_KEY=sk-ant-your-anthropic-key
GROQ_API_KEY=gsk-your-groq-key

# Admin Access
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-password

# WhatsApp Configuration
WHATSAPP_BRIDGE_URL=wss://your-bridge.com/ws
WHATSAPP_ALLOWED_NUMBERS=+1234567890,+0987654321
CONFIG_PATH=/tmp/config.json
NODE_ENV=production
```

## 📁 Function Configuration

The project uses mixed runtimes:

- **Go Functions**: `api/index.go` - Main PicoClaw handler
- **Node.js Functions**:
  - `api/chat.js` - Vercel Chat SDK integration
  - `api/whatsapp.js` - WhatsApp AI flows
- **Static Files**: `admin/index.html` - Web admin panel

## 🔗 Access URLs

After deployment, access your app at:

- **Main API**: `https://your-project.vercel.app/`
- **Chat API**: `https://your-project.vercel.app/api/chat`
- **WhatsApp API**: `https://your-project.vercel.app/api/whatsapp`
- **Admin Panel**: `https://your-project.vercel.app/admin`

## 🧪 Testing

Test your deployment:

```bash
# Test health check
curl https://your-project.vercel.app/health

# Test chat API
curl -X POST https://your-project.vercel.app/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello AI!","provider":"groq"}'

# Test admin access
curl https://your-project.vercel.app/admin
```

## 🎉 Success!

Your PicoClaw WhatsApp AI integration is now live! 🚀
