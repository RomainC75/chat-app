import {
  configureStore,
  type Action,
  type Store,
  type ThunkAction,
  type ThunkDispatch,
} from "@reduxjs/toolkit";
import type { AppState } from "./appState.ts";
import { gateways, type Gateways } from "../adapters/secondary/gatewaysConfig.ts";
import { AuthManagement } from "../core-logic/auth/reducers/authManagementReducer.ts";
import type { TLogin } from "../core-logic/auth/types/auth.type.ts";

export const initReduxStore = (config: {
    gateways?: Partial<Gateways>
    userInit?: TLogin
}) => {
  return configureStore({
    reducer: {
      authManagement: AuthManagement(config.userInit ?? null ),
    },
    middleware: (getDefaultMiddleware) => {
      const mergedGateways: Gateways = {
        ...gateways,
        ...(config.gateways ?? {}),
      } as Gateways;
      return getDefaultMiddleware({
        thunk: {
          extraArgument: mergedGateways,
        },
        serializableCheck: false,
      });
    },
    devTools: true,
  });
};

export type ReduxStore = Store<AppState> & {
  dispatch: ThunkDispatch<AppState, Gateways, Action>;
};

export type AppThunk<ReturnType = void> = ThunkAction<
  ReturnType,
  AppState,
  Gateways,
  Action
>;

export type AppDispatch = ThunkDispatch<AppState, Gateways, Action>;


// export const store = configureStore({
//     reducer: {
//         user: userReducer,
//         chat: chatReducer,
//     },
// });

// export type RootState = ReturnType<typeof store.getState>;
// export type AppDispatch = typeof store.dispatch;

// export * from "./thunks/user.thunk";

