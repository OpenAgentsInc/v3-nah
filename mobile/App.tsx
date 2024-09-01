import "text-encoding-polyfill"
import { StatusBar } from "expo-status-bar"
import React, { useCallback, useState } from "react"
import { SafeAreaView, View } from "react-native"
import { useStore } from "@/lib/store"
import PushToTalkButton from "./components/PushToTalkButton"
import Header from "./components/Header"
import TranscriptionDisplay from "./components/TranscriptionDisplay"
import { sendAudioToRelay } from "./lib/sendAudioToRelay"
import { useAudioRecording } from "./lib/useAudioRecording"
import { useNostrUser } from "./lib/useNostrUser"
import { useRelayConnection } from "./lib/useRelayConnection"
import { useAudioPermission } from "./lib/useAudioPermission"
import { useMessageHandler } from "./lib/useMessageHandler"
import { appStyles } from "./styles/appStyles"

interface Message {
  type: 'transcription' | 'agentResponse';
  content: string;
}

export default function App() {
  useNostrUser()
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

  return (
    <SafeAreaView style={appStyles.safeArea}>
      <View style={appStyles.container}>
        <Header isConnected={isConnected} />
        <View style={appStyles.content}>
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
      </View>
    </SafeAreaView>
  );
}