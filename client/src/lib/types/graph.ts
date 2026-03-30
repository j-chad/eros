export type GraphStatus = 'not_started' | 'in_progress' | 'completed';

export interface GraphSummary {
	id: string;
	title: string;
	description?: string;
	status?: GraphStatus;
	starting_at: Date;
	created_at: Date;
	updated_at: Date;
}

/** Full graph with accessible nodes and edges — returned by GET /api/graphs/:id */
export interface GraphDetail {
	id: string;
	title: string;
	description?: string;
	starting_at: string; // ISO 8601 date string
	nodes?: AnyNode[];
	edges?: Edge[];
	created_at: string; // ISO 8601 date string
	updated_at: string; // ISO 8601 date string
}

export enum NodeType {
	START = 'start',
	LOCATION = 'location',
	CODE = 'code',
	MANUAL = 'manual',
	REWARD = 'reward',
}

export enum RewardType {
	IMAGE = 'image',
	VIDEO = 'video',
	CALENDAR = 'calendar',
	WALLET = 'wallet',
	FAVOUR = 'favour',
	FILE = 'file',
	MARKDOWN = 'markdown',
}

export interface Node<Type extends NodeType = NodeType, Data = undefined> {
	id: string;
	type: Type;

	title: string;
	description?: string | null;

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
	reward_type: RewardType;
	payload: string;
	give_favours: number;
}>;

export type AnyNode = StartNode | LocationNode | CodeNode | ManualNode | RewardNode;

export interface Edge {
	id: string;
	from: string;
	to: string;
	choice_label?: string;
	created_at: string; // ISO 8601 date string
	updated_at: string; // ISO 8601 date string
}
