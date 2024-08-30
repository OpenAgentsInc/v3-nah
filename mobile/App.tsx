import "text-encoding-polyfill"
import { StatusBar } from "expo-status-bar"
import { nip19 } from "nostr-tools"
import { Image, LogBox, StyleSheet, Text, View } from "react-native"
import { useStore } from "@/lib/store"
import { useNostrUser } from "./lib/useNostrUser"

LogBox.ignoreLogs(["Promise"])

export default function App() {
  useNostrUser()
  const userPubkey = useStore(state => state.userPubkey)
  return (
    <View style={styles.container}>
      <Image source={require('./assets/sqlogo-t.png')} style={styles.image} resizeMode="contain" />
      {userPubkey && <Text style={styles.text}>{nip19.npubEncode(userPubkey)}</Text>}
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
  image: {
    width: 200,
    height: 200,
  },
  text: {
    color: '#fff',
    fontFamily: 'Courier New',
    fontSize: 18,
    paddingTop: 20,
    paddingHorizontal: 20,
    textAlign: 'center',
    fontWeight: 'bold'
  }
});
