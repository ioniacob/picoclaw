// Bridge WebSocket para WhatsApp - Ejemplo para Vercel
// Este es un ejemplo de cómo crear un bridge que conecte WhatsApp con PicoClaw

const { Client, LocalAuth } = require('whatsapp-web.js');
const WebSocket = require('ws');
const express = require('express');

const app = express();
const PORT = process.env.PORT || 3001;

// Configuración de WhatsApp
const client = new Client({
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

// WebSocket Server
let wss;
let wsClients = new Set();

function setupWebSocketServer() {
    wss = new WebSocket.Server({ port: PORT });
    
    wss.on('connection', (ws) => {
        console.log('Nuevo cliente WebSocket conectado');
        wsClients.add(ws);
        
        ws.on('message', (data) => {
            try {
                const message = JSON.parse(data);
                handleOutgoingMessage(message);
            } catch (error) {
                console.error('Error procesando mensaje:', error);
            }
        });
        
        ws.on('close', () => {
            console.log('Cliente WebSocket desconectado');
            wsClients.delete(ws);
        });
        
        // Enviar mensaje de bienvenida
        ws.send(JSON.stringify({
            type: 'status',
            message: 'Conectado a WhatsApp Bridge',
            timestamp: Date.now()
        }));
    });
}

// Manejar mensajes salientes (de PicoClaw a WhatsApp)
async function handleOutgoingMessage(message) {
    if (message.type === 'message' && message.to && message.content) {
        try {
            await client.sendMessage(message.to, message.content);
            console.log(`Mensaje enviado a ${message.to}: ${message.content}`);
        } catch (error) {
            console.error('Error enviando mensaje:', error);
            // Notificar error a los clientes
            broadcastToClients({
                type: 'error',
                message: 'Error enviando mensaje',
                error: error.message
            });
        }
    }
}

// Enviar mensajes a todos los clientes conectados
function broadcastToClients(message) {
    const messageStr = JSON.stringify(message);
    wsClients.forEach(ws => {
        if (ws.readyState === WebSocket.OPEN) {
            ws.send(messageStr);
        }
    });
}

// Eventos de WhatsApp
client.on('qr', (qr) => {
    console.log('QR Code generado, escanea con WhatsApp:');
    console.log(qr);
    
    // Enviar QR a los clientes si es necesario
    broadcastToClients({
        type: 'qr',
        data: qr
    });
});

client.on('ready', () => {
    console.log('WhatsApp está listo!');
    broadcastToClients({
        type: 'status',
        message: 'WhatsApp conectado'
    });
});

client.on('message', async (message) => {
    // No responder a mensajes de grupos o mensajes del propio bot
    if (message.fromMe || message.from.includes('@g.us')) {
        return;
    }
    
    console.log(`Mensaje recibido de ${message.from}: ${message.body}`);
    
    // Formatear mensaje para PicoClaw
    const formattedMessage = {
        type: 'message',
        from: message.from,
        chat: message.from,
        content: message.body,
        id: message.id.id,
        timestamp: message.timestamp
    };
    
    // Enviar a todos los clientes WebSocket
    broadcastToClients(formattedMessage);
});

client.on('disconnected', (reason) => {
    console.log('WhatsApp desconectado:', reason);
    broadcastToClients({
        type: 'status',
        message: 'WhatsApp desconectado',
        reason: reason
    });
});

// Health check endpoint
app.get('/health', (req, res) => {
    res.json({
        status: 'healthy',
        service: 'whatsapp-bridge',
        timestamp: Date.now(),
        clients: wsClients.size
    });
});

// Información del servicio
app.get('/', (req, res) => {
    res.json({
        service: 'WhatsApp Bridge for PicoClaw',
        version: '1.0.0',
        status: client.info ? 'connected' : 'disconnected',
        clients: wsClients.size,
        endpoints: {
            health: '/health',
            websocket: `ws://localhost:${PORT}`
        }
    });
});

// Iniciar servicios
async function start() {
    try {
        setupWebSocketServer();
        await client.initialize();
        
        console.log(`🚀 WhatsApp Bridge iniciado`);
        console.log(`📡 WebSocket server en ws://localhost:${PORT}`);
        console.log(`🏥 Health check en http://localhost:${PORT}/health`);
        
    } catch (error) {
        console.error('Error iniciando servicios:', error);
        process.exit(1);
    }
}

// Manejo de errores
process.on('uncaughtException', (error) => {
    console.error('Uncaught Exception:', error);
    process.exit(1);
});

process.on('unhandledRejection', (reason, promise) => {
    console.error('Unhandled Rejection at:', promise, 'reason:', reason);
    process.exit(1);
});

// Iniciar
start();

// Exportar para Vercel
module.exports = app;