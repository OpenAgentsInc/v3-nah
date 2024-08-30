import "text-encoding-polyfill"
import * as SecureStore from "expo-secure-store"
import { generatePrivateKey, getPublicKey } from "nostr-tools_1_1_1"
import { create } from "zustand"
import { createJSONStorage, persist } from "zustand/middleware"
import NDK, { NDKEvent, NDKPrivateKeySigner } from "@nostr-dev-kit/ndk"

const expoSecureStorage = {
  setItem: SecureStore.setItemAsync,
  getItem: SecureStore.getItemAsync,
  removeItem: SecureStore.deleteItemAsync,
}

interface State {
  userPubkey: string | null,
  userSecret: string | null,
  ndkInstance: NDK | null,
  events: NDKEvent[],
  setUserPubkey: (pubkey: string) => void,
  setUserSecret: (secret: string) => void,
  addEvent: (event: NDKEvent) => void,
  getEventsInReverseChronologicalOrder: () => NDKEvent[],
  initializeNDK: () => void
}

export const useStore = create<State>()(
  persist(
    (set, get) => ({
      userPubkey: null,
      userSecret: null,
      ndkInstance: null,
      events: [],
      setUserPubkey: (userPubkey) => set({ userPubkey }),
      setUserSecret: (userSecret) => set({ userSecret }),
      addEvent: (event: NDKEvent) => set((state) => ({
        events: [...state.events, event],
      })),
      getEventsInReverseChronologicalOrder: () => {
        return get().events.slice().sort((a, b) => b.created_at - a.created_at);
      },
      initializeNDK: async () => {
        const { userSecret, setUserPubkey, setUserSecret, addEvent } = get();

        let sk, pk;

        if (userSecret) {
          // Use existing keys
          sk = userSecret;
          pk = getPublicKey(sk);
          setUserPubkey(pk);
        } else {
          // Generate new keys
          sk = generatePrivateKey(); // `sk` is a hex string
          pk = getPublicKey(sk); // `pk` is a hex string
          setUserPubkey(pk);
          setUserSecret(sk);
        }

        const ndk = new NDK({
          explicitRelayUrls: [
            "wss://magency.nostr1.com",
          ],
          enableOutboxModel: true,
        });

        ndk.pool?.on("relay:connecting", (relay) => {
          console.log("🪄 MAIN POOL Connecting to relay", relay.url);
        });

        ndk.pool?.on("relay:connect", (relay) => {
          console.log("✅ MAIN POOL Connected to relay", relay.url);
        });

        ndk.signer = new NDKPrivateKeySigner(sk);

        set({ ndkInstance: ndk });

        await ndk.connect(5000);

        // Define the event kinds to subscribe to
        const eventKinds = [38000, 38001, 38002, 38003];

        // Subscribe to the specified event kinds
        const subscription = ndk.subscribe(
          { kinds: eventKinds },
          { closeOnEose: false }
        );

        // Listen for events and log them as they are received
        subscription.on("event", (event: NDKEvent) => {
          console.log("Received event:", event.id);
          addEvent(event); // Add the received event to the store
        });
      },
    }),
    {
      name: "nostr-keys",
      storage: createJSONStorage(() => expoSecureStorage),
    }
  )
);
