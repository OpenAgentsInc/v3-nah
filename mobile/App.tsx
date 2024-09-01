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

export default function App() {
  useNostrUser()
  const { isConnected, socket } = useRelayConnection()
  const { startRecording, stopRecording, isRecording } = useAudioRecording()
  const [transcription, setTranscription] = useState<string | null>(null)
  const [agentResponse, setAgentResponse] = useState<string | null>(null)
  const [isProcessing, setIsProcessing] = useState(false)
  const permissionStatus = useAudioPermission()

  useMessageHandler({ socket, setTranscription, setAgentResponse, setIsProcessing })

  const handlePressIn = useCallback(async () => {
    if (permissionStatus !== 'granted') {
      console.log('Audio permission not granted')
      return
    }
    setTranscription(null)
    setAgentResponse(null)
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
          setTranscription(receivedTranscription)
        })
      }
    } catch (error) {
      console.error('Error sending audio:', error)
      setTranscription('Error sending audio')
      setIsProcessing(false)
    }
  }, [socket, stopRecording])

  return (
    <SafeAreaView style={appStyles.safeArea}>
      <View style={appStyles.container}>
        <Header isConnected={isConnected} />
        <View style={appStyles.content}>
          <TranscriptionDisplay
            transcription={transcription}
            agentResponse={agentResponse}
          />
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