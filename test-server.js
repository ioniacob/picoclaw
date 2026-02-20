import express from 'express';
import cors from 'cors';
import path from 'path';
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const app = express();
const PORT = process.env.PORT || 3001;

app.use(cors());
app.use(express.json());
app.use(express.static('public'));
app.use('/admin', express.static('admin'));

// Health check endpoint
app.get('/health', (req, res) => {
  res.json({ status: 'ok', timestamp: Date.now() });
});

// Test chat endpoint with admin authentication
app.post('/api/chat', async (req, res) => {
  try {
    const { message, provider = 'openai', action } = req.body;
    
    // Handle admin authentication
    if (action === 'admin_login') {
      const { username, password } = message;
      
      if (username === 'admin' && password === 'picoclaw123') {
        const token = 'admin_token_' + Date.now();
        return res.json({ 
          success: true, 
          token,
          message: 'Login successful'
        });
      } else {
        return res.status(401).json({ error: 'Invalid credentials' });
      }
    }
    
    if (action === 'admin_logout') {
      return res.json({ success: true });
    }
    
    // Regular chat functionality
    if (!message) {
      return res.status(400).json({ error: 'Message required' });
    }

    // Mock response for testing
    const mockResponse = {
      response: `Mock response from ${provider} provider for: "${message}"`,
      provider,
      timestamp: Date.now(),
      sessionId: 'test-session-123'
    };

    res.json(mockResponse);
  } catch (error) {
    console.error('Error in chat:', error);
    res.status(500).json({ error: 'Error processing message' });
  }
});

// Test WhatsApp endpoint
app.post('/api/whatsapp', (req, res) => {
  try {
    const { message, to, action, phone_number_id, access_token, template_name, language_code } = req.body;
    
    // Handle admin login
    if (action === 'admin_login') {
      const { username, password } = req.body.data || {};
      
      // Simple validation (in production use bcrypt)
      if (username === 'admin' && password === 'picoclaw123') {
        const token = 'admin_' + Math.random().toString(36).substr(2, 16);
        
        console.log('🔐 Admin login successful:', username);
        
        return res.json({
          success: true,
          token,
          message: 'Login successful',
          whatsapp_ready: false // Test server doesn't have real WhatsApp
        });
      }
      
      return res.status(401).json({ error: 'Invalid credentials' });
    }
    
    // Handle Facebook WhatsApp Business API actions
    if (action === 'test_facebook_api') {
      if (!phone_number_id || !access_token) {
        return res.status(400).json({ 
          error: 'Facebook WhatsApp Business API requires phone_number_id and access_token',
          required: ['phone_number_id', 'access_token', 'recipient']
        });
      }

      // Validate Facebook API credentials (mock validation)
      if (phone_number_id === 'YOUR_PHONE_NUMBER_ID' || access_token === 'YOUR_ACCESS_TOKEN') {
        return res.status(400).json({ 
          error: 'Please replace placeholder credentials with actual Facebook WhatsApp Business API credentials',
          message: 'Update your .env.local file with real Facebook credentials'
        });
      }

      // Mock Facebook API response
      const mockResponse = {
        success: true,
        provider: 'facebook_whatsapp_business_api',
        messageId: 'fb_msg_' + Date.now(),
        status: 'sent',
        timestamp: Date.now(),
        details: {
          phone_number_id: phone_number_id,
          recipient: to || 'RECIPIENT_PHONE_NUMBER',
          template_name: template_name || 'hello_world',
          language_code: language_code || 'en_US',
          api_version: 'v22.0'
        },
        note: 'This is a mock response. In production, this would send via Facebook Graph API.'
      };

      console.log('📱 Facebook WhatsApp API test:', {
        phone_number_id: phone_number_id,
        recipient: to || 'RECIPIENT_PHONE_NUMBER',
        template: template_name || 'hello_world'
      });

      return res.json(mockResponse);
    }
    
    // Regular WhatsApp message handling
    if (!message || !to) {
      return res.status(400).json({ error: 'Message and destination number required' });
    }

    // Mock WhatsApp response
    const mockResponse = {
      success: true,
      messageId: 'whatsapp_' + Date.now(),
      status: 'sent',
      timestamp: Date.now()
    };

    res.json(mockResponse);
  } catch (error) {
    console.error('Error in WhatsApp:', error);
    res.status(500).json({ error: 'Error processing WhatsApp message' });
  }
});

