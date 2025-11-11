import { IsRoom, ObjectKeys, Rule } from "../../shared/rule";
import type { TAvailableRoom, TPrivateRoom } from "../../types/chat.type";
import { EWsMessageIn } from "../../types/socket.type";

export type TMessageIn = {
  type: string;
};

type TRawContent = {
  [key: string]: string;
};

type TRawMessageIn = {
  type: string;
  content: TRawContent;
};

export type TParsedContent = {
  [key: string]: object | string;
};

export const ParseContent = (content: TRawContent): TParsedContent => {
  console.log("-> content", content, Object.keys(content));
  const result: TParsedContent = {};
  Object.keys(content).forEach((key) => {
    console.log("parser : ", content[key]);
    result[key] = JSON.parse(content[key]);
  });
  return result;
};

export type MessageInHandlerDeps = {
  privateRoomCreated: (privateRoom: TPrivateRoom) => void;
  roomListReceived: (availableRooms: TAvailableRoom[]) => void;
};

export const MessageInHandler =
  ({ privateRoomCreated, roomListReceived }: MessageInHandlerDeps) =>
  (messageIn: TRawMessageIn) => {
    const parsedContent = ParseContent(messageIn.content);
    switch (messageIn.type) {
      case EWsMessageIn.HELLO:
        console.log("-> HELLO connected");
        return;
      case EWsMessageIn.ROOM_CREATED:
        Rule.checkAll([
          new ObjectKeys(parsedContent, ["room_id", "room_name", "users"]),
        ]);
        privateRoomCreated(parsedContent.content as TPrivateRoom);
        return;
      case EWsMessageIn.ROOMS_LIST:
        Rule.checkAll([new ObjectKeys(parsedContent, ["rooms_list"])]);
        // IsArray.of(parsedContent, IsRoom)
        for (const el of parsedContent["rooms_list"] as TPrivateRoom[]) {
          Rule.checkAll([new IsRoom(el)]);
        }
        roomListReceived(parsedContent["rooms_list"] as TAvailableRoom[]);
        return;
      case EWsMessageIn.NEW_USER_CONNECTED_TO_CHAT:
        console.log("-> new user connected to chat")
        return;
      default:
        console.log("not found ! ");
    }
  };
