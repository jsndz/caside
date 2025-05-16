**Caside**

### 1. **Project Planning**

- **Define the Scope**: Start by outlining the core features Caside will offer, including:

  - Video call support (peer-to-peer or SFU-based)
  - Media recording
  - Uploading videos to cloud storage (like S3)
  - Basic signaling via WebSocket

- **Tech Stack**: Decide on the technologies you'll use for each part of the system:

  - **WebRTC** for video streaming
  - **WebSocket** for signaling
  - **Go or Node.js** for backend services
  - **S3 (or other cloud storage)** for storing video files
  - **JWT** for authentication

### 2. **Architecture Design**

- **System Components**: Break the system into distinct components:

  1. **Signaling Server** (WebSocket server to manage WebRTC connections)
  2. **Media Server** (for managing media streams and recording)
  3. **Storage Integration** (for uploading recorded videos to S3)
  4. **Authentication** (JWT for managing secure access)
  5. **Database** (PostgreSQL or MongoDB for storing metadata like room info, user data, etc.)

- **Interfaces**: Define how these components interact (e.g., WebSocket messages for signaling, REST API for metadata, cloud storage APIs for video uploads).

### 3. **Set Up the Basic Backend Structure**

- **Choose a Framework**: Decide on a backend framework (e.g., Go with `fiber` or Node.js with `express`).
- **Create Microservices**: Structure the backend as independent services:

  - Signaling Service
  - Media Service
  - Storage Service

- **Database Schema**: Design a database schema that can store information like user data, room metadata, and video recording details.

### 4. **Build the Signaling Server**

- **WebSocket Server**: Implement a WebSocket server that will handle signaling between clients for establishing WebRTC connections.

  - Handle messages for joining rooms, sending/receiving SDP offers and ICE candidates, and notifying participants of room status.

- **Session Management**: Maintain a list of active rooms and users in each room.
- **Error Handling**: Make sure that you handle errors gracefully, like disconnecting users or handling failed connections.

### 5. **Implement Media Handling**

- **WebRTC Integration**: Set up WebRTC on the backend to manage video streams. Decide whether to implement:

  - **Peer-to-peer**: Each client connects directly to others.
  - **SFU (Selective Forwarding Unit)**: A server-based model where the server forwards video streams to participants (more scalable).

- **Recording**: Implement functionality for recording the media streams.

  - On the client side, use `MediaRecorder` to capture the media.
  - On the server side, if using SFU, manage recording and storage.

- **Handling Video Data**: Break recorded video data into chunks and prepare it for uploading to cloud storage.

### 6. **Cloud Storage Integration**

- **Set Up Cloud Storage**: Choose a cloud storage provider like **S3**.

  - Set up an S3 bucket for storing recorded videos.
  - Ensure the necessary access rights and configurations are in place (e.g., IAM roles for access).

- **Upload Mechanism**: Implement a system where video files are uploaded in chunks during the call.

  - Handle metadata like video duration, participants, and timestamps.
  - Use a service (e.g., AWS SDK) to upload the video after the call ends.

### 7. **Authentication & Authorization**

- **JWT Authentication**: Implement JWT-based authentication for secure access to your services.

  - Issue JWT tokens when users authenticate and use these tokens to authorize access to APIs (e.g., for starting a video call or accessing the recorded videos).

- **Session Management**: Ensure that each user session is properly authenticated and tracked.

### 8. **REST API for Metadata**

- **API Endpoints**: Set up REST API endpoints for managing metadata:

  - **Create room**: Allow users to create new video call rooms.
  - **Start/Stop recording**: Allow users to start and stop recording.
  - **Get video metadata**: Retrieve information about recorded videos (like duration, participants, etc.).

- **Database Interaction**: Store metadata in your chosen database (e.g., PostgreSQL or MongoDB).

### 9. **Integrate with Frontend**

- **Frontend Connection**: Provide clear documentation on how frontend applications can integrate with Caside.

  - **WebSocket signaling**: Explain how the frontend can connect to your signaling server to initiate a video call.
  - **WebRTC handling**: Provide guidelines on how to set up the WebRTC connection on the client side.
  - **Recording integration**: Guide frontend developers on how to handle the video recording process (via `MediaRecorder` or receiving server-side recordings).

