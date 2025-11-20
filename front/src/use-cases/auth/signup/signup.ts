import { createAction } from "@reduxjs/toolkit";
import type { AppThunk } from "../../../store/store";
import { errorRaised, isLoading } from "../../error";

export const userSignedUp =
  createAction<void>("SIGNED_UP");


export const signup =
  (email: string, password: string): AppThunk<Promise<void>> =>
  async (dispatch, getState, { authGateway }) => {
    // if (getState().authManagement.data != null) {
    //   dispatch(alreadyValidated());
    //   return;
    // }
    dispatch(isLoading())
    try {
      await authGateway.signup(
        email,
        password,
      );
      dispatch(userSignedUp());
    } catch (error) {
      const message = error instanceof Error ? error.message : String(error)
      dispatch(errorRaised(message))
    }
  };
