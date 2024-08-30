import "text-encoding-polyfill"
import { StatusBar } from "expo-status-bar"
import { Image, StyleSheet, View } from "react-native"
import { useStore } from "@/lib/store"
import { useNostrUser } from "./lib/useNostrUser"

export default function App() {
  useNostrUser()
  const userPubkey = useStore(state => state.userPubkey);
  console.log(userPubkey);
  return (
    <View style={styles.container}>
      <Image source={require('./assets/sqlogo-t.png')} style={styles.image} resizeMode="contain" />
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
  }
});
