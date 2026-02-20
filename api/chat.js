import { anthropic } from '@ai-sdk/anthropic';
import { groq } from '@ai-sdk/groq';
import { openai } from '@ai-sdk/openai';
import { generateText, streamText } from 'ai';

// Configuración de providers con manejo de errores
const providers = {};

try {
  if (process.env.OPENAI_API_KEY) {
    providers.openai = openai('gpt-4-turbo');
  }
  if (process.env.ANTHROPIC_API_KEY) {
    providers.anthropic = anthropic('claude-3-sonnet-20240229');
  }
  if (process.env.GROQ_API_KEY) {
    providers.groq = groq('mixtral-8x7b-32768');
  }
} catch (error) {
  console.error('Error initializing AI providers:', error);
}

// Proveedor mock para desarrollo
const mockProvider = {
  generateText: async ({ messages }) => {
    const lastMessage = messages[messages.length - 1]?.content || 'Hello';
    return {
      text: `Mock AI response to: "${lastMessage}". This is a development environment response.`
    };
  }
};

// Sistema de autenticación simple
const adminTokens = new Set();
const userSessions = new Map();

// Gestión de flujos de conversación
const conversationFlows = new Map();

/**
 * Handler principal para chat con AI
 */
export async function POST(request) {
  try {
    const { message, provider = 'openai', sessionId, action } = await request.json();
    
    if (!message) {
      return new Response(JSON.stringify({ error: 'Mensaje requerido' }), {
        status: 400,
        headers: { 'Content-Type': 'application/json' }
      });
    }

    // Verificar autenticación para acciones de admin
    if (action === 'admin_login' || action === 'admin_logout') {
      return handleAdminAction(action, message);
    }

    // Procesar mensaje normal
    const response = await processChatMessage(message, provider, sessionId);
    
    return new Response(JSON.stringify({ 
      response, 
      provider,
      timestamp: Date.now(),
      sessionId: sessionId || generateSessionId()
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  } catch (error) {
    console.error('Error en chat:', error);
    return new Response(JSON.stringify({ error: 'Error procesando mensaje' }), {
      status: 500,
      headers: { 'Content-Type': 'application/json' }
    });
  }
}

/**
 * Handler para streaming de respuestas
 */
export async function GET(request) {
  const { searchParams } = new URL(request.url);
  const message = searchParams.get('message');
  const provider = searchParams.get('provider') || 'openai';
  
  if (!message) {
    return new Response('Mensaje requerido', { status: 400 });
  }

  try {
    const result = await streamText({
      model: providers[provider] || mockProvider,
      messages: [{ role: 'user', content: message }],
      maxTokens: 1000,
      temperature: 0.7,
    });

    return result.toDataStreamResponse();
  } catch (error) {
    console.error('Error en streaming:', error);
    return new Response('Error en streaming', { status: 500 });
  }
}

/**
 * Procesar mensaje de chat con AI
 */
async function processChatMessage(message, provider, sessionId) {
  // Obtener contexto de la conversación
  const context = getConversationContext(sessionId);
  
  // Usar provider disponible o mock
  const aiProvider = providers[provider] || mockProvider;
  
  // Generar respuesta
  const result = await generateText({
    model: aiProvider,
    messages: [
      { role: 'system', content: getSystemPrompt() },
      ...context,
      { role: 'user', content: message }
    ],
    maxTokens: 1000,
    temperature: 0.7,
  });

  // Guardar en el contexto
  saveToConversation(sessionId, message, result.text);

  return result.text;
}

/**
 * Handler para acciones de administrador
 */
function handleAdminAction(action, credentials) {
  if (action === 'admin_login') {
    const { username, password } = credentials;
    
    // Validación simple (en producción usar bcrypt)
    if (username === process.env.ADMIN_USERNAME && 
        password === process.env.ADMIN_PASSWORD) {
      const token = generateAdminToken();
      adminTokens.add(token);
      
      return new Response(JSON.stringify({ 
        success: true, 
        token,
        message: 'Login exitoso'
      }), {
        headers: { 'Content-Type': 'application/json' }
      });
    }
    
    return new Response(JSON.stringify({ error: 'Credenciales inválidas' }), {
      status: 401,
      headers: { 'Content-Type': 'application/json' }
    });
  }
  
  if (action === 'admin_logout') {
    const { token } = credentials;
    adminTokens.delete(token);
    
    return new Response(JSON.stringify({ success: true }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }
}

/**
 * Obtener contexto de conversación
 */
function getConversationContext(sessionId) {
  if (!sessionId) return [];
  
  const session = userSessions.get(sessionId);
  if (!session) return [];
  
  // Retornar últimos 10 mensajes
  return session.messages.slice(-10).map(msg => ({
    role: msg.role,
    content: msg.content
  }));
}

/**
 * Guardar en conversación
 */
function saveToConversation(sessionId, userMessage, aiResponse) {
  if (!sessionId) return;
  
  if (!userSessions.has(sessionId)) {
    userSessions.set(sessionId, {
      messages: [],
      createdAt: Date.now()
    });
  }
  
  const session = userSessions.get(sessionId);
  session.messages.push(
    { role: 'user', content: userMessage, timestamp: Date.now() },
    { role: 'assistant', content: aiResponse, timestamp: Date.now() }
  );
  
  // Limpiar mensajes antiguos (mantener últimos 50)
  if (session.messages.length > 50) {
    session.messages = session.messages.slice(-50);
  }
}

/**
 * Generar ID de sesión
 */
function generateSessionId() {
  return 'session_' + Math.random().toString(36).substr(2, 9);
}

/**
 * Generar token de administrador
 */
function generateAdminToken() {
  return 'admin_' + Math.random().toString(36).substr(2, 16);
}

/**
 * Handler GET para health checks y root path
 */
export async function GET(request) {
  const url = new URL(request.url);
  
  // Health check endpoint
  if (url.pathname === '/health' || url.pathname === '/') {
    return new Response(JSON.stringify({
      status: 'healthy',
      service: 'PicoClaw WhatsApp AI Integration',
      version: '1.0.0',
      timestamp: new Date().toISOString(),
      features: {
        chat: true,
        whatsapp: true,
        admin: true,
        ai_providers: Object.keys(providers),
        mock_mode: Object.keys(providers).length === 0
      }
    }), {
      headers: { 'Content-Type': 'application/json' }
    });
  }
  
  // Default response for unknown GET paths
  return new Response(JSON.stringify({
    message: 'PicoClaw WhatsApp AI Integration API',
    endpoints: {
      chat: '/api/chat',
      whatsapp: '/api/whatsapp', 
      admin: '/admin',
      health: '/health'
    }
  }), {
    headers: { 'Content-Type': 'application/json' }
  });
}

/**
 * Obtener prompt del sistema
 */
function getSystemPrompt() {
  return `Eres un asistente AI inteligente y útil. 
  Responde de manera concisa y clara.
  Si te preguntan sobre WhatsApp, explica que puedes ayudar a configurar flujos automáticos.
  Mantén un tono profesional pero amigable.`;
}