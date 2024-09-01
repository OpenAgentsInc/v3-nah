import { Audio } from "expo-av"
import * as FileSystem from "expo-file-system"

export async function sendAudioToRelay(audioUri: string, socket: WebSocket, onTranscriptionReceived: (transcription: string) => void): Promise<void> {
  return new Promise(async (resolve, reject) => {
    try {
      const audioContent = await FileSystem.readAsStringAsync(audioUri, {
        encoding: FileSystem.EncodingType.Base64,
      });

      const event = {
        kind: 5252, // NIP-90 range for audio events; we'll use 5252 for speech-to-text
        content: JSON.stringify({
          audio: audioContent,
          format: 'm4a',
        }),
        created_at: Math.floor(Date.now() / 1000),
        tags: [],
      };

      const message = JSON.stringify(["EVENT", event]);

      socket.send(message);

      // Set up a listener for the 6252 event (speech-to-text response)
      const listener = (event: MessageEvent) => {
        const data = JSON.parse(event.data);
        console.log("Received data:", data);
        if (data.type === 'EVENT' && data.data.kind === 6252) {
          const transcription = data.data.content;
          onTranscriptionReceived(transcription);

          // Send the 5838 event (agent command request)
          const agentCommandEvent = {
            kind: 5838, // NIP-90 kind for agent command request
            content: JSON.stringify({
              command: transcription,
            }),
            created_at: Math.floor(Date.now() / 1000),
            tags: [],
          };

          const agentCommandMessage = JSON.stringify(["EVENT", agentCommandEvent]);

          socket.send(agentCommandMessage);

          // Remove the listener after processing the 6252 event
          socket.removeEventListener('message', listener);
          resolve();
        }
      };

      socket.addEventListener('message', listener);
    } catch (error) {
      console.error('Error sending audio to relay:', error);
      reject(error);
    }
  });
}