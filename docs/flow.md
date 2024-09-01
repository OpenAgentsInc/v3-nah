# Example flow

- User of mobile app records voice command
  - App sends NIP-90 kind 5252 event to connected relay(s)
- Relay receives kind 5252 event, transcribes audio via Whisper/Groq, responds with kind 6252
- App receives 6252, sends NIP-90 kind 5838 event

## Our custom NIP-90 kinds
(Request/Response)
- 5252/6252 - Speech-to-text
- 5838/6838 - Agent command
