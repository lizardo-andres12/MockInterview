# 🧪 Mock Interview Platform

A browser-based platform for conducting 1-on-1 mock interviews. Users can join a video call, share their screen, use an interview timer, and (optionally) schedule sessions using a calendar.

---

## 📌 Features

- 🎥 Peer-to-peer video and audio calling
- 🖥️ Screen sharing
- ⏱️ Real-time interview countdown timer
- 🧑‍🎨 Shared whiteboard (in progress)
- 📅 Schedule mock interviews via calendar (optional)
- 🧩 Built with modular components and real-time signaling

---

## 💻 Tech Stack

### Frontend
- **Framework**: [React](https://reactjs.org/)
- **Styling**: [Tailwind CSS](https://tailwindcss.com/)
- **Real-time**: WebRTC, WebSockets
- **Libraries**:
  - `simple-peer`
  - `socket.io-client` or native WebSocket
  - `react-canvas-draw` *(planned for whiteboard)*
  - `react-big-calendar` *(planned for calendar)*

### Backend
- **Framework**: [Go](https://fastapi.tiangolo.com/)
- **Server**: Docker
- **Real-time**: WebSocket signaling for peer connections
- **Database**: _(planned — PostgreSQL or SQLite for user data/bookings)_
- **Auth**: _(planned — JWT or OAuth2)_

---

## 🚦 Getting Started

### Backend Setup

cd server
go mod tidy
go run cmd/app/main.go

### Frontend Setup
bash
cd client
npm install
npm run dev
🎯 Usage
Start both frontend and backend servers.

Navigate to:

http://localhost:5173/room123#init (initiator)

http://localhost:5173/room123 (receiver)

Allow camera and mic access.

The video call will start automatically.

🔮 Roadmap
 Video/audio call setup

 Basic room-based WebRTC signaling

 Shared whiteboard drawing

 Chat or notes sidebar

 Timer syncing across clients

 Authentication and protected routes

 Calendar booking and availability system

📃 License
MIT License © 2025 Lizardo Hernandez
