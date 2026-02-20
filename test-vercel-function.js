// Test the WhatsApp function locally
import('./api/whatsapp.js').then(module => {
  console.log('✅ WhatsApp module loaded successfully!');
  console.log('Available exports:', Object.keys(module));
}).catch(error => {
  console.error('❌ Error loading WhatsApp module:');
  console.error(error.message);
});