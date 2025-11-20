import { createReducer } from "@reduxjs/toolkit";
import type { AppState } from "../../../store/appState";
import { userLoggedIn } from "../../../use-cases/auth/login/login";
import type { TLogin } from "../types/auth.type";
import { errorRaised, isLoading } from "../../../use-cases/error";
import { userSignedUp } from "../../../use-cases/auth/signup/signup";


const setInitialState = (user: TLogin | null): AppState["authManagement"] => ({
  data: user,
  error: null,
  isLoading: false
})


const authManagementReducer = (user: TLogin | null) => createReducer(
  setInitialState(user),
  (builder) => {
    builder
      .addCase(userSignedUp, (state) => {
        state.data = null;
        state.isLoading = false ;
        state.error = null;
      })
      .addCase(userLoggedIn, (state, action) => {
        state.data = action.payload;
        state.isLoading = false;
        state.error = null;
      })
      .addCase(errorRaised, (state, action)=>{
        state.data = null;
        state.isLoading = false;
        state.error = action.payload
      })
      .addCase(isLoading, (state)=>{
        state.isLoading =true;
      })
  },
);

export const AuthManagement = (user: TLogin | null = null) =>
  authManagementReducer(user);
