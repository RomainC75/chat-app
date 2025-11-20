import { createReducer } from "@reduxjs/toolkit";
import type { AppState } from "../../../store/appState";
import { userSignedUp } from "../use-cases/signup";
import { userLoggedIn } from "../../../use-cases/login/login";
import type { TLogin } from "../types/auth.type";


const setInitialState = (user: TLogin | null): AppState["authManagement"] => ({
  data: user,
  error: null,
})


const authManagementReducer = (user: TLogin | null) => createReducer(
  setInitialState(user),
  (builder) => {
    builder
      .addCase(userSignedUp, (state) => {
        state.error = null;
      })
      .addCase(userLoggedIn, (state, action) => {
        state.data = action.payload;
        state.error = null;
      })
  },
);

export const AuthManagement = (user: TLogin | null = null) =>
  authManagementReducer(user);
