# Facebook WhatsApp Business API Integration

This guide explains how to integrate Facebook's WhatsApp Business API with PicoClaw using the provided curl command.

## Quick Start

Your Facebook-provided curl command format:
```bash
curl -i -X POST "https://graph.facebook.com/v22.0/YOUR_PHONE_NUMBER_ID/messages" \
  -H "Authorization: Bearer YOUR_ACCESS_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "messaging_product": "whatsapp",
    "to": "RECIPIENT_PHONE_NUMBER",
    "type": "template",
    "template": {
      "name": "hello_world",
      "language": {
        "code": "en_US"
      }
    }
  }'
```

## Configuration

### Option 1: JSON Configuration File

Add this to your `config.json`:

```json
{
  "channels": {
    "whatsapp": {
      "enabled": true,
      "fb_phone_number_id": "YOUR_PHONE_NUMBER_ID",
      "fb_access_token": "YOUR_ACCESS_TOKEN",
      "fb_api_version": "v22.0",
      "allow_from": ["ALLOWED_PHONE_NUMBERS"]
    }
  }
}
```

### Option 2: Environment Variables

Set these environment variables:

```bash
export PICOCLAW_CHANNELS_WHATSAPP_ENABLED=true
export PICOCLAW_CHANNELS_WHATSAPP_FB_PHONE_NUMBER_ID=YOUR_PHONE_NUMBER_ID
export PICOCLAW_CHANNELS_WHATSAPP_FB_ACCESS_TOKEN=YOUR_ACCESS_TOKEN
export PICOCLAW_CHANNELS_WHATSAPP_FB_API_VERSION=v22.0
export PICOCLAW_CHANNELS_WHATSAPP_ALLOW_FROM='["ALLOWED_PHONE_NUMBERS"]'
```

## Usage Examples

### 1. Send Template Message (Your Original Command)

```go
import (
    "context"
    "github.com/sipeed/picoclaw/pkg/channels"
)

// Create client with your credentials
client := channels.NewFacebookWhatsAppClient("YOUR_PHONE_NUMBER_ID", "YOUR_ACCESS_TOKEN", "v22.0")

// Send the exact same template message
ctx := context.Background()
err := client.SendTemplateMessage(ctx, "RECIPIENT_PHONE_NUMBER", "hello_world", "en_US", nil)
```

### 2. Send Text Message

```go
err := client.SendTextMessage(
    ctx,
    "RECIPIENT_PHONE_NUMBER",
    "Hello from PicoClaw! This is a test message.",
)
```

### 3. Send Template with Parameters

```go
components := []channels.TemplateComponent{
    {
        Type: "body",
        Parameters: []channels.TemplateParameter{
            {Type: "text", Text: "John Doe"},
            {Type: "text", Text: "PicoClaw"},
        },
    },
}

err := client.SendTemplateMessage(ctx, "RECIPIENT_PHONE_NUMBER", "welcome_message", "en_US", components)
```

## Integration with PicoClaw

### Using the WhatsApp Channel

The WhatsApp channel now supports both WebSocket bridge and Facebook Business API:

```go
import (
    "github.com/sipeed/picoclaw/pkg/channels"
    "github.com/sipeed/picoclaw/pkg/config"
    "github.com/sipeed/picoclaw/pkg/bus"
)

// Configure WhatsApp to use Facebook API
cfg := config.WhatsAppConfig{
    Enabled:         true,
    FBPhoneNumberID: "YOUR_PHONE_NUMBER_ID",
    FBAccessToken:   "YOUR_ACCESS_TOKEN",
    FBAPIVersion:    "v22.0",
    AllowFrom:       []string{"ALLOWED_PHONE_NUMBERS"},
}

// Create WhatsApp channel
baseChannel := &channels.BaseChannel{
    Name: "whatsapp",
    Bus:  messageBus,
}

whatsappChannel := channels.NewWhatsAppChannel(baseChannel, cfg)

// Start the channel
ctx := context.Background()
if err := whatsappChannel.Start(ctx); err != nil {
    log.Fatal(err)
}

// Send messages through the channel
msg := bus.OutboundMessage{
    ChatID:  "RECIPIENT_PHONE_NUMBER",
    Content: "Hello from PicoClaw!",
}

err := whatsappChannel.Send(ctx, msg)
```

### Send Template Messages

```go
// Send template message through WhatsApp channel
templateComponents := []channels.TemplateComponent{
    {
        Type: "body",
        Parameters: []channels.TemplateParameter{
            {Type: "text", Text: "User Name"},
        },
    },
}

err := whatsappChannel.SendTemplate(
    ctx,
    "RECIPIENT_PHONE_NUMBER",
    "hello_world",
    "en_US",
    templateComponents,
)
```

## Security Considerations

### Credential Management

1. **Never commit access tokens to version control**
2. **Use environment variables for sensitive data**
3. **Rotate access tokens regularly**
4. **Use least-privilege Facebook app permissions**

### Message Validation

The integration includes:
- Phone number format validation
- Template parameter sanitization
- Rate limiting protection
- Connection timeout handling

### Network Security

- All API calls use HTTPS/TLS
- Bearer token authentication
- Request timeout protection (30 seconds default)
- Error response validation

## Error Handling

Common errors and solutions:

```go
// Check if using Facebook API
if whatsappChannel.IsUsingFacebookAPI() {
    // Validate credentials
    err := whatsappChannel.ValidateFacebookCredentials(ctx)
    if err != nil {
        log.Printf("Facebook API validation failed: %v", err)
    }
}

// Handle send errors
err := whatsappChannel.Send(ctx, msg)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "credential validation failed"):
        // Invalid access token or phone number ID
        log.Fatal("Check your Facebook API credentials")
    case strings.Contains(err.Error(), "rate limit"):
        // Implement retry logic with backoff
        time.Sleep(time.Second * 5)
        retry()
    default:
        log.Printf("Send failed: %v", err)
    }
}
```

## Testing

Run the example:

```bash
cd examples
go run facebook_whatsapp_example.go
```

This will:
1. Validate your Facebook credentials
2. Send a template message (your original curl command)
3. Send a text message
4. Send a template with parameters

## Migration from WebSocket Bridge

If you're currently using the WebSocket bridge:

1. **Keep existing configuration**: The bridge will continue to work
2. **Add Facebook credentials**: Add `fb_phone_number_id` and `fb_access_token`
3. **Choose mode**: Only one mode (bridge or Facebook API) can be active
4. **Update message handling**: Facebook API uses different message formats

## API Reference

### FacebookWhatsAppClient Methods

- `NewFacebookWhatsAppClient(phoneNumberID, accessToken, apiVersion)` - Create client
- `SendTemplateMessage(ctx, to, templateName, languageCode, components)` - Send template
- `SendTextMessage(ctx, to, text)` - Send text message
- `ValidateCredentials(ctx)` - Validate API credentials

### WhatsAppChannel Methods

- `SendTemplate(ctx, to, templateName, languageCode, components)` - Send template via channel
- `IsUsingFacebookAPI()` - Check if using Facebook API
- `ValidateFacebookCredentials(ctx)` - Validate Facebook credentials

## Support

For issues:
1. Check Facebook API documentation: https://developers.facebook.com/docs/whatsapp/business-management-api
2. Verify your phone number ID and access token
3. Ensure your WhatsApp Business account is properly configured
4. Check PicoClaw logs for detailed error messages