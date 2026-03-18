import type { GraphSummary } from '$lib/types/graph';

export const DB_NAME = 'eros';
export const DB_VERSION = 2;

// Graphs are stored with date fields as ISO strings (as returned by the API).
export type StoredGraph = Omit<GraphSummary, 'starting_at' | 'created_at' | 'updated_at'> & {
	starting_at: string;
	created_at: string;
	updated_at: string;
};

export interface DBSchema {
	kv: {
		key: string;
		value: unknown;
	};
	graphs: {
		key: string;
		value: StoredGraph;
	};
}

interface StoreParameters extends IDBObjectStoreParameters {
	indexes?: { name: string; keyPath: string | string[]; options?: IDBIndexParameters }[];
}

const STORES:Record<keyof DBSchema, StoreParameters> = {
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
			}
		]
	}
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
