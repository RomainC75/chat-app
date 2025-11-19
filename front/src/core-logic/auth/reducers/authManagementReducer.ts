import { createReducer } from "@reduxjs/toolkit";
import type { AppState } from "../../../store/appState";
import type { TUser } from "../../../types/user.type";
import { userSignedUp } from "../use-cases/signup";


const setInitialState = (user: TUser | null): AppState["authManagement"] => ({
  data: user,
  error: null,
})


const authManagementReducer = (user: TUser | null) => createReducer(
  setInitialState(user),
  (builder) => {
    builder
      .addCase(userSignedUp, (state, action) => {
        state.data = action.payload;
        state.error = null;
      })
  },
);

export const AuthManagement = (user: TUser|null = null) =>
  authManagementReducer(user);
