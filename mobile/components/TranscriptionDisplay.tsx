import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

interface TranscriptionDisplayProps {
  transcription: string | null;
  agentResponse: string | null;
}

const TranscriptionDisplay: React.FC<TranscriptionDisplayProps> = ({ transcription, agentResponse }) => {
  return (
    <View style={styles.transcriptionContainer}>
      {transcription && (
        <Text style={styles.transcription}>Transcription: {transcription}</Text>
      )}
      {agentResponse && (
        <Text style={styles.agentResponse}>Agent Response: {agentResponse}</Text>
      )}
    </View>
  );
};

const styles = StyleSheet.create({
  transcriptionContainer: {
    flex: 1,
    justifyContent: 'center',
    alignItems: 'center',
    marginBottom: 100,
  },
  transcription: {
    color: '#fff',
    fontFamily: 'JetBrainsMono',
    fontSize: 18,
    textAlign: 'center',
    marginBottom: 10,
  },
  agentResponse: {
    color: '#0f0',
    fontFamily: 'JetBrainsMono',
    fontSize: 18,
    textAlign: 'center',
  },
});

export default TranscriptionDisplay;