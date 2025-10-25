import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import { loginUser, signupUser, verify } from "../thunks/user.thunk";
import type { TLoginResponse, TUserSlice, TVerifyResponse } from "../../types/user.type";

const initialState: TUserSlice = {
  isLoading: false,
  errorMessage: null,
  user: null,
  isConnected: false
};

const userSlice = createSlice({
  name: "user",
  initialState,

  reducers: {},

  extraReducers(builder) {
    builder.addCase(signupUser.pending, (state) => {
      state.isLoading = true;
    });
    builder.addCase(
      signupUser.fulfilled,
      (state, action: PayloadAction<string>) => {
        state.isLoading = false;
        console.log("--> action ", action);
        // state.user = action.payload;
      }
    );
    builder.addCase(signupUser.rejected, (state) => {
      state.isLoading = false;
    });

    // --------- LOGIN ---------
    builder.addCase(loginUser.pending, (state) => {
      state.isLoading = true;
    });
    builder.addCase(
      loginUser.fulfilled,
      (state, action: PayloadAction<TLoginResponse>) => {
        state.isLoading = false;
        localStorage.setItem("token", action.payload.token)
      }
    );
    builder.addCase(loginUser.rejected, (state) => {
      state.isLoading = false;
    });

    // --------- VERIFY ---------
    builder.addCase(verify.pending, (state) => {
      state.isLoading = true;
    });
    builder.addCase(
      verify.fulfilled,
      (state, action: PayloadAction<TVerifyResponse>) => {
        state.isLoading = false;
        state.isConnected = true;
      }
    );
    builder.addCase(verify.rejected, (state) => {
      state.isLoading = false;
      state.isConnected = false;
    });
  },
});

export default userSlice.reducer;
