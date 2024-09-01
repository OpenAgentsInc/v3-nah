import "text-encoding-polyfill"
import { Audio } from "expo-av"
import { StatusBar } from "expo-status-bar"
import React, { useCallback, useEffect, useState } from "react"
import { Image, SafeAreaView, StyleSheet, Text, View } from "react-native"
import { useStore } from "@/lib/store"
import PushToTalkButton from "./components/PushToTalkButton"
import RelayStatusIcon from "./components/RelayStatusIcon"
import { sendAudioToRelay } from "./lib/sendAudioToRelay"
import { useAudioRecording } from "./lib/useAudioRecording"
import { useNostrUser } from "./lib/useNostrUser"
import { useRelayConnection } from "./lib/useRelayConnection"

export default function App() {
  useNostrUser()
  const { isConnected, socket } = useRelayConnection()
  const { startRecording, stopRecording, isRecording } = useAudioRecording()
  const [transcription, setTranscription] = useState<string | null>(null)
  const [agentResponse, setAgentResponse] = useState<string | null>(null)
  const [isProcessing, setIsProcessing] = useState(false)
  const [permissionStatus, setPermissionStatus] = useState<'granted' | 'denied' | 'pending'>('pending')

  useEffect(() => {
    (async () => {
      const { status } = await Audio.requestPermissionsAsync()
      setPermissionStatus(status === 'granted' ? 'granted' : 'denied')
    })()
  }, [])

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
          // The 5838 event is now sent automatically in sendAudioToRelay
        })
      }
    } catch (error) {
      console.error('Error sending audio:', error)
      setTranscription('Error sending audio')
      setIsProcessing(false)
    }
  }, [socket, stopRecording])

  useEffect(() => {
    if (socket) {
      const messageHandler = (event: MessageEvent) => {
        const data = JSON.parse(event.data)
        console.log("Received data:", data)
        if (Array.isArray(data) && data[0] === "EVENT") {
          const eventData = data[1]
          console.log("Event data:", eventData)
          if (eventData.kind === 6252) {
            setTranscription(eventData.content)
          } else if (eventData.kind === 6838) {
            setAgentResponse(eventData.content)
            setIsProcessing(false)
          }
        }
      }
      socket.addEventListener('message', messageHandler)
      return () => {
        socket.removeEventListener('message', messageHandler)
      }
    }
  }, [socket])

  return (
    <SafeAreaView style={styles.safeArea}>
      <View style={styles.container}>
        <View style={styles.header}>
          <Image source={require('./assets/sqlogo-t.png')} style={styles.logo} resizeMode="contain" />
          <RelayStatusIcon isConnected={isConnected} />
        </View>
        <View style={styles.content}>
          <View style={styles.transcriptionContainer}>
            {transcription && (
              <Text style={styles.transcription}>Transcription: {transcription}</Text>
            )}
            {agentResponse && (
              <Text style={styles.agentResponse}>Agent Response: {agentResponse}</Text>
            )}
          </View>
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

const styles = StyleSheet.create({
  safeArea: {
    flex: 1,
    backgroundColor: '#000',
  },
  container: {
    flex: 1,
    backgroundColor: '#000',
  },
  header: {
    flexDirection: 'row',
    justifyContent: 'space-between',
    alignItems: 'center',
    paddingTop: 0,
    paddingHorizontal: 20,
    height: 50,
  },
  logo: {
    width: 40,
    height: 40,
  },
  content: {
    flex: 1,
    justifyContent: 'flex-start',
    alignItems: 'center',
    paddingHorizontal: 20,
    marginTop: 10,
  },
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