
import { createSlice, type PayloadAction } from "@reduxjs/toolkit";
import { signupUser } from "../thunks/user.thunk";
import type { TUserSlice } from "../../types/user.type";

const initialState: TUserSlice = {
    isLoading: false,
    errorMessage: null,
    user: null,
}

const userSlice = createSlice({
    name: "user",
    initialState,
    reducers: {},

    extraReducers(builder) {
        builder.addCase(signupUser.pending, (state) => {
            state.isLoading = true;
        });
        builder.addCase(signupUser.fulfilled, (state, action: PayloadAction<string>) => {
            state.isLoading = false;
            console.log("--> action ", action)
            // state.user = action.payload;
        });
        builder.addCase(signupUser.rejected, (state) => {
            state.isLoading = false;
        });
    }
});

export default userSlice.reducer;
