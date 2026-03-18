export interface GraphSummary {
	id: string;
	title: string;
	description?: string;
	starting_at: Date;
	created_at: Date;
	updated_at: Date;
}

export interface GraphTypes {
	id: string;
	title: string;
	description?: string;
	starting_at: Date;

	nodes?: AnyNode[];
	edges?: Edge[];

	created_at: Date;
	updated_at: Date;
}

export enum NodeType {
	START = 'start',
	LOCATION = 'location',
	CODE = 'code',
	MANUAL = 'manual',
	REWARD = 'reward',
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
	payload: string;
	give_favours: number;
}>;

export type AnyNode = StartNode | LocationNode | CodeNode | ManualNode | RewardNode;

export interface Edge {
	id: string;

	from: string;
	to: string;

	choice_label?: string;

	created_at: Date;
	updated_at: Date;
}
