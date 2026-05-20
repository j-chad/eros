export interface PushMessage {
    title: string
    body: string
    tag?: string
    data?: {
        url?: string
    }
}