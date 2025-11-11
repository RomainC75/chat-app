import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { TAvailableRoom, TChatSlice, TPrivateRoom } from "../../types/chat.type";

const initialState: TChatSlice = {
  publicRoom: null,
  privateRoom: null,
  availableRooms: [
    {
        id: "1",
        name: "General Chat",
        description: "General discussion for everyone",
        memberCount: 24,
        isPrivate: false,
        createdAt: new Date("2024-01-15"),
        lastActivity: new Date(),
      },
  ],
  isConnected: false
};

const chatSlice = createSlice({
  name: "chat",
  initialState,

  reducers: {
    setConnected(state, action: PayloadAction<boolean>) {
          state.isConnected = action.payload
        },
    privateRoomCreated(state, action: PayloadAction<TPrivateRoom>){
      state.privateRoom = {
        id: action.payload.id,
        name: action.payload.name,
        
      } as TPrivateRoom
    },
    roomListReceived(state, action: PayloadAction<TAvailableRoom[]>){
      state.availableRooms = action.payload
    }
  },
});

export default chatSlice.reducer;

export const {
    setConnected,
    privateRoomCreated,
    roomListReceived,
} = chatSlice.actions;
