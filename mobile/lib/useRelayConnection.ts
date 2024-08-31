import { useState, useEffect } from 'react';
import { useStore } from '@/lib/store';

const RELAY_URL = 'ws://localhost:8080'; // Replace with your actual relay URL

export function useRelayConnection() {
  const [socket, setSocket] = useState<WebSocket | null>(null);
  const [isConnected, setIsConnected] = useState(false);
  const userPubkey = useStore(state => state.userPubkey);

  useEffect(() => {
    const ws = new WebSocket(RELAY_URL);

    ws.onopen = () => {
      console.log('Connected to relay');
      setIsConnected(true);
      setSocket(ws);
    };

    ws.onclose = () => {
      console.log('Disconnected from relay');
      setIsConnected(false);
      setSocket(null);
    };

    ws.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log('Received message:', message);
      // Handle incoming messages here
    };

    return () => {
      ws.close();
    };
  }, []);

  useEffect(() => {
    if (socket && isConnected && userPubkey) {
      // Subscribe to events for the user's pubkey
      const subscriptionMessage = JSON.stringify(['REQ', 'user_events', { authors: [userPubkey] }]);
      socket.send(subscriptionMessage);
    }
  }, [socket, isConnected, userPubkey]);

  return { isConnected, socket };
}