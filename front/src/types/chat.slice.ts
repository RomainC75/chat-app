
export type TChatSlice = {
    rooms: TRoom[];
    isConnected: boolean;
}

export type TRoom = {
    users: any;
    id: string;

}