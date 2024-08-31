import * as FileSystem from 'expo-file-system';

export async function sendAudioToRelay(audioUri: string, socket: WebSocket): Promise<string> {
  return new Promise(async (resolve, reject) => {
    try {
      // Read the audio file
      const audioContent = await FileSystem.readAsStringAsync(audioUri, {
        encoding: FileSystem.EncodingType.Base64,
      });

      // Prepare the message
      const message = JSON.stringify([
        'EVENT',
        {
          kind: 1234, // Custom event kind for audio processing
          content: JSON.stringify({
            audio: audioContent,
            format: 'm4a', // Use m4a format for all platforms
          }),
          created_at: Math.floor(Date.now() / 1000),
          tags: [],
        },
      ]);

      // Set up a one-time event listener for the response
      const messageHandler = (event: MessageEvent) => {
        const response = JSON.parse(event.data);
        if (response[0] === 'EVENT' && response[1].kind === 1235) { // Custom event kind for transcription response
          socket.removeEventListener('message', messageHandler);
          resolve(JSON.parse(response[1].content).transcription);
        }
      };

      socket.addEventListener('message', messageHandler);

      // Send the message
      socket.send(message);

      // Set a timeout for the response
      setTimeout(() => {
        socket.removeEventListener('message', messageHandler);
        reject(new Error('Timeout waiting for transcription'));
      }, 30000); // 30 seconds timeout

    } catch (error) {
      console.error('Error sending audio to relay:', error);
      reject(error);
    }
  });
}