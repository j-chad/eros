export interface APIError {
    code: string;
    message: string;
    details?: Record<string, unknown>;
	internal?: string;
}

function isAPIError(error: unknown): error is APIError {
	return typeof error === 'object' && error !== null && 'code' in error && 'message' in error;
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
    TIME = 'time',
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

export interface Node<Type extends NodeType = NodeType, Data = undefined> {
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
export interface LocationHint {
    latitude: number;
    longitude: number;
    radius_m: number;
}

export type LocationNode = Node<NodeType.LOCATION, {
    latitude: number;
    longitude: number;
    radius_m: number;
    show_hint: boolean;
    hint?: LocationHint | null;
}>;
export type CodeNode = Node<NodeType.CODE, {
    codes: string[];
}>;
export type ManualNode = Node<NodeType.MANUAL, {
	instructions: string;
	unlocked_at?: string | null; // ISO 8601 date string
}>;
export type TimeNode = Node<NodeType.TIME, {
	unlock_at: string; // ISO 8601 date string (UTC)
}>;

export enum RewardType {
	IMAGE = 'image',
	VIDEO = 'video',
	CALENDAR = 'calendar',
	WALLET = 'wallet',
	FAVOUR = 'favour',
	FILE = 'file',
	MARKDOWN = 'markdown',
}

export type RewardNode = Node<NodeType.REWARD, {
    reward_type: RewardType;
    payload: string;
    give_favours: number;
}>;

export type AnyNode = StartNode | LocationNode | CodeNode | ManualNode | TimeNode | RewardNode;
export type NodeByType<T extends NodeType> = Extract<AnyNode, { type: T }>;
export type NodeDataByType<T extends NodeType> = NodeByType<T>['data'];

export interface NewGraph {
    title: string;
    description?: string | null;
    starting_at: string; // ISO 8601 date string
}
