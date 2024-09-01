import { Audio } from "expo-av"
import * as FileSystem from "expo-file-system"

export async function sendAudioToRelay(audioUri: string, socket: WebSocket): Promise<void> {
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

      // Instead of waiting for a response here, we'll resolve immediately
      // The App component will handle the response through its WebSocket listener
      resolve();
    } catch (error) {
      console.error('Error sending audio to relay:', error);
      reject(error);
    }
  });
}
