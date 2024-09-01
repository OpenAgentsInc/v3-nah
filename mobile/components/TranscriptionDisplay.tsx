import React from 'react';
import { View, Text, StyleSheet, ScrollView } from 'react-native';

interface Message {
  type: 'transcription' | 'agentResponse';
  content: string;
}

interface TranscriptionDisplayProps {
  messages: Message[];
}

const TranscriptionDisplay: React.FC<TranscriptionDisplayProps> = ({ messages }) => {
  return (
    <ScrollView style={styles.scrollView} contentContainerStyle={styles.transcriptionContainer}>
      {messages.map((message, index) => (
        <Text key={index} style={message.type === 'transcription' ? styles.transcription : styles.agentResponse}>
          {message.type === 'transcription' ? 'Transcription: ' : 'Agent Response: '}
          {message.content}
        </Text>
      ))}
    </ScrollView>
  );
};

const styles = StyleSheet.create({
  scrollView: {
    flex: 1,
    width: '100%',
  },
  transcriptionContainer: {
    paddingVertical: 20,
    paddingHorizontal: 10,
    alignItems: 'flex-start',
  },
  transcription: {
    color: '#fff',
    fontFamily: 'JetBrainsMono',
    fontSize: 16,
    marginBottom: 10,
  },
  agentResponse: {
    color: '#0f0',
    fontFamily: 'JetBrainsMono',
    fontSize: 16,
    marginBottom: 20,
  },
});

export default TranscriptionDisplay;