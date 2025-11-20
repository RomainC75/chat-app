import { createAction } from "@reduxjs/toolkit";


export const errorRaised =
  createAction<string>("ERROR_RAISED");