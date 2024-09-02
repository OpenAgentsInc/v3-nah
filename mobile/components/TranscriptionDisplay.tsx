import React from 'react';
import { View, Text, StyleSheet } from 'react-native';

interface Message {
  type: 'transcription' | 'agentResponse';
  content: string;
}

interface TranscriptionDisplayProps {
  messages: Message[];
}

const TranscriptionDisplay: React.FC<TranscriptionDisplayProps> = ({ messages }) => {
  return (
    <View style={styles.container}>
      {messages.map((message, index) => (
        <Text key={index} style={styles.message}>
          {message.type === 'transcription' ? 'You: ' : 'Agent: '}
          {message.content}
        </Text>
      ))}
    </View>
  );
};

const styles = StyleSheet.create({
  container: {
    padding: 10,
    pointerEvents: 'box-none',
  },
  message: {
    color: 'white',
    marginBottom: 5,
  },
});

export default TranscriptionDisplay;