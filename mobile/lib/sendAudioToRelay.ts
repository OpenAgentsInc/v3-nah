import { Audio } from "expo-av"
import * as FileSystem from "expo-file-system"

export async function sendAudioToRelay(audioUri: string, socket: WebSocket, onTranscriptionReceived: (transcription: string) => void): Promise<void> {
  return new Promise(async (resolve, reject) => {
    try {
      const audioContent = await FileSystem.readAsStringAsync(audioUri, {
        encoding: FileSystem.EncodingType.Base64,
      });

      const message = JSON.stringify([
        'EVENT',
        {
          kind: 5252, // NIP-90 range for audio events; we'll use 5252 for speech-to-text
          content: JSON.stringify({
            audio: audioContent,
            format: 'm4a',
          }),
          created_at: Math.floor(Date.now() / 1000),
          tags: [],
        },
      ]);

      socket.send(message);

      // Set up a listener for the 6252 event (speech-to-text response)
      const listener = (event: MessageEvent) => {
        const data = JSON.parse(event.data);
        if (data[0] === 'EVENT' && data[1].kind === 6252) {
          const transcription = JSON.parse(data[1].content).text;
          onTranscriptionReceived(transcription);

          // Send the 5838 event (agent command request)
          const agentCommandMessage = JSON.stringify([
            'EVENT',
            {
              kind: 5838, // NIP-90 kind for agent command request
              content: JSON.stringify({
                command: transcription,
              }),
              created_at: Math.floor(Date.now() / 1000),
              tags: [],
            },
          ]);

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