import { createAction } from "@reduxjs/toolkit";
import type { AppThunk } from "../../store/store";
import type { TLogin } from "../../core-logic/auth/types/auth.type";

export const userLoggedIn =
  createAction<TLogin>("LOGGED_ID");

export const alreadyValidated = createAction<void>("ALREADY_VALIDATED");

export const login =
  (email: string, password: string): AppThunk<Promise<void>> =>
  async (dispatch, getState, { authGateway }) => {
    if (getState().authManagement.data != null) {
      dispatch(alreadyValidated());
      return;
    }
    const logValidation = await authGateway.login(
      email,
      password,
    );
    dispatch(userLoggedIn(logValidation));
  };
