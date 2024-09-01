# Example flow

- User of mobile app records voice command
  - App sends NIP-90 kind 5252 event to connected relay(s)
- Relay receives kind 5252 event, transcribes audio via Whisper/Groq, responds with kind 6252
- App receives 6252, sends NIP-90 kind 5838 event
- Relay receives kind 5838 event, begins agent command routing flow:
  - Identify what if any tools should be used to handle this request

## Our custom NIP-90 kinds
- 5252 - Speech-to-text request
- 6252 - Speech-to-text response
- 5838 - Agent command request
- 5838 - Agent command response
