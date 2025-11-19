import type { TUser } from "../types/user.type";

// export interface AppState {
//   user: TChatSlice;
//   chat: TUserSlice;
// }

export interface AppState {
  authManagement: {
    data: TUser | null,
    error: string | null;
  };
}


export type TPokemon = {
    name:string;
    url: string;
}