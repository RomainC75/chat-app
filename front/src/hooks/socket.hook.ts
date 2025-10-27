import { useEffect } from 'react'
import useWebSocket, { ReadyState } from 'react-use-websocket'
import { useDispatch } from 'react-redux'
import { type AppDispatch } from '../store/store'
import { EWsMessageIn, type IwebSocketMessageIn } from '../types/socket.type'
import { privateRoomCreated, setConnected } from '../store/slices/chat.slice'

export const useSocket = () => {
    const dispatch = useDispatch<AppDispatch>()

        const token = localStorage.getItem("token")!
        const SOCKET_URL = import.meta.env.VITE_SOCKET_URL;
        console.log("->W socker url ", SOCKET_URL)
        const { sendMessage: sendWsMessage, lastMessage, readyState } =
            useWebSocket<IwebSocketMessageIn>(
              `${SOCKET_URL}/api/chat/ws?token=${token}`, {
    shouldReconnect: () => true,
    reconnectAttempts: 10,
    reconnectInterval: 3000,
  }
            );

  useEffect(() => {
    dispatch(setConnected(readyState === ReadyState.OPEN))
  }, [readyState, dispatch])

  
  useEffect(() => {
    if (lastMessage) {
        try {
            const message = JSON.parse(lastMessage.data)
            console.log("last message : ", message)
        switch(message.type){
            case EWsMessageIn.ROOM_CREATED:
                dispatch(privateRoomCreated(message.content))
            
        }
        // dispatch(addMessage(data))
      } catch {
        console.error("error trying to decode this message : ", lastMessage)
      }
    }
  }, [lastMessage, dispatch])

  return { sendWsMessage, readyState }
}