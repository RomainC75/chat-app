import type { TChatSlice } from "../types/chat.type";
import type { TUserSlice } from "../types/user.type";

export interface AppState {
  user: TChatSlice;
  chat: TUserSlice;
}