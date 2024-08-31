import { Platform } from 'react-native';
import * as FileSystem from 'expo-file-system';

export async function sendAudioToRelay(audioUri: string): Promise<string> {
  const relayUrl = 'https://your-relay-url.com/process-audio'; // Replace with your actual relay URL

  try {
    // Read the audio file
    const audioContent = await FileSystem.readAsStringAsync(audioUri, {
      encoding: FileSystem.EncodingType.Base64,
    });

    // Prepare the request body
    const body = JSON.stringify({
      audio: audioContent,
      format: Platform.OS === 'ios' ? 'caf' : 'webm', // Adjust based on your actual audio format
    });

    // Send the request to the relay
    const response = await fetch(relayUrl, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body,
    });

    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`);
    }

    const result = await response.json();
    return result.transcription; // Assuming the relay returns a transcription

  } catch (error) {
    console.error('Error sending audio to relay:', error);
    throw error;
  }
}