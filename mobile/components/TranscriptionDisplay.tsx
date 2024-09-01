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
          {message.content.trim()}
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
    color: 'rgba(255, 255, 255, 0.75)', // White with 75% opacity
    fontFamily: 'JetBrainsMono',
    fontSize: 16,
    marginBottom: 10,
  },
  agentResponse: {
    color: '#fff', // Solid white
    fontFamily: 'JetBrainsMono',
    fontSize: 16,
    marginBottom: 20,
  },
});

export default TranscriptionDisplay;