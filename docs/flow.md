# Example flow

- User of mobile app records voice command
  - App sends NIP-90 request via kind 5252 event to connected relay(s)
- Relay receives kind 5252 event, transcribes audio via Whisper/Groq, responds with kind 6252
