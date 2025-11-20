import { createAction } from "@reduxjs/toolkit";


export const errorRaised =
  createAction<string>("ERROR_RAISED");

export const isLoading = 
      createAction<void>("IS_LOADING");    