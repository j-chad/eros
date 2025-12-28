export interface APIError {
    code: string;
    message: string;
    details?: Record<string, unknown>;
}

export interface RegistrationToken {
    code: string;
    created_at: string; // ISO 8601 date string
    expires_at: string; // ISO 8601 date string
}

export interface Device {
    id: string;
    name: string;
    registered_at: string; // ISO 8601 date string
    last_seen_at: string; // ISO 8601 date string
    expires_at: string; // ISO 8601 date string
}