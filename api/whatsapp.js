import { Client, LocalAuth } from 'whatsapp-web.js';
import { generateText } from 'ai';
import { openai } from '@ai-sdk/openai';
import { anthropic } from '@ai-sdk/anthropic';
import { groq } from '@ai-sdk/groq';

// Configuración de providers para WhatsApp
const providers = {
  openai: openai('gpt-4-turbo'),
  anthropic: anthropic('claude-3-sonnet-20240229'),
  groq: groq('mixtral-8x7b-32768')
};

// Estado de WhatsApp
let whatsappClient = null;
let isWhatsAppReady = false;
let adminUsers = new Set();

/**
 * Handler principal para WhatsApp con AI
 */
export default async function handler(req, res) {
  // Configurar CORS
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization');

  if (req.method === 'OPTIONS') {
    return res.status(200).end();
  }

  try {
    switch (req.method) {
      case 'GET':
        return handleGet(req, res);
      case 'POST':
        return handlePost(req, res);
      default:
        return res.status(405).json({ error: 'Método no permitido' });
    }
  } catch (error) {
    console.error('Error en handler:', error);
    return res.status(500).json({ error: 'Error interno del servidor' });
  }
}

/**
 * Manejar peticiones GET
 */
async function handleGet(req, res) {
  const { action } = req.query;

  switch (action) {
    case 'status':
      return res.status(200).json({
        service: 'PicoClaw WhatsApp AI',
        version: '1.0.0',
        whatsapp: {
          ready: isWhatsAppReady,
          connected: whatsappClient?.info ? true : false
        },
        providers: Object.keys(providers),
        timestamp: Date.now()
      });

    case 'qr':
      if (!isWhatsAppReady) {
        return res.status(200).json({ 
          status: 'waiting_for_qr',
          message: 'WhatsApp no está listo, esperando QR code'
        });
      }
      return res.status(200).json({ 
        status: 'ready',
        message: 'WhatsApp está conectado'
      });

    default:
      return res.status(200).json({
        service: 'PicoClaw WhatsApp AI',
        endpoints: {
          'GET /api/whatsapp?action=status': 'Estado del servicio',
          'GET /api/whatsapp?action=qr': 'Estado de QR',
          'POST /api/whatsapp': 'Webhook para mensajes',
          'POST /api/whatsapp/admin/login': 'Login de administrador'
        },
        documentation: 'Ver README.md para más información'
      });
  }
}

/**
 * Manejar peticiones POST
 */
async function handlePost(req, res) {
  const { action, data } = req.body;

  switch (action) {
    case 'admin_login':
      return handleAdminLogin(data, res);
    
    case 'admin_logout':
      return handleAdminLogout(data, res);
    
    case 'send_message':
      return handleSendMessage(data, res);
    
    case 'webhook':
      return handleWebhook(data, res);
    
    default:
      return res.status(400).json({ error: 'Acción no válida' });
  }
}

/**
 * Login de administrador
 */
async function handleAdminLogin(data, res) {
  const { username, password } = data;
  
  // Validación simple (en producción usar bcrypt)
  if (username === process.env.ADMIN_USERNAME && 
      password === process.env.ADMIN_PASSWORD) {
    
    const token = generateToken();
    adminUsers.add(token);
    
    // Iniciar WhatsApp si no está iniciado
    if (!whatsappClient) {
      await initializeWhatsApp();
    }
    
    return res.status(200).json({
      success: true,
      token,
      message: 'Login exitoso',
      whatsapp_ready: isWhatsAppReady
    });
  }
  
  return res.status(401).json({ error: 'Credenciales inválidas' });
}

/**
 * Logout de administrador
 */
async function handleAdminLogout(data, res) {
  const { token } = data;
  adminUsers.delete(token);
  
  return res.status(200).json({ success: true });
}

/**
 * Enviar mensaje a través de WhatsApp
 */
async function handleSendMessage(data, res) {
  const { to, message, provider = 'openai' } = data;
  
  if (!isWhatsAppReady) {
    return res.status(503).json({ error: 'WhatsApp no está conectado' });
  }
  
  try {
    // Si el mensaje contiene una pregunta para AI, procesarla
    if (message.includes('@ai') || message.includes('/ai')) {
      const aiResponse = await processWithAI(message.replace('@ai', '').replace('/ai', ''), provider);
      await whatsappClient.sendMessage(to, aiResponse);
    } else {
      await whatsappClient.sendMessage(to, message);
    }
    
    return res.status(200).json({ success: true });
  } catch (error) {
    console.error('Error enviando mensaje:', error);
    return res.status(500).json({ error: 'Error enviando mensaje' });
  }
}

/**
 * Procesar mensaje con AI
 */
async function processWithAI(message, providerName = 'openai') {
  const provider = providers[providerName] || providers.openai;
  
  try {
    const result = await generateText({
      model: provider,
      messages: [
        { role: 'system', content: getSystemPrompt() },
        { role: 'user', content: message }
      ],
      maxTokens: 500,
      temperature: 0.7,
    });
    
    return result.text;
  } catch (error) {
    console.error('Error con AI:', error);
    return 'Lo siento, hubo un error procesando tu mensaje con AI.';
  }
}

/**
 * Inicializar WhatsApp
 */
async function initializeWhatsApp() {
  whatsappClient = new Client({
    authStrategy: new LocalAuth(),
    puppeteer: {
      headless: true,
      args: [
        '--no-sandbox',
        '--disable-setuid-sandbox',
        '--disable-dev-shm-usage',
        '--disable-accelerated-2d-canvas',
        '--no-first-run',
        '--no-zygote',
        '--single-process',
        '--disable-gpu'
      ]
    }
  });

  whatsappClient.on('qr', (qr) => {
    console.log('QR Code generado, escanea con WhatsApp');
    // En un entorno real, podrías enviar este QR a través de WebSocket
  });

  whatsappClient.on('ready', () => {
    console.log('WhatsApp está listo!');
    isWhatsAppReady = true;
  });

  whatsappClient.on('message', async (message) => {
    // No responder a mensajes de grupos o mensajes del propio bot
    if (message.fromMe || message.from.includes('@g.us')) {
      return;
    }
    
    console.log(`Mensaje recibido de ${message.from}: ${message.body}`);
    
    // Procesar mensajes que mencionan AI
    if (message.body.includes('@ai') || message.body.includes('/ai')) {
      const aiResponse = await processWithAI(
        message.body.replace('@ai', '').replace('/ai', ''),
        'openai' // Por defecto usar OpenAI
      );
      
      await whatsappClient.sendMessage(message.from, aiResponse);
    }
  });

  whatsappClient.on('disconnected', (reason) => {
    console.log('WhatsApp desconectado:', reason);
    isWhatsAppReady = false;
  });

  await whatsappClient.initialize();
}

/**
 * Obtener prompt del sistema
 */
function getSystemPrompt() {
  return `Eres un asistente AI inteligente que responde a través de WhatsApp.
  Responde de manera concisa y clara, adecuada para mensajes de WhatsApp.
  Usa emojis cuando sea apropiado para hacer la conversación más amigable.
  Mantén un tono profesional pero amigable.`;
}

/**
 * Generar token
 */
function generateToken() {
  return 'admin_' + Math.random().toString(36).substr(2, 16);
}