// Admin login endpoint
app.post('/api/whatsapp/admin/login', (req, res) => {
  try {
    const { action, data } = req.body;
    
    if (action === 'admin_login') {
      const { username, password } = data;
      
      // Simple validation (in production use bcrypt)
      if (username === 'admin' && password === 'picoclaw123') {
        const token = 'admin_' + Math.random().toString(36).substr(2, 16);
        
        console.log('🔐 Admin login successful:', username);
        
        return res.json({
          success: true,
          token,
          message: 'Login successful',
          whatsapp_ready: false // Test server doesn't have real WhatsApp
        });
      }
      
      return res.status(401).json({ error: 'Invalid credentials' });
    }
    
    res.status(400).json({ error: 'Invalid action' });
  } catch (error) {
    console.error('Error in admin login:', error);
    res.status(500).json({ error: 'Error processing login' });
  }
});

// WhatsApp status endpoint with Facebook API support
app.get('/api/whatsapp', (req, res) => {
  const { action } = req.query;
  
  switch (action) {
    case 'status':
      const facebookConfig = {
        enabled: process.env.PICOCLAW_CHANNELS_WHATSAPP_ENABLED === 'true',
        phone_number_id: process.env.PICOCLAW_CHANNELS_WHATSAPP_FB_PHONE_NUMBER_ID,
        api_version: process.env.PICOCLAW_CHANNELS_WHATSAPP_FB_API_VERSION || 'v22.0',
        using_facebook_api: process.env.PICOCLAW_CHANNELS_WHATSAPP_FB_PHONE_NUMBER_ID && 
                           process.env.PICOCLAW_CHANNELS_WHATSAPP_FB_ACCESS_TOKEN
      };
      
      return res.json({
        service: 'PicoClaw WhatsApp',
        version: '2.0.0',
        mode: facebookConfig.using_facebook_api ? 'facebook_business_api' : 'websocket_bridge',
        facebook_config: facebookConfig,
        providers: ['openai', 'anthropic', 'groq'],
        timestamp: Date.now()
      });
      
    case 'facebook_config':
      return res.json({
        required_env_vars: [
          'PICOCLAW_CHANNELS_WHATSAPP_FB_PHONE_NUMBER_ID',
          'PICOCLAW_CHANNELS_WHATSAPP_FB_ACCESS_TOKEN',
          'PICOCLAW_CHANNELS_WHATSAPP_FB_API_VERSION'
        ],
        example_curl: 'curl -X POST http://localhost:3001/api/whatsapp -H "Content-Type: application/json" -d \'{"action": "test_facebook_api", "phone_number_id": "YOUR_PHONE_NUMBER_ID", "access_token": "YOUR_ACCESS_TOKEN", "to": "RECIPIENT_PHONE_NUMBER", "template_name": "hello_world"}\'',
        note: 'Replace placeholder values with actual Facebook WhatsApp Business API credentials'
      });
      
    default:
      return res.json({
        service: 'PicoClaw WhatsApp',
        endpoints: {
          'POST /api/whatsapp': 'Send WhatsApp message',
          'POST /api/whatsapp (action=test_facebook_api)': 'Test Facebook WhatsApp Business API',
          'GET /api/whatsapp?action=status': 'Get service status',
          'GET /api/whatsapp?action=facebook_config': 'Get Facebook API configuration'
        },
        documentation: 'See README_VERCEL_CHAT_SDK.md for Facebook WhatsApp Business API integration'
      });
  }
});

// Start server
app.listen(PORT, () => {
  console.log(`🚀 PicoClaw Test Server running on port ${PORT}`);
  console.log(`📊 Health check: http://localhost:${PORT}/health`);
  console.log(`🤖 Chat API: http://localhost:${PORT}/api/chat`);
  console.log(`📱 WhatsApp API: http://localhost:${PORT}/api/whatsapp`);
  console.log(`🔐 Admin Panel: http://localhost:${PORT}/admin`);
  console.log('');
  console.log('✅ Ready for development!');
});