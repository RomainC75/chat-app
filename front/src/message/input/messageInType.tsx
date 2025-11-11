export class MessageInType {
  static valideTypes = {
    ROOM_MESSAGE: "ROOM_MESSAGE",
    BROADCAST_MESSAGE: "BROADCAST_MESSAGE",
    CONNECT_TO_ROOM: "CONNECT_TO_ROOM",
    CREATE_ROOM: "CREATE_ROOM",
    DISCONNECT_FROM_ROOM: "DISCONNECT_FROM_ROOM",
  };

  private constructor(private readonly miType: string) {}

  static fromString(miType: string): MessageInType {
    if (!Object.values(this.valideTypes).find(v=>v===miType)) {
      throw Error("messageType invalid");
    }
    return new MessageInType(miType);
  }

  asString(): string {
    return this.miType
  }
};