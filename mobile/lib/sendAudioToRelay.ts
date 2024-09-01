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
        content: "",
        created_at: Math.floor(Date.now() / 1000),
        tags: [
          ["i", audioContent, "text"],
          ["param", "format", "m4a"],
          ["output", "text/plain"],
          ["bid", "0"]
        ],
      };

      const message = JSON.stringify(["EVENT", event]);

      socket.send(message);

      // Set up a listener for the 6252 event (speech-to-text response)
      const listener = (event: MessageEvent) => {
        const data = JSON.parse(event.data);
        console.log("Received data:", data);
        if (Array.isArray(data) && data[0] === "EVENT") {
          const eventData = data[1];
          if (eventData.kind === 6252) {
            const transcription = eventData.content;
            onTranscriptionReceived(transcription);

            // Send the 5838 event (agent command request)
            const agentCommandEvent = {
              kind: 5838, // NIP-90 kind for agent command request
              content: "",
              created_at: Math.floor(Date.now() / 1000),
              tags: [
                ["i", transcription, "text"],
                ["output", "text/plain"],
                ["bid", "0"],
                ["t", "agent_command"],
                ["param", "repo", "https://github.com/OpenAgentsInc/v3"]
              ],
            };

            const agentCommandMessage = JSON.stringify(["EVENT", agentCommandEvent]);

            socket.send(agentCommandMessage);

            // Remove the listener after processing the 6252 event
            socket.removeEventListener('message', listener);
            resolve();
          }
        }
      };

      socket.addEventListener('message', listener);
    } catch (error) {
      console.error('Error sending audio to relay:', error);
      reject(error);
    }
  });
}