### 10. **Testing and Deployment**

- **Testing**: Thoroughly test all parts of the system:

  - WebSocket connection (signaling)
  - Media handling (WebRTC and recording)
  - Cloud storage uploads
  - Authentication and API security

- **CI/CD Setup**: Set up continuous integration and continuous deployment (CI/CD) pipelines to automatically deploy updates to your backend.
- **Documentation**: Write clear and concise documentation on how to integrate the Caside backend into any frontend application. Include setup instructions, API details, and code examples.

### 11. **Future Improvements**

- **Scalability**: As your application grows, consider improving the scalability of your system. You can achieve this by:

  - Implementing load balancing for WebSocket servers.
  - Moving media handling to dedicated microservices if needed.
  - Using cloud-based storage solutions (e.g., Amazon S3) to scale video storage.

- **Advanced Features**: Once the basic version is up and running, you can add more advanced features like:

  - Video analytics
  - Multi-party recording
  - Transcription or captioning services
  - Enhanced real-time video effects (e.g., background blur).

### 12. **Documentation and Community Support**

- **Create Documentation**: Document every aspect of the project, including:

  - **API endpoints** (REST)
  - **WebSocket message formats**
  - **Integration guidelines for frontend developers**

- **Community Involvement**: If you want to open-source Caside, set up a GitHub repository and encourage other developers to contribute, report issues, or ask for help.

## 🚀 Caside MVP Feature List

### ✅ 1. **Auth (Basic JWT)**

- `POST /login` → returns JWT
- Middleware to protect routes
- Use mock users or hardcoded accounts for now

> **Goal**: Make the service usable by other apps securely.

---

### ✅ 2. **Room Management (In-Memory or DB)**

- `POST /rooms` → create room
- `GET /rooms/:id` → get room details
- Store participants, room state (active/inactive)

> **Goal**: Separate sessions for each video call, like Riverside rooms.

---

### ✅ 3. **WebSocket Signaling Server**

- `ws://.../signal`
- Handles offer/answer, ICE candidate exchange
- Room ID passed on connect
- Supports 1-to-1 or small group signaling

> **Goal**: Bootstrap peer-to-peer connection via WebRTC.

---

### ✅ 4. **Client-side WebRTC + MediaRecorder Integration (Docs Only)**

- Provide a sample frontend or code snippet (not part of backend)
- Let frontend call `MediaRecorder` and `fetch('/upload')` with chunks

> **Goal**: Let devs integrate recording quickly in their frontend.

---

### ✅ 5. **Video Upload API**

- `POST /upload` (auth + multipart or JSON blob)
- Accepts `.webm` chunks or final file
- Stores to:

  - S3 (default)
  - Local FS (fallback option)

> **Goal**: Record locally, upload server-side. Simple, cost-free approach.

---

### ✅ 6. **Recording Metadata**

- `GET /recordings` → list recordings
- `GET /recordings/:id` → get playback URL
- Store filename, size, room ID, timestamp

> **Goal**: Let apps fetch or replay recordings post-call.

---

### ✅ 7. **Dev Setup & Deployment**

- Dockerfile + docker-compose
- `.env` with S3 credentials, JWT secret
- Local dev: `docker compose up`
- Cloud-ready: Deploy on Render, Railway, or VPS

> **Goal**: Devs clone → run → use in under 10 mins.

---

### ✅ 8. **README with Integration Guide**

- Add a **“How to integrate in your app”** section
- Show how to:

  - Connect WebSocket
  - Send SDP offers
  - Record with MediaRecorder
  - Upload video
  - Fetch recording links

> **Goal**: Build trust and usability from day one.

---

## 🧱 Optional (Post-MVP)

- Server-side recording (via SFU or ffmpeg)
- Video stitching + muxing
- User dashboard
- Webhook on recording complete
- Upload progress reporting

---

## 💡 Final Advice

Keep it boring but solid:
👉 Simple REST + WebSocket backend
👉 No SFU or mixing in MVP
👉 Win by making integration easier than setting up Jitsi

---

