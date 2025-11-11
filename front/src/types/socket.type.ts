export const EWsMessageIn = {
  NEW_ROOM_MESSAGE: "NEW_ROOM_MESSAGE",
  NEW_BROADCAST_MESSAGE: "NEW_BROADCAST_MESSAGE",
  MEMBER_JOINED: "MEMBER_JOINED",
  MEMBER_LEAVED: "MEMBER_LEAVED",
  ROOM_CREATED: "ROOM_CREATED",
  CONNECTED_TO_ROOM: "CONNECTED_TO_ROOM",
  DISCONNECTED_FROM_ROOM: "DISCONNECTED_FROM_ROOM",
  NEW_USER_CONNECTED_TO_CHAT: "NEW_USER_CONNECTED_TO_CHAT",
  ROOMS_LIST: "ROOMS_LIST",
  HELLO: "HELLO",
} as const;

export type EWsMessageIn = (typeof EWsMessageIn)[keyof typeof EWsMessageIn];

export const EWsMessageOut = {
  roomMessage: "MESSAGE",
  broadcastMessage: "BROADCAST_MESSAGE",
  createRoom: "CREATE_ROOM",
  connectToRoom: "CONNECT_TO_ROOM",
  disconnectFromRoom: "DISCONNECT_FROM_ROOM",
} as const;

export type EWsMessageOut = (typeof EWsMessageOut)[keyof typeof EWsMessageOut];

export interface IwebSocketMessageOut {
  type: EWsMessageOut;
  content: Record<string, unknown>;
}

export interface IwebSocketMessageIn {
  type: EWsMessageIn;
  content: IWebSocketMessageContent;
}

export interface IWebSocketMessageContent {
  message: string;
  userEmail: string;
  userId: string;
}

export interface WSMessage {
  type: string;
  content: {
    message: string;
    userEmail?: string;
    userId?: string;
  };
}

export interface IRoom {
  name: string;
  id: string;
}

export interface IGame {
  name: string;
  id: string;
  playerNumber: number;
}

export interface IGameStateReducer {
  status: string;
  state: IGameState | null;
}

export interface IGameState {
  bait: IPosition;
  players: [IPlayerState, IPlayerState];
  lever: number;
  game_config: {
    size: number;
    speed_ms: number;
  };
  lastCommands: number[];
}

export interface IGameConfig {
  size: number;
  speed_ms: number;
}

export interface IPlayerState {
  score: number;
  positions: IPosition[];
  direction: number;
}

export interface IPosition {
  x: number;
  y: number;
}

export interface IGridDot {
  color: string | undefined;
}
