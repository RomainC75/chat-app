export type TSignup = {
    id: string
    email: string
}

export type TLogin = TSignup & {
    token: string
}