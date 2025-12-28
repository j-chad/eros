export interface APIError {
    code: string;
    message: string;
    details?: Record<string, unknown>;
}

export interface RegistrationToken {
    token: string;
    expiresAt: string; // ISO 8601 date string
}