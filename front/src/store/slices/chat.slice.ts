import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { TAvailableRoom, TChatSlice } from "../../types/chat.type";

const initialState: TChatSlice = {
  publicRoom: null,
  privateRoom: null,
  availableRooms: [
    
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
    privateRoomCreated(state, action: PayloadAction<TAvailableRoom>){
      state.privateRoom = action.payload
    },
    roomListReceived(state, action: PayloadAction<TAvailableRoom[]>){
      console.log("-> roomList received : ", action.payload)
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
