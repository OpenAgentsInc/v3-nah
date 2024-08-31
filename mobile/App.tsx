import "text-encoding-polyfill"
import { StatusBar } from "expo-status-bar"
import { nip19 } from "nostr-tools"
import { StyleSheet, Text, View } from "react-native"
import { useStore } from "@/lib/store"
import { useNostrUser } from "./lib/useNostrUser"
import { useRelayConnection } from "./lib/useRelayConnection"
import { useAudioRecording } from "./lib/useAudioRecording"
import PushToTalkButton from "./components/PushToTalkButton"

export default function App() {
  useNostrUser()
  const userPubkey = useStore(state => state.userPubkey)
  const { isConnected } = useRelayConnection()
  const { startRecording, stopRecording, isRecording } = useAudioRecording()

  const handlePressIn = async () => {
    await startRecording()
  }

  const handlePressOut = async () => {
    const audioUri = await stopRecording()
    if (audioUri) {
      console.log('Audio recorded:', audioUri)
      // TODO: Send audio to relay for processing
    }
  }

  return (
    <View style={styles.container}>
      {userPubkey && <Text style={styles.text}>{nip19.npubEncode(userPubkey)}</Text>}
      <Text style={styles.connectionStatus}>
        Relay: {isConnected ? 'Connected' : 'Disconnected'}
      </Text>
      <Text style={styles.recordingStatus}>
        {isRecording ? 'Recording...' : 'Press and hold to speak'}
      </Text>
      <PushToTalkButton onPressIn={handlePressIn} onPressOut={handlePressOut} />
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
  }
});