| Area            | Tech                            | Why                                     |
| --------------- | ------------------------------- | --------------------------------------- |
| Language        | Go (Golang)                     | Speed, concurrency, low memory use      |
| REST API        | [Fiber](https://gofiber.io/)    | Fast, Express-like, easy to use         |
| WebSocket       | Gorilla WebSocket or uWebSocket | Reliable and performant                 |
| Media Upload    | Multipart upload + \[MinIO/S3]  | S3-compatible for video chunk storage   |
| Media Handling  | ffmpeg (invoked via Go)         | For post-processing & muxing            |
| JWT Auth        | `github.com/golang-jwt/jwt/v5`  | Clean JWT parsing and validation        |
| Database        | PostgreSQL + GORM               | Relational, robust, GORM simplifies ORM |
| Configuration   | Viper or Env + `.env`           | Centralized config management           |
| File Storage    | S3, MinIO (dev)                 | Flexible, scalable object storage       |
| Observability   | Zerolog or Zap (logging)        | Structured logs, fast output            |
| Testing         | `stretchr/testify`              | Clean, readable unit tests              |
| Dev Environment | Docker Compose                  | Local orchestration for all services    |

caside/
│
├── cmd/ # Entrypoints for each service
│ ├── signaling/ # WebSocket signaling server
│ ├── api/ # REST API server (room, metadata)
│ └── recorder/ # SFU / media service (optional, later)
│
├── internal/ # Private logic for each service
│ ├── auth/ # JWT validation, user parsing
│ ├── signaling/ # WebRTC signaling logic
│ ├── media/ # Media handling, S3 upload
│ ├── recording/ # Recording logic (e.g., ffmpeg, muxing)
│ └── db/ # PostgreSQL repository layer
│
├── pkg/ # Shared utilities
│ ├── logger/ # Structured logging
│ ├── config/ # Env/Config loading
│ └── utils/ # Misc helpers
│
├── api/ # OpenAPI specs or HTTP handler definitions
│
├── scripts/ # Dev or deploy scripts
│
├── deployments/ # Dockerfiles, k8s, docker-compose
│
├── go.mod
└── README.md

### Stage 2 : Media Handling

[Client A] [Signaling Server] [Client B]
| | |

1.  |----- Create Offer -------->| |
2.  |<------ Send Offer ---------| |
3.  |---------------------------->|----- Receive Offer --------|
4.  | |<----- Create Answer --------|
5.  |<------ Send Answer -------| |
6.  |----- Exchange ICE -------->| |
7.  |<----- Exchange ICE -------| |
8.  | |----- Exchange ICE --------->|
    |<--- Media Flow Starts --->|<--- Media Flow Starts ------>|

What You Should Do Next

1.  Integrate Signaling with Media

    When a client sends an SDP offer:

        Use media.Manager to create or join a session.

        Pass the SDP offer to the Session.HandleSDP method.

        Send the SDP answer back via signaling (Client.Send).

    When receiving ICE candidates:

        Forward these candidates to the relevant Session via Session.AddICECandidate.

    Handle client disconnections by cleaning up sessions and removing clients from rooms.

2.  Implement WebSocket Server Entrypoint

    In cmd/signaling/main.go:

        Set up HTTP server.

        Upgrade HTTP connections to WebSocket.

        Create Client instances for each WS connection.

        Attach clients to the Hub.

        Start client ReadPump and WritePump.

        Initialize media.Manager and signaling.Hub.

3.  Establish Room <-> Session Mapping

    Right now rooms contain clients, and media sessions exist separately.

    Connect a signaling room to a media session by:

        Storing session ID in the room.

        On "join" or "offer", check if a media session exists or create one.

        Relay media signaling between clients via the session.

4.  Track Management & Media Forwarding

    Implement handling of incoming media tracks on PeerConnection.OnTrack event.

    Use your Track abstraction to manage RTP packets.

    Forward media streams as needed to other participants.

5.  Client Side / Frontend

    Build HTML + JavaScript that:

        Connects via WebSocket to signaling server.

        Exchanges SDP offers/answers and ICE candidates.

        Captures local media and streams via WebRTC.

        Displays remote media streams.

6.  Error Handling and Cleanup

    Gracefully handle errors on signaling and media layers.

    Close sessions and tracks properly on client disconnect or room empty.

    Manage resource cleanup.
