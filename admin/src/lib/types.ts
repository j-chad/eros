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
    device_info: string
}

export interface FavourChoice {
    id: string;
    label: string;
    description: string;
    can_message: boolean;
    created_at: string; // ISO 8601 date string
    updated_at: string; // ISO 8601 date string
}
export type CreateFavourChoicePayload = Omit<FavourChoice, 'id' | 'created_at' | 'updated_at'>;

export type FavourCount = {
    total: number;
    remaining: number;
};

export type Favour = {
    id: string;
    choice_id: string;
    choice_label: string;
    choice_description: string | null;
    message: string | null;
    requested_at: string;
    fulfilled_at: string | null;
};

export enum NodeType {
    START = 'start',
}

export interface Graph {
    id: string;
    title: string;
    description?: string | null;

    viewport?: {
        x: number;
        y: number;
        zoom: number;
    }

    nodes?: Array<Node>;
    edges?: Array<Edge>;

    starting_at: string; // ISO 8601 date string
    created_at: string; // ISO 8601 date string
    updated_at: string; // ISO 8601 date string
}

export interface Edge {
    id: string;
    from: string;
    to: string;

    choice_label?: string;

    created_at: string; // ISO 8601 date string
    updated_at: string; // ISO 8601 date string
}

export interface Node<T extends NodeType = NodeType, D extends Record<string, unknown> = never> {
    id: string;
    type: T,

    title: string;
    description?: string | null;

    ui_position?: {
        x: number;
        y: number;
    };

    data?: D;

    unlocked_at?: string | null; // ISO 8601 date string

    created_at: string; // ISO 8601 date string
    updated_at: string; // ISO 8601 date string
}

export interface NewGraph {
    title: string;
    description?: string | null;
    starting_at: string; // ISO 8601 date string
}