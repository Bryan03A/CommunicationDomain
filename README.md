## 8️⃣ **Chat Service** (Python / FastAPI + WebSocket + MongoDB)
- **🧠 Purpose**: Provides real-time bidirectional messaging between two users via WebSocket, with message persistence in MongoDB and user validation via an external SOAP service.
- **🧪 Port**: `5010`
- **🧰 Tech Stack**:
  - Language: Python
  - Framework: FastAPI
  - DB: MongoDB (via Motor - async driver)
  - Realtime: WebSocket
  - Env Config: `dotenv`
  - External Call: `requests` to a SOAP microservice
- **🛢️ Database**:
  - Type: NoSQL
  - Engine: MongoDB
  - Collection: `chats`
  - Stored Data:
    - `chat_id`: Unique conversation ID (based on usernames)
    - `messages`: List of `{ sender, text }` dictionaries
- **🔐 Security**:
  - User existence validated against external `user-search-service` (SOAP) before WebSocket is established
  - Chat IDs are deterministic to prevent duplication or user spoofing
  - `.env`-based secret loading (e.g., `MONGO_URI`, `SECRET_KEY`)
- **📡 Communication**: WebSocket (JSON)
  - URL format: `ws://host/ws/{user1}/{user2}`
  - Sends:
    - Historical messages upon connection
    - Real-time messages received from peers
  - Receives:
    - JSON payloads like `{ "sender": "alice", "text": "Hi!" }`
- **📂 Endpoints**:
  - `GET /ws/{user1}/{user2}` — Initiates a WebSocket chat session between two users
- **♻️ Message Flow**:
  1. Verifies both users exist via HTTP call to `user-search-service`
  2. Accepts WebSocket and assigns connection to a `chat_id`
  3. Loads previous messages from MongoDB (if any)
  4. Listens for messages, stores them, and broadcasts to connected clients
  5. Handles disconnection gracefully and cleans up unused connections
- **🏗️ Architecture**: Single-layer with separation of:
  - Connection management (`active_chats`)
  - External service validation (`user_exists`)
  - Persistence (MongoDB `update_one`, `find_one`)
- **🛠️ Notes**:
  - MongoDB connection is verified at startup with a `ping`
  - Ensures chat history is not lost between restarts
  - Uses consistent chat IDs (alphabetically sorted usernames)
  - Avoids echoing sender’s message back to themselves
---
