# Notice

For those I have discussed this idea with, I want to inform you that I will be pausing any further development on this project to concentrate fully on my internship (thank you Uber!).

Along with this update, I left some architecture diagrams and implementation plans. I intended to use these as a refresher for when I return to the project, but I guess it's also a good opportunity to get some feedback on my design choices. Feel free to reach out with feedback and I will happily get back to you as soon as I can give your suggestions the time it deserves.

I will resume development at full force once my intership concludes on 8/1. If you have any pressing concerns, questions, or suggestions, please reach out to me, but please keep in mind the longer response times.

Enjoy your Summer!

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
