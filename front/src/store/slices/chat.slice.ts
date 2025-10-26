import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import type { TChatSlice, TRoom } from "../../types/chat.slice";

const initialState: TChatSlice = {
  rooms: [],
  isConnected: false
};

const chatSlice = createSlice({
  name: "chat",
  initialState,

  reducers: {
    addRoom(state, action: PayloadAction<TRoom>) {
          state.rooms.push(action.payload)
        },
  },

  // extraReducers(builder) {

  // },
});

export default chatSlice.reducer;
