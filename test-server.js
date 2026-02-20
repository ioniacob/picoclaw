const express = require('express');
const cors = require('cors');
const path = require('path');

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
    const { message, to } = req.body;
    
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