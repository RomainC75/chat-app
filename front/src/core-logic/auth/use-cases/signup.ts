import { createAction } from "@reduxjs/toolkit";
import type { AppThunk } from "../../../store/store.ts";
import type { TUser } from "../../../types/user.type.ts";

export const userSignedUp = createAction<TUser>("SIGNUP");

export const signup =
  (): AppThunk<Promise<void>> =>
  async (dispatch, _, { authGateway }) => {
    const pokemons = await authGateway.signup();
    console.log("-> retrieved : ", pokemons)
    dispatch(userSignedUp(pokemons));
  };