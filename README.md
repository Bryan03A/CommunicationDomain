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
## 9️⃣ **Redis Search History Service** (Node.js / Express / Redis)
- **🧠 Purpose**: Manages user search history by saving and retrieving recent search queries in Redis for quick access and performance optimization.
- **🧪 Port**: `5006`
- **🧰 Tech Stack**:
- &nbsp; - Language: JavaScript (Node.js)
- &nbsp; - Framework: Express.js
- &nbsp; - DB: Redis (in-memory key-value store)
- **🛢️ Database**:
- &nbsp; - Type: In-memory data structure store
- &nbsp; - Engine: Redis, hosted on remote server (`54.161.44.165:6379`)
- &nbsp; - Uses Redis lists (`lPush`, `lRange`) for storing per-user search histories
- **🔐 Security**:
- &nbsp; - Redis connection secured with password authentication (`ADMIN123`)
- &nbsp; - Basic input validation on API endpoints to prevent missing required fields
- **📡 Communication**: REST (JSON)
/ - Endpoints for saving and retrieving user search histories
- **🌍 Endpoints**:
- &nbsp; - `GET /health` — Health check endpoint returning service status
- &nbsp; - `POST /save-search` — Save a search query for a specific user; expects JSON body with `query`, `creator`, `username`, `firstModelName`
- &nbsp; - `GET /recent-searches` — Retrieve last 10 searches for a user by query param `username`
- **🎨 Design Pattern**: Simple REST API with clear separation of concerns; connection and error handling abstracted
- **🏗️ Architecture**: Single-layer service
- &nbsp; - Express routes handle request validation, business logic (save/retrieve), and error handling
- &nbsp; - Redis client manages persistent interaction with remote Redis instance
- **🛠️ Notes**:
- &nbsp; - Uses async/await for asynchronous Redis operations
- &nbsp; - Handles Redis connection lifecycle and errors robustly
- &nbsp; - List-based storage allows efficient insertion and retrieval of recent search data
- &nbsp; - Suitable for caching and fast retrieval use cases with minimal latency
