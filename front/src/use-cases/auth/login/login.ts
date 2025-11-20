import { createAction } from "@reduxjs/toolkit";
import type { AppThunk } from "../../../store/store";
import type { TLogin } from "../../../core-logic/auth/types/auth.type";
import { errorRaised, isLoading } from "../../error";

export const userLoggedIn =
  createAction<TLogin>("LOGGED_ID");

export const alreadyValidated = createAction<void>("ALREADY_VALIDATED");

export const login =
  (email: string, password: string): AppThunk<Promise<void>> =>
  async (dispatch, getState, { authGateway }) => {
    // if (getState().authManagement.data != null) {
    //   dispatch(alreadyValidated());
    //   return;
    // }
    dispatch(isLoading())
    try {
      const logValidation = await authGateway.login(
        email,
        password,
      );
      
      dispatch(userLoggedIn(logValidation));
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error)

      dispatch(errorRaised(message))
    }
  };
