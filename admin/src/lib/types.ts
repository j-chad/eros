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
    cost: number;
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
    LOCATION = 'location',
    CODE = 'code',
    MANUAL = 'manual',
    REWARD = 'reward',
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

    nodes: AnyNode[];
    edges: Edge[];

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

export interface Node<Type extends NodeType = NodeType, Data = never> {
    id: string;
    type: Type,

    title: string;
    description?: string | null;

    ui_position?: {
        x: number;
        y: number;
    };

    data?: Data;

    unlocked_at?: string | null; // ISO 8601 date string

    created_at: string; // ISO 8601 date string
    updated_at: string; // ISO 8601 date string
}

export type StartNode = Node<NodeType.START>;
export type LocationNode = Node<NodeType.LOCATION, {
    latitude: number;
    longitude: number;
    radius_m: number;
}>;
export type CodeNode = Node<NodeType.CODE, {
    code: string;
}>;
export type ManualNode = Node<NodeType.MANUAL, {
	instructions: string;
	unlocked_at?: string | null; // ISO 8601 date string
}>;
export type RewardNode = Node<NodeType.REWARD, {
    reward_type: string; // not sure what the possible types are yet
    content: string;
    media_type: string;
    give_favours: number;
}>;

export type AnyNode = StartNode | LocationNode | CodeNode | ManualNode | RewardNode;

export interface NewGraph {
    title: string;
    description?: string | null;
    starting_at: string; // ISO 8601 date string
}
