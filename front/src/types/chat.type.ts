export type TChatSlice = {
  publicRoom: TAvailableRoom | null;
  privateRoom: TAvailableRoom | null;
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
