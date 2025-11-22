import type { TLogin } from "../core-logic/auth/types/auth.type";

// export interface AppState {
//   user: TChatSlice;
//   chat: TUserSlice;
// }

export interface AppState {
  authManagement: {
    data: TLogin | null,
    error: string | null;
    isLoading: boolean;
  };
}


export type TPokemon = {
    name:string;
    url: string;
}