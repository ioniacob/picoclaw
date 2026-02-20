import { anthropic } from '@ai-sdk/anthropic';
import { groq } from '@ai-sdk/groq';
import { openai } from '@ai-sdk/openai';
import { generateText, streamText } from 'ai';

// Configuración de providers
const providers = {
  openai: openai('gpt-4-turbo'),
  anthropic: anthropic('claude-3-sonnet-20240229'),
  groq: groq('mixtral-8x7b-32768')
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
      model: providers[provider],
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
  
  // Generar respuesta
  const result = await generateText({
    model: providers[provider],
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
 * Manejar acciones de administrador
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
  
  const conversation = conversationFlows.get(sessionId);
  if (!conversation) return [];
  
  // Retornar últimos 10 mensajes
  return conversation.messages.slice(-10).map(msg => ({
    role: msg.role,
    content: msg.content
  }));
}

/**
 * Guardar mensaje en conversación
 */
function saveToConversation(sessionId, userMessage, aiResponse) {
  if (!sessionId) return;
  
  let conversation = conversationFlows.get(sessionId);
  if (!conversation) {
    conversation = {
      sessionId,
      messages: [],
      createdAt: Date.now(),
      lastActivity: Date.now()
    };
    conversationFlows.set(sessionId, conversation);
  }
  
  conversation.messages.push(
    { role: 'user', content: userMessage, timestamp: Date.now() },
    { role: 'assistant', content: aiResponse, timestamp: Date.now() }
  );
  conversation.lastActivity = Date.now();
  
  // Limpiar conversaciones antiguas (más de 24 horas)
  cleanupOldConversations();
}

/**
 * Limpiar conversaciones antiguas
 */
function cleanupOldConversations() {
  const now = Date.now();
  const maxAge = 24 * 60 * 60 * 1000; // 24 horas
  
  for (const [sessionId, conversation] of conversationFlows.entries()) {
    if (now - conversation.lastActivity > maxAge) {
      conversationFlows.delete(sessionId);
    }
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
        ai_providers: ['openai', 'anthropic', 'groq']
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