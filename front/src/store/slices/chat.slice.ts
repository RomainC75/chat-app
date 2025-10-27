import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { TChatSlice, TPrivateRoom } from "../../types/chat.type";

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
        users: [],
        messages: [],
      } as TPrivateRoom
    }
    // addRoom(state, action: PayloadAction<TRoom>) {
    //       state.rooms.push(action.payload)
    //     },
  },
});

export default chatSlice.reducer;

export const {
    setConnected,
    privateRoomCreated
} = chatSlice.actions;
