# Vercel Deployment Configuration

## 🚀 Quick Deploy to Vercel

### 1. Connect to Vercel Dashboard

1. Visit: https://vercel.com
2. Click "New Project"
3. Import Git Repository: `https://github.com/ioniacob/picoclaw`
4. Configure settings (see below)
5. Deploy!

### 2. Project Settings

- **Framework**: Vercel Functions
- **Root Directory**: `/`
- **Build Command**: `npm install`
- **Output Directory**: `public`
- **Install Command**: `npm install`
- **Development Command**: `vercel dev`

### 3. Environment Variables

Add these in Vercel dashboard:

```bash
# AI Providers (Required)
OPENAI_API_KEY=sk-your-openai-key
ANTHROPIC_API_KEY=sk-ant-your-anthropic-key
GROQ_API_KEY=gsk-your-groq-key

# Admin Access (Required)
ADMIN_USERNAME=admin
ADMIN_PASSWORD=your-secure-password

# System Settings
NODE_ENV=production
CONFIG_PATH=/tmp/config.json
```

### 4. Get API Keys

- **OpenAI**: https://platform.openai.com/api-keys
- **Anthropic**: https://console.anthropic.com/
- **Groq**: https://console.groq.com/keys

## 🌐 Live URLs After Deployment

- **Main Site**: `https://your-project.vercel.app/`
- **Admin Panel**: `https://your-project.vercel.app/admin`
- **AI Chat API**: `https://your-project.vercel.app/api/chat`
- **WhatsApp API**: `https://your-project.vercel.app/api/whatsapp`
- **Health Check**: `https://your-project.vercel.app/health`

## 🧪 Testing Your Deployment

### Test Health Endpoint

```bash
curl https://your-project.vercel.app/health
```

### Test AI Chat

```bash
curl -X POST https://your-project.vercel.app/api/chat \
  -H "Content-Type: application/json" \
  -d '{"message":"Hello AI!","provider":"groq"}'
```

### Test Admin Login

```bash
curl -X POST https://your-project.vercel.app/api/chat \
  -H "Content-Type: application/json" \
  -d '{
    "action": "admin_login",
    "credentials": {"username": "admin", "password": "your-password"}
  }'
```

## 🔧 Troubleshooting

### Build Issues

- Ensure all dependencies are installed: `npm install`
- Check Node.js version (18.x recommended)
- Verify environment variables are set

### Runtime Issues

- Check Vercel function logs in dashboard
- Verify API keys are valid
- Test endpoints locally first: `npm run dev`

### Deployment Issues

- Clear build cache in Vercel dashboard
- Check vercel.json configuration
- Ensure GitHub repository is connected properly

## 🎉 Success!

Your PicoClaw WhatsApp AI Integration is now live! 🚀

Visit your admin panel to start configuring your automated WhatsApp flows with AI capabilities.
