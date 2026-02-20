export interface Graph {
	id: string;
	title: string;
	description?: string;
	starting_at: Date;

	nodes?: Node[];
	edges?: Edge[];

	created_at: Date;
	updated_at: Date;
}

export interface Node<T extends NodeType, D> {
	id: string;
	type: T;

	title: string;
	description?: string;

	data?: D;

	created_at: Date;
	updated_at: Date;
	unlocked_at?: Date;
}

export interface Edge {
	id: string;

	from: string;
	to: string;

	choice_label?: string;

	created_at: Date;
	updated_at: Date;
}
