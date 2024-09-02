import "text-encoding-polyfill"
import { StatusBar } from "expo-status-bar"
import React, { useCallback, useState } from "react"
import { SafeAreaView, View, StyleSheet } from "react-native"
import { useStore } from "@/lib/store"
import {
  JetBrainsMono_400Regular, useFonts
} from "@expo-google-fonts/jetbrains-mono"
import Header from "./components/Header"
import PushToTalkButton from "./components/PushToTalkButton"
import TranscriptionDisplay from "./components/TranscriptionDisplay"
import GraphCanvas from "./components/GraphCanvas"
import { sendAudioToRelay } from "./lib/sendAudioToRelay"
import { useAudioPermission } from "./lib/useAudioPermission"
import { useAudioRecording } from "./lib/useAudioRecording"
import { useMessageHandler } from "./lib/useMessageHandler"
import { useNostrUser } from "./lib/useNostrUser"
import { useRelayConnection } from "./lib/useRelayConnection"
import { appStyles } from "./styles/appStyles"

interface Message {
  type: 'transcription' | 'agentResponse';
  content: string;
}

export default function App() {
  useNostrUser()
  let [fontsLoaded] = useFonts({
    JetBrainsMono_400Regular,
  });
  const { isConnected, socket } = useRelayConnection()
  const { startRecording, stopRecording, isRecording } = useAudioRecording()
  const [messages, setMessages] = useState<Message[]>([])
  const [isProcessing, setIsProcessing] = useState(false)
  const permissionStatus = useAudioPermission()

  const addMessage = useCallback((type: 'transcription' | 'agentResponse', content: string) => {
    setMessages(prevMessages => [...prevMessages, { type, content }])
  }, [])

  useMessageHandler({
    socket,
    setAgentResponse: (response) => addMessage('agentResponse', response),
    setIsProcessing
  })

  const handlePressIn = useCallback(async () => {
    if (permissionStatus !== 'granted') {
      console.log('Audio permission not granted')
      return
    }
    await startRecording()
  }, [permissionStatus, startRecording])

  const handlePressOut = useCallback(async () => {
    if (!socket) {
      console.error('No WebSocket connection available')
      return
    }

    setIsProcessing(true)
    try {
      const audioUri = await stopRecording()
      if (audioUri) {
        console.log('Audio recorded:', audioUri)
        await sendAudioToRelay(audioUri, socket, (receivedTranscription) => {
          addMessage('transcription', receivedTranscription)
        })
      }
    } catch (error) {
      console.error('Error sending audio:', error)
      addMessage('transcription', 'Error sending audio')
      setIsProcessing(false)
    }
  }, [socket, stopRecording, addMessage])

  // Sample data for the graph
  const nodes = [
    { id: '1', position: [0, 0, 0] },
    { id: '2', position: [1, 1, 1] },
    { id: '3', position: [-1, -1, -1] },
  ]
  const edges = [
    { source: '1', target: '2' },
    { source: '2', target: '3' },
    { source: '3', target: '1' },
  ]

  return (
    <View style={styles.container}>
      <GraphCanvas nodes={nodes} edges={edges} />
      <SafeAreaView style={styles.content}>
        <Header isConnected={isConnected} />
        <View style={styles.transcriptionContainer}>
          <TranscriptionDisplay messages={messages} />
        </View>
        <PushToTalkButton
          onPressIn={handlePressIn}
          onPressOut={handlePressOut}
          disabled={permissionStatus !== 'granted'}
          isRecording={isRecording}
          isProcessing={isProcessing}
        />
        <StatusBar style="light" />
      </SafeAreaView>
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
  },
  content: {
    flex: 1,
    zIndex: 1,
  },
  transcriptionContainer: {
    flex: 1,
    justifyContent: 'flex-end',
    marginBottom: 20,
  },
});