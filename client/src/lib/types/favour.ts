export interface FavourChoice {
	id: string;
	label: string;
	description: string | null;
	cost: number;
	can_message: boolean;
	created_at: string;
	updated_at: string;
}

export interface FavourCount {
	total: number;
	remaining: number;
}

export interface FavourRequest {
	id: string;
	choice_id: string;
	choice_label: string;
	choice_description: string | null;
	message: string | null;
	requested_at: string;
	fulfilled_at: string | null;
}
