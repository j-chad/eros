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

/** Full graph with accessible nodes and edges - returned by GET /api/graphs/:id */
export interface GraphDetail {
	id: string;
	title: string;
	description?: string;
	starting_at: string; // ISO 8601 date string
	nodes: AnyNode[];
	edges: Edge[];
	created_at: string; // ISO 8601 date string
	updated_at: string; // ISO 8601 date string
}

export enum NodeType {
	START = 'start',
	LOCATION = 'location',
	CODE = 'code',
	MANUAL = 'manual',
	TIME = 'time',
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
	URL = 'url',
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
export interface LocationHint {
	latitude: number;
	longitude: number;
	radius_m: number;
}

export type LocationNode = Node<NodeType.LOCATION, {
	radius_m: number;
	hint?: LocationHint | null;
}>;
export type CodeNode = Node<NodeType.CODE, {
	code: string;
}>;
export type ManualNode = Node<NodeType.MANUAL, {
	instructions: string;
	unlocked_at?: string | null; // ISO 8601 date string
}>;
export type TimeNode = Node<NodeType.TIME, {
	unlock_at: string; // ISO 8601 date string (UTC)
}>;
export interface FileInfo {
	id: string;
	filename: string;
	mime_type: string;
	size_bytes: number;
	url: string;
	url_expires_at?: string;
}

export type RewardNode = Node<NodeType.REWARD, {
	reward_type: RewardType;
	payload: string;
	give_favours: number;
	file?: FileInfo;
}>;

export type AnyNode = StartNode | LocationNode | CodeNode | ManualNode | TimeNode | RewardNode;

export interface Edge {
	id: string;
	from: string;
	to: string;
	choice_label?: string;
	created_at: string; // ISO 8601 date string
	updated_at: string; // ISO 8601 date string
}
