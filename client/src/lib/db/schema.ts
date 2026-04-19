import type { GraphDetail, GraphSummary } from '$lib/types/graph';
import type { FavourChoice, FavourRequest } from '$lib/types/favour';

export const DB_NAME = 'eros';
export const DB_VERSION = 5;

export interface DBSchema {
	kv: {
		key: string;
		value: unknown;
	};
	graphs: {
		key: string;
		value: GraphSummary;
	};
	graphDetails: {
		key: string;
		value: GraphDetail;
	};
	favourChoices: {
		key: string;
		value: FavourChoice;
	};
	favourRequests: {
		key: string;
		value: FavourRequest;
	};
}

interface StoreParameters extends IDBObjectStoreParameters {
	indexes?: { name: string; keyPath: string | string[]; options?: IDBIndexParameters }[];
}

const STORES: Record<keyof DBSchema, StoreParameters> = {
	kv: {
		keyPath: 'key',
	},
	graphs: {
		keyPath: 'id',
		indexes: [
			{
				name: 'byStartDate',
				keyPath: 'starting_at',
				options: { unique: false },
			},
		],
	},
	graphDetails: {
		keyPath: 'id',
	},
	favourChoices: {
		keyPath: 'id',
	},
	favourRequests: {
		keyPath: 'id',
	},
}

export function createObjectStores(db: IDBDatabase) {
	for (const [name, options] of Object.entries(STORES)) {
		if (!db.objectStoreNames.contains(name)) {
			const store = db.createObjectStore(name, options);

			if (options.indexes) {
				for (const index of options.indexes) {
					store.createIndex(index.name, index.keyPath, index.options);
				}
			}

			console.log(`Created object store: ${name}`);
		}
	}
}
