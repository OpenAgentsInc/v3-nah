import { useEffect } from "react"
import { useStore } from "@/lib/store"

export function useNostrUser() {
  const initializeNDK = useStore(state => state.initializeNDK);
  const userPubkey = useStore(state => state.userPubkey);

  useEffect(() => {
    initializeNDK();
  }, [initializeNDK]);

  useEffect(() => {
    if (userPubkey) {
      console.log(`You are ${userPubkey}`);
    }
  }, [userPubkey]);
}
