import type { TParsedContent } from "../message/input/messageInHandler";
import type { TPrivateRoom } from "../types/chat.type";

// ==

export abstract class Rule {
  static checkAll(rules: Rule[]) {
    rules.forEach((rule) => {
      if (!rule.isRespected) {
        rule.raiseError();
      }
    });
  }
  abstract isRespected(): boolean;
  protected abstract raiseError(): Error;
}

export class ObjectKeys extends Rule {
  constructor(
    private readonly obj: TParsedContent,
    private readonly keys: string[]
  ) {
    super();
  }
  isRespected(): boolean {
    this.keys.forEach((key) => {
      if (!this.obj[key]) {
        return false;
      }
    });
    return true;
  }
  raiseError(): Error {
    return new Error(`a key is missing`);
  }
}

export class IsRoom extends Rule {
  constructor(private readonly obj: TPrivateRoom) {
    super();
  }
  isRespected(): boolean {
      const {id, name, description, created_at} = this.obj
    if (
      !id || !name || !description || !created_at
    ) {
      return false;
    }
    return true;
  }
  raiseError(): Error {
    return new Error(`a key is missing`);
  }
}