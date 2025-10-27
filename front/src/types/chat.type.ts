export type TChatSlice = {
  publicRoom: TPublicRoom | null;
  privateRoom: TPrivateRoom | null;
  availableRooms: TAvailableRoom[];
  isConnected: boolean;
};

export type TAvailableRoom = {
  id: string;
  name: string;
  description?: string;
  memberCount: number;
  isPrivate: boolean;
  createdAt: Date;
  lastActivity: Date;
};

export type TPublicRoom = {

  users: TChatUsers[];
  messages: TChatMessage[];
};

export type TPrivateRoom = TPublicRoom & {
  id: string;
  name: string;
};

export type TChatUsers = {
  id: string;
  email: string;
};

export type TChatMessage = {
  id: string;
  sender: TChatUsers;
  content: string;
  date: Date;
};
