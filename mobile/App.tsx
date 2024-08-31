import "text-encoding-polyfill"
import React, { useState, useEffect } from 'react'
import { StatusBar } from "expo-status-bar"
import { nip19 } from "nostr-tools"
import { StyleSheet, Text, View, Image } from "react-native"
import { Audio } from 'expo-av'
import { useStore } from "@/lib/store"
import { useNostrUser } from "./lib/useNostrUser"
import { useRelayConnection } from "./lib/useRelayConnection"
import { useAudioRecording } from "./lib/useAudioRecording"
import { sendAudioToRelay } from "./lib/sendAudioToRelay"
import PushToTalkButton from "./components/PushToTalkButton"

export default function App() {
  useNostrUser()
  const userPubkey = useStore(state => state.userPubkey)
  const { isConnected, socket } = useRelayConnection()
  const { startRecording, stopRecording, isRecording } = useAudioRecording()
  const [transcription, setTranscription] = useState<string | null>(null)
  const [isProcessing, setIsProcessing] = useState(false)
  const [permissionStatus, setPermissionStatus] = useState<'granted' | 'denied' | 'pending'>('pending')

  useEffect(() => {
    (async () => {
      const { status } = await Audio.requestPermissionsAsync()
      setPermissionStatus(status === 'granted' ? 'granted' : 'denied')
    })()
  }, [])

  const handlePressIn = async () => {
    if (permissionStatus !== 'granted') {
      console.log('Audio permission not granted')
      return
    }
    setTranscription(null)
    await startRecording()
  }

  const handlePressOut = async () => {
    if (!socket) {
      console.error('No WebSocket connection available')
      return
    }

    setIsProcessing(true)
    try {
      const audioUri = await stopRecording()
      if (audioUri) {
        console.log('Audio recorded:', audioUri)
        const result = await sendAudioToRelay(audioUri, socket)
        setTranscription(result)
      }
    } catch (error) {
      console.error('Error processing audio:', error)
      setTranscription('Error processing audio')
    } finally {
      setIsProcessing(false)
    }
  }

  return (
    <View style={styles.container}>
      <Image source={require('./assets/sqlogo-t.png')} style={styles.logo} resizeMode="contain" />
      {userPubkey && <Text style={styles.text}>{nip19.npubEncode(userPubkey)}</Text>}
      <Text style={styles.connectionStatus}>
        Relay: {isConnected ? 'Connected' : 'Disconnected'}
      </Text>
      <Text style={styles.recordingStatus}>
        {isRecording ? 'Recording...' : isProcessing ? 'Processing...' : 'Press and hold to speak'}
      </Text>
      {transcription && (
        <Text style={styles.transcription}>Transcription: {transcription}</Text>
      )}
      <PushToTalkButton 
        onPressIn={handlePressIn} 
        onPressOut={handlePressOut}
        disabled={permissionStatus !== 'granted'}
      />
      <StatusBar style="light" />
    </View>
  );
}

const styles = StyleSheet.create({
  container: {
    flex: 1,
    backgroundColor: '#000',
    alignItems: 'center',
    justifyContent: 'center',
  },
  logo: {
    width: 200,
    height: 200,
    marginBottom: 20,
  },
  text: {
    color: '#fff',
    fontFamily: 'Courier New',
    fontSize: 18,
    paddingTop: 20,
    paddingHorizontal: 20,
    textAlign: 'center',
    fontWeight: 'bold'
  },
  connectionStatus: {
    color: '#fff',
    fontFamily: 'Courier New',
    fontSize: 14,
    paddingTop: 10,
    textAlign: 'center',
  },
  recordingStatus: {
    color: '#fff',
    fontFamily: 'Courier New',
    fontSize: 16,
    paddingTop: 20,
    textAlign: 'center',
  },
  transcription: {
    color: '#fff',
    fontFamily: 'Courier New',
    fontSize: 14,
    paddingTop: 20,
    paddingHorizontal: 20,
    textAlign: 'center',
  }
});