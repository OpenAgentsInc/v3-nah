import { useEffect } from 'react';

interface MessageHandlerProps {
  socket: WebSocket | null;
  setTranscription: (transcription: string) => void;
  setAgentResponse: (response: string) => void;
  setIsProcessing: (isProcessing: boolean) => void;
}

export const useMessageHandler = ({
  socket,
  setTranscription,
  setAgentResponse,
  setIsProcessing,
}: MessageHandlerProps) => {
  useEffect(() => {
    if (socket) {
      const messageHandler = (event: MessageEvent) => {
        const data = JSON.parse(event.data);
        console.log("Received data:", data);
        if (Array.isArray(data) && data[0] === "EVENT") {
          const eventData = data[1];
          console.log("Event data:", eventData);
          if (eventData.kind === 6252) {
            setTranscription(eventData.content);
          } else if (eventData.kind === 6838) {
            setAgentResponse(eventData.content);
            setIsProcessing(false);
          }
        }
      };
      socket.addEventListener('message', messageHandler);
      return () => {
        socket.removeEventListener('message', messageHandler);
      };
    }
  }, [socket, setTranscription, setAgentResponse, setIsProcessing]);
};