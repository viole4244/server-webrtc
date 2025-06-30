# 🚀 WebRTC Chat & File Share

A real-time chat and file sharing application built with Go (WebRTC signaling server) and vanilla JavaScript (client). This project demonstrates a hybrid WebRTC architecture where the server acts as a relay between peers using WebRTC DataChannels.

## ✨ Features

- 🏠 **Room-based Communication**: Create or join rooms with custom room IDs
- 💬 **Real-time Chat**: Instant messaging between connected peers
- 📁 **File Sharing**: Transfer files of any size with automatic chunking (64KB chunks)
- 🔄 **Server-mediated Relay**: Uses WebRTC DataChannels with Go server as intermediary
- 🎨 **Modern UI**: Clean, dark-themed interface with responsive design
- ⚡ **Binary & Text Support**: Handles both chat messages and binary file data

## 🏗️ Architecture

```
Client A ←→ WebSocket ←→ Go Server ←→ WebRTC DataChannel ←→ Go Server ←→ WebSocket ←→ Client B
```

This application uses a **hybrid WebRTC approach**:
- WebRTC is used for establishing reliable DataChannels
- The Go server acts as a relay (not true peer-to-peer)
- All data flows through the server for simplified NAT traversal and control

## 🛠️ Tech Stack

**Backend:**
- Go 1.19+
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - WebSocket handling
- [Pion WebRTC](https://github.com/pion/webrtc) - WebRTC implementation

**Frontend:**
- Vanilla JavaScript (ES6+)
- HTML5 File API
- WebSocket API
- Modern CSS with Flexbox

## 🚀 Quick Start

### Prerequisites

- Go 1.19 or higher
- Modern web browser with WebRTC support

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/yourusername/webrtc-chat-file-share.git
   cd webrtc-chat-file-share
   ```

2. **Initialize Go module**
   ```bash
   go mod init webrtc-chat
   go mod tidy
   ```

3. **Install dependencies**
   ```bash
   go get github.com/gorilla/websocket
   go get github.com/pion/webrtc/v3
   ```

4. **Create static directory**
   ```bash
   mkdir static
   # Move the HTML file to static/index.html
   ```

5. **Run the server**
   ```bash
   go run main.go
   ```

6. **Open your browser**
   Navigate to `http://localhost:8080`

### Project Structure

```
webrtc-chat-file-share/
├── main.go              # Go WebRTC signaling server
├── static/
│   └── index.html       # Client application
├── go.mod
├── go.sum
└── README.md
```

## 📖 Usage

### Creating a Room

1. Enter a unique room ID in the input field
2. Click **"Create"** button
3. Wait for another peer to join
4. Start chatting and sharing files!

### Joining a Room

1. Enter the room ID shared by the room creator
2. Click **"Join"** button
3. Connection will be established automatically
4. Begin communication!

### File Sharing

1. Click the **"Share File"** button
2. Select any file from your device
3. File will be automatically chunked and transferred
4. Recipient will see a download link when transfer completes

## 🔧 How It Works

### Connection Flow

1. **WebSocket Connection**: Clients connect to Go server via WebSocket
2. **Room Management**: Server creates/manages room structs for peer coordination
3. **WebRTC Handshake**: Server facilitates SDP offer/answer exchange and ICE gathering
4. **DataChannel Setup**: WebRTC DataChannels established between server and each client
5. **Message Relay**: Server forwards all messages between peers through DataChannels

### Message Types

| Type | Protocol | Description |
|------|----------|-------------|
| `create` | WebSocket JSON | Create a new room |
| `join` | WebSocket JSON | Join existing room |
| `chat` | WebSocket JSON → DataChannel | Text messages |
| `file-meta` | WebSocket JSON → DataChannel | File metadata (name, size, type) |
| File chunks | WebSocket Binary → DataChannel | 64KB binary chunks |

### File Transfer Process

1. **Metadata**: File info sent as JSON message
2. **Chunking**: File split into 64KB chunks using FileReader API
3. **Binary Transfer**: Each chunk sent as binary WebSocket message
4. **Relay**: Server forwards chunks through WebRTC DataChannels
5. **Reconstruction**: Receiving client reassembles chunks into downloadable file

## ⚙️ Configuration

### Server Configuration

```go
// Modify these constants in main.go
const (
    PORT = ":8080"                    // Server port
    STUN_SERVER = "stun:stun.l.google.com:19302"  // STUN server for ICE
)
```

### Client Configuration

```javascript
// Modify these constants in index.html
const CHUNK_SIZE = 64 * 1024;  // File chunk size (64KB)
```

## 🧪 Development

### Running in Development Mode

```bash
# Run with auto-reload (if using air)
air

# Or run normally
go run main.go
```

### Testing File Transfer

1. Open two browser tabs/windows
2. Create a room in one, join in the other
3. Try sharing different file types and sizes
4. Monitor browser console for debugging info

## 🔒 Security Considerations

⚠️ **This is a development/demo application. For production use, consider:**

- Authentication and authorization
- Rate limiting
- File type and size restrictions
- Input validation and sanitization
- HTTPS/WSS encryption
- CORS configuration
- Room cleanup and expiration

## 🐛 Troubleshooting

### Common Issues

**Connection fails:**
- Check if port 8080 is available
- Ensure WebRTC is supported in your browser
- Try disabling browser extensions

**File transfer stuck:**
- Check browser console for errors
- Try smaller files first
- Ensure stable network connection

**WebRTC errors:**
- STUN server might be unreachable
- Try different STUN servers in the configuration

## 📝 API Reference

### WebSocket Message Format

```json
{
  "type": "message_type",
  "payload": "message_payload",
  "roomId": "room_identifier"
}
```

### Supported Message Types

- `create`: Create new room
- `join`: Join existing room  
- `chat`: Send chat message
- `file-meta`: Send file metadata
- `status`: Server status updates
- `error`: Error messages

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Pion WebRTC](https://github.com/pion/webrtc) - Excellent Go WebRTC implementation
- [Gorilla WebSocket](https://github.com/gorilla/websocket) - Reliable WebSocket library
- [Google STUN servers](https://developers.google.com/web/fundamentals/connectivity/webrtc) - Free STUN services

## 📚 Learn More

- [WebRTC Fundamentals](https://webrtc.org/getting-started/overview)
- [Pion WebRTC Examples](https://github.com/pion/webrtc/tree/master/examples)
- [Go WebSocket Tutorial](https://gorilla.github.io/websocket/)

---

**Built with ❤️ using Go and WebRTC**
